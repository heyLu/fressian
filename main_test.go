package main

import (
	"bytes"
	"testing"

	tu "github.com/klingtnet/gol/util/testing"
)

func TestReadInt(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x01}))
	i := r.readInt()
	tu.RequireNil(t, r.Err())
	tu.ExpectEqual(t, i, 1)
}
