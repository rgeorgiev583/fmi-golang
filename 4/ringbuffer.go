package main

type RingBuffer struct {
    buffer []interface{}
    beginPos, endPos int
}

func (rb *RingBuffer) Slice() []interface{} {
    if rb.beginPos <= rb.endPos {
        return rb.buffer[rb.beginPos:rb.endPos+1]
    } else {
        return append(buffer[rb.beginPos:], rb.buffer[:rb.endPos+1])
    }
}

func (rb *RingBuffer) Append(item interface{}) {
    if rb.beginPos == rb.endPos + 1 {
        rb.beginPos++
    }

    if rb.endPos == len(rb.buffer) - 1 {
        rb.endPos = 0
    } else {
        rb.endPos++
    }

    if rb.beginPos == len(rb.buffer) - 1 {
        rb.beginPos = 0
    }

    rb.buffer[rb.endPos] = item
}

func NewRingBuffer(size int) *RingBuffer {
    return &RingBuffer{buffer: make([]interface{}, size)}
}
