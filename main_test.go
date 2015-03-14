package main

import (
	"bytes"
	"testing"
	"time"

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

	readObjectBytes(t, []byte{BYTES, 0x03, 0x01, 0x0A, 0xf9}, []byte{0x01, 0x0A, 0xf9})

	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START}, "")
	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START + 1, 0x61}, "a")
	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START + 3, 0x61, 0x62, 0x63}, "abc")
	expectReadObject(t, []byte{STRING_PACKED_LENGTH_START + 5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}, "hello")
	expectReadObject(t, []byte{STRING, 0x0D, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21}, "Hello, World!")

	readObjectList(t, []byte{LIST_PACKED_LENGTH_START}, []interface{}{})
	readObjectList(t, []byte{LIST_PACKED_LENGTH_START + 1, 0x01}, []interface{}{1})
	readObjectList(t, []byte{LIST_PACKED_LENGTH_START + 3, 0x07, 0x04, 0x09}, []interface{}{7, 4, 9})
	readObjectList(t, []byte{LIST, 0x00}, []interface{}{})
	readObjectList(t, []byte{LIST, 0x0A, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}, []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	readObjectMap(t, []byte{MAP, LIST_PACKED_LENGTH_START}, map[interface{}]interface{}{})
	readObjectMap(t, []byte{MAP, LIST_PACKED_LENGTH_START + 2, 0x01, 0x02},
		map[interface{}]interface{}{1: 2})
	readObjectMap(t, []byte{MAP, LIST_PACKED_LENGTH_START + 6, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		map[interface{}]interface{}{1: 2, 3: 4, 5: 6})
	readObjectMap(t, []byte{MAP, LIST_PACKED_LENGTH_START + 4, STRING_PACKED_LENGTH_START + 3, 0x61, 0x62, 0x63, 0x2a, 0x07, 0x08},
		map[interface{}]interface{}{"abc": 42, 7: 8})

	obj := readObject(t, []byte{INST, 0x7b, 0x4c, 0x0f, 0x1e, 0xcd, 0x76})
	date, ok := obj.(time.Time)
	if !ok {
		t.Fatalf("expected a time.Time, but got %#v", obj)
	}
	tu.ExpectEqual(t, date.Unix(), int64(1426182819190))

	readObjectTagged(t, []byte{KEY, STRING_PACKED_LENGTH_START + 2, 0x61, 0x62, STRING_PACKED_LENGTH_START + 1, 0x63}, Key{namespace: "ab", name: "c"})
}

func expectReadObject(t *testing.T, bs []byte, res interface{}) {
	r := newReader(bs)
	obj := r.readObject()
	tu.RequireNil(t, r.Err())
	tu.ExpectEqual(t, obj, res)
}

func readObjectBytes(t *testing.T, bs []byte, res []byte) []byte {
	bytes := readObject(t, bs).([]byte)
	tu.RequireEqual(t, len(bytes), len(res))
	for i, b := range res {
		tu.ExpectEqual(t, bytes[i], b)
	}
	return bytes
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

func readObjectMap(t *testing.T, bs []byte, res map[interface{}]interface{}) map[interface{}]interface{} {
	obj := readObject(t, bs)
	m, ok := obj.(map[interface{}]interface{})
	if !ok {
		t.Fatal("readObject did not return a map!", obj)
	}
	if len(m) != len(res) {
		t.Fatalf("len(m) = %d != %d", len(m), len(res))
	}
	for k, v := range res {
		mv, ok := m[k]
		if !ok {
			t.Errorf("key %#v not present", k)
			continue
		}
		tu.ExpectEqual(t, mv, v)
	}
	return m
}

func readObjectTagged(t *testing.T, bs []byte, res Tagged) Tagged {
	tagged := readObject(t, bs).(Tagged)
	tu.ExpectEqual(t, tagged.Key(), res.Key())
	tu.RequireEqual(t, len(tagged.Value()), len(res.Value()))
	for i, val := range res.Value() {
		tu.ExpectEqual(t, tagged.Value()[i], val)
	}
	return tagged
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
