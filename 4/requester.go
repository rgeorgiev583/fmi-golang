package main

import (
    "fmt"
    "sync"
)

const (
    bufferSize = 100
)

type Requester interface {
    // Добавя заявка за изпълнение и я изпълнява, ако това е необходимо, при първа възможност.
    AddRequest(request Request)

    // Спира 'Заявчика'. Това означава, че изчаква всички вече започнали заявки да завършат
    // и извиква `SetResult` на тези заявки, които вече са били добавени, но "равни" на тях вече са били изпълнявание.
    // Нови заявки не трябва да бъдат започвани през това време, нито вече започнати, равни на тях, да бъдат добавяни за извикване на `SetResult`.
    Stop()
}

type SimpleRequester struct {
    queue chan Request
    cache *RingBuffer
    throttle chan struct{}
    lock sync.RWMutex
    waiter sync.WaitGroup
    isStopped bool
    classes map[string]chan Request
}

type CachedResult struct {
    ID string
    Result interface{}
    Err error
}

type RequestNotFoundError struct {
    ID string
}

func (rnfe *RequestNotFoundError) Error() string {
    return fmt.Sprintf("could not find result of request with ID %v in cache", rnfe.ID)
}

func (sr *SimpleRequester) FindResult(id string) (*CachedResult, error) {
    cacheLen := sr.cache.Length()

    for i := 0; i < cacheLen; i++ {
        cacheItem, _ := sr.cache.Item(i)

        if result, ok := cacheItem.(CachedResult); ok && result.ID == id {
            return &result, nil
        }
    }

    return nil, &RequestNotFoundError{id}
}


func (sr *SimpleRequester) AddRequest(request Request) {
    if !sr.isStopped {
        sr.queue <- request
    }
}

func (sr *SimpleRequester) Stop() {
    sr.isStopped = true
    close(sr.queue)
    for range sr.queue {}
    sr.waiter.Wait()
    close(sr.throttle)
    for range sr.throttle {}
}

// Връща нов заявчик, който кешира отговорите на до cacheSize заявки,
// изпълнявайки не повече от throttleSize заявки едновременно.
func NewRequester(cacheSize int, throttleSize int) Requester {
    sr := &SimpleRequester{
        cache: NewRingBuffer(cacheSize),
        queue: make(chan Request, bufferSize),
        throttle: make(chan struct{}, throttleSize),
        classes: make(map[string]chan Request),
    }

    for i := 0; i < throttleSize; i++ {
        sr.throttle <- struct{}{}
    }

    sr.waiter.Add(1)

    go func() {
        for request := range sr.queue {
            sr.waiter.Add(1)

            go func() {
                defer sr.waiter.Done()
                id := request.ID()
                sr.lock.Lock()
                class, ok := sr.classes[id]
                sr.lock.Unlock()

                if ok {
                    class <- request
                    return
                }

                cachedResult, _ := sr.FindResult(id)

                if cachedResult != nil {
                    request.SetResult(cachedResult.Result, cachedResult.Err)
                    return
                }

                <-sr.throttle
                class = make(chan Request, bufferSize)
                sr.lock.Lock()
                sr.classes[id] = class
                sr.lock.Unlock()
                result, err := request.Run()
                close(class)

                if request.Cacheable() {
                    sr.cache.Append(&CachedResult{
                        ID: id,
                        Result: result,
                        Err: err,
                    })

                    for identicalRequest := range class {
                        identicalRequest.SetResult(result, err)
                    }
                } else if !sr.isStopped {
                    for identicalRequest := range class {
                        sr.queue <- identicalRequest
                    }
                }

                sr.lock.Lock()
                delete(sr.classes, id)
                sr.lock.Unlock()
                sr.throttle <- struct{}{}
            }()
        }

        sr.waiter.Done()
    }()

    return sr
}
