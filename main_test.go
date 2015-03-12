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

func TestReadObject(t *testing.T) {
	expectReadObject(t, []byte{0x01}, 1)
	expectReadObject(t, []byte{0x05}, 5)
	expectReadObject(t, []byte{0x09}, 9)
	expectReadObject(t, []byte{0x21}, 33)
	expectReadObject(t, []byte{0x3A}, 58)
	// INT_PACKED_2
	expectReadObject(t, []byte{0x40, 0x00}, -4096)
	expectReadObject(t, []byte{0x50, 0x00}, 0)
	expectReadObject(t, []byte{0x5F, 0xFF}, 4095)
	// INT_PACKED_3
	expectReadObject(t, []byte{0x60, 0x00, 0x00}, -524288)
	expectReadObject(t, []byte{0x68, 0x00, 0x00}, 0)
	expectReadObject(t, []byte{0x6F, 0xFF, 0xFF}, 524287)
	// INT_PACKED_4
	expectReadObject(t, []byte{0x70, 0x00, 0x00, 0x00}, -33554432)
	expectReadObject(t, []byte{0x72, 0x00, 0x00, 0x00}, 0)
	expectReadObject(t, []byte{0x73, 0xFF, 0xFF, 0xFF}, 33554431)
	// INT_PACKED_5
	expectReadObject(t, []byte{0x74, 0x00, 0x00, 0x00, 0x00}, -8589934592)
	expectReadObject(t, []byte{0x77, 0xFF, 0xFF, 0xFF, 0xFF}, 8589934591)
	// INT_PACKED_6
	expectReadObject(t, []byte{0x78, 0x00, 0x00, 0x00, 0x00, 0x00}, -2199023255552)
	expectReadObject(t, []byte{0x7B, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, 2199023255551)
	// INT_PACKED_7
	expectReadObject(t, []byte{0x7C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, -562949953421312)
	expectReadObject(t, []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, 562949953421311)

	expectReadObject(t, []byte{TRUE}, true)
	expectReadObject(t, []byte{FALSE}, false)

	expectReadObject(t, []byte{NULL}, nil)

	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START}, "")
	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START + 1, 0x61}, "a")
	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START + 3, 0x61, 0x62, 0x63}, "abc")
	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START + 5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}, "hello")

	readObjectList(t, []byte{LIST_PACKED_LENGTH_START}, []interface{}{})
	readObjectList(t, []byte{LIST_PACKED_LENGTH_START + 1, 0x01}, []interface{}{1})
	readObjectList(t, []byte{LIST_PACKED_LENGTH_START + 3, 0x07, 0x04, 0x09}, []interface{}{7, 4, 9})
}

func expectReadObject(t *testing.T, bs []byte, res interface{}) {
	r := newReader(bs)
	obj := r.readObject()
	tu.RequireNil(t, r.Err())
	tu.ExpectEqual(t, obj, res)
}

func readObjectList(t *testing.T, bs []byte, res []interface{}) []interface{} {
	obj := readObject(t, bs)
	list, ok := obj.([]interface{})
	if !ok {
		t.Fatal("readObject did not return a list!", obj)
	}
	if len(list) != len(res) {
		t.Fatalf("len(list) = %d != %d", len(list), len(res))
	}
	for i, expected := range res {
		tu.ExpectEqual(t, list[i], expected)
	}
	return list
}

func readObject(t *testing.T, bs []byte) interface{} {
	r := newReader(bs)
	obj := r.readObject()
	tu.RequireNil(t, r.Err())
	return obj
}

func newReader(bs []byte) *Reader {
	return NewReader(bytes.NewReader(bs))
}
