package main

import (
	"bytes"
	"testing"

	tu "github.com/klingtnet/gol/util/testing"
)

func TestReadInt(t *testing.T) {
	expectReadInt(t, []byte{0x01}, 1)
}

func expectReadInt(t *testing.T, bs []byte, res int) {
	r := newReader(bs)
	i := r.readInt()
	tu.RequireNil(t, r.Err())
	tu.ExpectEqual(t, i, res)
}

func newReader(bs []byte) *Reader {
	return NewReader(bytes.NewReader(bs))
}
