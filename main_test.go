package main

import (
	"bytes"
	"testing"
)

func TestReadInt(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x01}))
	i := r.readInt()
	if r.Err() != nil {
		t.Fatal(r.Err())
	}
	if i != 1 {
		t.Errorf("%d != %d", i, 1)
	}
}
