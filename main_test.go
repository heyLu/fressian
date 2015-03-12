package main

import (
	"bytes"
	"testing"

	tu "github.com/klingtnet/gol/util/testing"
)

func TestReadInt(t *testing.T) {
	expectReadInt(t, []byte{0x01}, 1)
	expectReadInt(t, []byte{0x02}, 2)
	expectReadInt(t, []byte{0x03}, 3)
	expectReadInt(t, []byte{0x04}, 4)
	expectReadInt(t, []byte{0x05}, 5)
	expectReadInt(t, []byte{0x06}, 6)
	expectReadInt(t, []byte{0x07}, 7)
	expectReadInt(t, []byte{0x08}, 8)
	expectReadInt(t, []byte{0x09}, 9)
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
