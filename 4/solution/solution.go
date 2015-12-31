package main

import (
	"fmt"
    "sync"
)

const (
	bufferSize = 100
)

type Request interface {
	// Връща идентификатор за заявката. Ако две заявки имат еднакви идентификатори
	// то те са "равни".
	ID() string

	// Блокира докато изпълнява заявката.
	// Връща резултата или грешка ако изпълнението е неуспешно.
	// Резултата и грешката не трябва да бъдат подавани на SetResult
	// за текущата заявка - те се запазват вътрешно преди да бъдат върнати.
	Run() (result interface{}, err error)

	// Връща дали заявката е кешируерма.
	// Метода има неопределено поведение, ако бъде извикан преди `Run`.
	Cacheable() bool

	// Задава резултата на заявката.
	// Не трябва да се извиква за заявки, за които е бил извикан `Run`.
	SetResult(result interface{}, err error)
}

type RingBuffer struct {
	buffer           []interface{}
	beginPos, endPos int
}

type InvalidIndexError struct {
	index int
}

func (iie *InvalidIndexError) Error() string {
	return fmt.Sprintf("%d is not a valid index for an element in this buffer", iie.index)
}

func (rb *RingBuffer) Length() int {
	if rb.beginPos < rb.endPos {
		return rb.endPos - rb.beginPos + 1
	} else {
		return len(rb.buffer)
	}
}

func (rb *RingBuffer) Item(index int) (interface{}, error) {
	if index < 0 || index >= len(rb.buffer) || rb.beginPos == 0 && index > rb.endPos {
		return nil, &InvalidIndexError{index}
	}

	pos := rb.endPos - index

	if pos < 0 {
		pos += len(rb.buffer)
	}

	return rb.buffer[pos], nil
}

func (rb *RingBuffer) Append(value interface{}) {
	if rb.beginPos == -1 || rb.endPos == -1 {
		rb.beginPos = 0
		rb.endPos = 0
		rb.buffer[0] = value
		return
	}

	if rb.endPos == len(rb.buffer)-1 {
		rb.endPos = 0
	} else {
		rb.endPos++
	}

	if rb.beginPos == len(rb.buffer)-1 {
		rb.beginPos = 0
	} else if rb.beginPos == rb.endPos {
		rb.beginPos++
	}

	rb.buffer[rb.endPos] = value
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer:   make([]interface{}, size),
		beginPos: -1,
		endPos:   -1,
	}
}

type Requester interface {
	// Добавя заявка за изпълнение и я изпълнява, ако това е необходимо, при първа възможност.
	AddRequest(request Request)

	// Спира 'Заявчика'. Това означава, че изчаква всички вече започнали заявки да завършат
	// и извиква `SetResult` на тези заявки, които вече са били добавени, но "равни" на тях вече са били изпълнявание.
	// Нови заявки не трябва да бъдат започвани през това време, нито вече започнати, равни на тях, да бъдат добавяни за извикване на `SetResult`.
	Stop()
}

type SimpleRequester struct {
	queue     chan Request
	cache     *RingBuffer
	throttle  chan struct{}
	lock      sync.RWMutex
	waiter    sync.WaitGroup
	isStopped bool
	classes   map[string]chan Request
}

type CachedResult struct {
	ID     string
	Result interface{}
	Err    error
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
	for range sr.queue {
	}
	sr.waiter.Wait()
	close(sr.throttle)
	for range sr.throttle {
	}
}

// Връща нов заявчик, който кешира отговорите на до cacheSize заявки,
// изпълнявайки не повече от throttleSize заявки едновременно.
func NewRequester(cacheSize int, throttleSize int) Requester {
	sr := &SimpleRequester{
		cache:    NewRingBuffer(cacheSize),
		queue:    make(chan Request, bufferSize),
		throttle: make(chan struct{}, throttleSize),
		classes:  make(map[string]chan Request),
	}

	for i := 0; i < throttleSize; i++ {
		sr.throttle <- struct{}{}
	}

	sr.waiter.Add(1)

	go func() {
		for request := range sr.queue {
			sr.waiter.Add(1)

			go func(request Request) {
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
						ID:     id,
						Result: result,
						Err:    err,
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
			}(request)
		}

		sr.waiter.Done()
	}()

	return sr
}
