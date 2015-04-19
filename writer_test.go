package fressian

import (
	"bytes"
	"testing"

	tu "github.com/klingtnet/gol/util/testing"
)

func TestWriteInt(t *testing.T) {
	testWriteInt(t, 0)
	testWriteInt(t, 1)
	testWriteInt(t, 10)
	testWriteInt(t, 583)
	testWriteInt(t, 36342523521)
}

func testWriteInt(t *testing.T, i int) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf, nil)
	w.WriteInt(i)
	w.Flush()

	tu.ExpectNil(t, w.Error())

	r := NewReader(buf, nil)
	val, err := r.ReadObject()
	tu.ExpectNil(t, err)
	tu.ExpectEqual(t, i, val)
}

func TestWriteBool(t *testing.T) {
	testWriteBool(t, true)
	testWriteBool(t, false)
}

func testWriteBool(t *testing.T, b bool) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf, nil)
	w.WriteBool(b)
	w.Flush()

	tu.ExpectNil(t, w.Error())

	r := NewReader(buf, nil)
	val, err := r.ReadObject()
	tu.ExpectNil(t, err)
	tu.ExpectEqual(t, b, val)
}
