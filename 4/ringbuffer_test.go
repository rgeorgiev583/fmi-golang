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

func appendNumbers(rb *RingBuffer, s []interface{}, count int, val *int, pos *int) {
	for i := 0; i < count; i++ {
		if *pos == -1 {
			for j := len(s) - 1; j > 0; j-- {
				s[j] = s[j-1]
			}

			*pos++
		}

		rb.Append(*val)
		s[*pos] = *val
		*val++
		*pos--
	}
}

func logRingBuffer(t *testing.T, rb *RingBuffer, s []interface{}, message string) {
	t.Log(message)
	t.Logf("ring buffer: %v", rb)
	t.Logf("slice: %v", s)
	t.Logf("items in ring buffer: %v", rb)
	t.Logf("items in slice: %v", s)
}

func testEqualityRingBufferSlice(t *testing.T, rb *RingBuffer, s []interface{}) {
	if len(s) != len(rb.buffer) {
		t.Errorf("ring buffer %v and slice %v should have had the same length", rb, s)
	}

	for i := range s {
		rbi, _ := rb.Item(i)

		if s[i] != rbi {
			t.Errorf("ring buffer %v and slice %v should have had the same data", rb, s)
		}
	}
}

func TestRingBufferAppend(t *testing.T) {
	rb := NewRingBuffer(10)
	s := make([]interface{}, 10)

	val := 1
	pos := len(s) - 1

	appendNumbers(rb, s, 10, &val, &pos)
	logRingBuffer(t, rb, s, "after adding 1..10")
	testEqualityRingBufferSlice(t, rb, s)

	appendNumbers(rb, s, 5, &val, &pos)
	logRingBuffer(t, rb, s, "after adding 11..15")
	testEqualityRingBufferSlice(t, rb, s)

	appendNumbers(rb, s, 10, &val, &pos)
	logRingBuffer(t, rb, s, "after adding 16..25")
	testEqualityRingBufferSlice(t, rb, s)
}
