package main

type RingBuffer struct {
    buffer []interface{}
    beginPos, endPos int
}

func (rb *RingBuffer) Slice() []interface{} {
    if rb.beginPos <= rb.endPos {
        return rb.buffer[rb.beginPos:rb.endPos+1]
    } else {
        return append(rb.buffer[rb.beginPos:], rb.buffer[:rb.endPos+1])
    }
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

func NewRingBuffer(size int) *RingBuffer {
    return &RingBuffer{
        buffer: make([]interface{}, size),
        beginPos: -1,
        endPos: -1,
    }
}
