package main

import (
    "bytes"
    "fmt"
)

type RingBuffer struct {
    buffer []interface{}
    beginPos, endPos int
}

type InvalidIndexError struct {
    index int
}

func (iie *InvalidIndexError) Error() string {
    return fmt.Sprintf("%d is not a valid index for a ring buffer element", iie.index)
}

func (rb *RingBuffer) Item(index int) (interface{}, error) {
    if index < 0 || index >= len(rb.buffer) || beginPos == 0 && index > endPos {
        return nil, &InvalidIndexError{index}
    }

    pos := rb.beginPos + index

    if pos >= len(rb.buffer) {
        pos -= len(rb.buffer)
    }

    return rb.buffer[pos], nil
}

func (rb *RingBuffer) Append(item interface{}) {
    if rb.beginPos == -1 || rb.endPos == -1 {
        rb.beginPos = 0
        rb.endPos = 0
        rb.buffer[0] = item
        return
    }

    if rb.endPos == len(rb.buffer) - 1 {
        rb.endPos = 0
    } else {
        rb.endPos++
    }

    if rb.beginPos == len(rb.buffer) - 1 {
        rb.beginPos = 0
    } else if rb.beginPos == rb.endPos {
        rb.beginPos++
    }

    rb.buffer[rb.endPos] = item
}

func (rb *RingBuffer) String() string {
    var buf bytes.Buffer
    buf.WriteString("[ ")

    for i := 0; i < len(rb.buffer) - 1; i++ {
        rbi, _ := rb.Item(i)
        buf.WriteString(fmt.Sprintf("%v ", rbi))
    }

    rbi, _ := rb.Item(len(rb.buffer) - 1)
    buf.WriteString(fmt.Sprintf("%v]", rbi))
    return buf.String()
}

func NewRingBuffer(size int) *RingBuffer {
    return &RingBuffer{
        buffer: make([]interface{}, size),
        beginPos: -1,
        endPos: -1,
    }
}
