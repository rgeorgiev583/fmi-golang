package main

import (
    "testing"
)

func TestNewRingBuffer(t *testing.T) {
    rb := NewRingBuffer(10)

    if rb == nil {
        t.Errorf("could not create RingBuffer struct")
    }
}

func testEqualitySlices(t *testing.T, s1 []interface{}, s2 []interface{}) {
    for i := range s1 {
        if s1[i] != s2[i] {
            t.Errorf("slices %v and %v should have had the same data", s1, s2)
        }
    }
}

func TestRingBufferAppend(t *testing.T) {
    rb := NewRingBuffer(10)
    s := make([]interface{}, 10)

    for i := 1; i <= 10; i++ {
        rb.Append(i)
        s = append(s, i)
    }

    testEqualitySlices(t, s, rb.Slice())

    for i := 0; i <= 4; i++ {
        val := 11 + i
        rb.Append(val)
        s[i] = val
    }

    testEqualitySlices(t, s, rb.Slice())

    for i := 0; i <= 4; i++ {
        val := 16 + i
        rb.Append(val)
        s[5 + i] = val
    }

    for i := 0; i <= 4; i++ {
        val := 21 + i
        rb.Append(val)
        s[i] = val
    }

    testEqualitySlices(t, s, rb.Slice())
}
