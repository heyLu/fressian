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

	expectReadObject(t, []byte{TRUE}, true)
	expectReadObject(t, []byte{FALSE}, false)

	expectReadObject(t, []byte{NULL}, nil)

	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START}, "")

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
