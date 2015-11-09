package main

import (
    "fmt"
)

func makeSemaphore(value int) chan struct{} {
    sem := make(chan struct{}, value)

    for i := 0; i < value; i++ {
        sem <- struct{}{}
    }

    return sem
}

func OrderedLogDrainer(logs chan (chan string)) chan string {
    mergedLogs := make(chan string, 100)

    go func() {
        logsSyncer := make(chan struct{})

        go func() {
            for log := range logs {
                go func() {
                    logSyncer := make(chan struct{})

                    go func() {
                        for logEntry := range log {
                            mergedLogs <- logEntry
                        }

                        <-logSyncer
                    }()

                    logSyncer <- struct{}{}
                    close(logSyncer)
                }()
            }

            <-logsSyncer
        }()

        logsSyncer <- struct{}{}
        close(logsSyncer)
        close(mergedLogs)
    }()

    return mergedLogs
}

func main() {
    logs := make(chan (chan string))
    orderedLog := OrderedLogDrainer(logs)

    first := make(chan string)
    logs <- first
    second := make(chan string)
    logs <- second

    first <- "test message 1 in first"
    second <- "test message 1 in second"
    second <- "test message 2 in second"
    first <- "test message 2 in first"
    first <- "test message 3 in first"
    // Print the first message now just because we can
    fmt.Println(<-orderedLog)

    third := make(chan string)
    logs <- third

    third <- "test message 1 in third"
    first <- "test message 4 in first"
    close(first)
    second <- "test message 3 in second"
    close(third)
    close(logs)

    second <- "test message 4 in second"
    close(second)

    // Print all the rest of the messages
    for logEntry := range orderedLog {
        fmt.Println(logEntry)
    }
}
