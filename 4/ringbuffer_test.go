package main

import (
    "bytes"
    "fmt"
    "testing"
)

func TestNewRingBuffer(t *testing.T) {
    rb := NewRingBuffer(10)

    if rb == nil {
        t.Errorf("could not create RingBuffer struct")
    }
}

func appendNumbers(rb *RingBuffer, s []interface{}, count int, val *int, pos *int) {
    for i := 0; i < count; i++ {
        if *pos == len(s) {
            for j := 0; j < len(s) - 1; j++ {
                s[j] = s[j + 1]
            }

            *pos--
        }

        rb.Append(*val)
        s[*pos] = *val
        *val++
        *pos++
    }
}

func logRingBuffer(t *testing.T, rb *RingBuffer, s []interface{}, message string) {
    t.Log(message)
    t.Logf("ring buffer: %v", rb)
    t.Logf("slice: %v", s)

    var buf bytes.Buffer
    buf.WriteString("items in ring buffer:")
    rbs := rb.Slice()

    for i := range rbs {
        buf.WriteString(fmt.Sprintf(", %d", rbs[i]))
    }

    t.Log(buf.String())

    buf.Reset()
    buf.WriteString("items in slice:")

    for i := range s {
        buf.WriteString(fmt.Sprintf(", %d", s[i]))
    }
}

func testEqualitySlices(t *testing.T, s1 []interface{}, s2 []interface{}) {
    if len(s1) != len(s2) {
        t.Errorf("slices %v and %v should have had the same length", s1, s2)
    }

    for i := range s1 {
        if s1[i] != s2[i] {
            t.Errorf("slices %v and %v should have had the same data", s1, s2)
        }
    }
}

func TestRingBufferAppend(t *testing.T) {
    rb := NewRingBuffer(10)
    s := make([]interface{}, 10)

    val := 1
    pos := 0

    appendNumbers(rb, s, 10, &val, &pos)
    logRingBuffer(t, rb, s, "after adding 1..10")
    testEqualitySlices(t, rb.Slice(), s)

    appendNumbers(rb, s, 5, &val, &pos)
    logRingBuffer(t, rb, s, "after adding 11..15")
    testEqualitySlices(t, rb.Slice(), s)

    appendNumbers(rb, s, 10, &val, &pos)
    logRingBuffer(t, rb, s, "after adding 16..25")
    testEqualitySlices(t, rb.Slice(), s)
}
