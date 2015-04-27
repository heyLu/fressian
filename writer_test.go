package fressian

import (
	"bytes"
	"reflect"
	"testing"

	tu "github.com/klingtnet/gol/util/testing"
)

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
	val, err := r.ReadValue()
	tu.ExpectNil(t, err)
	tu.ExpectEqual(t, b, val)
}

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
	val, err := r.ReadValue()
	tu.ExpectNil(t, err)
	tu.ExpectEqual(t, i, val)
}

func TestWriteValue(t *testing.T) {
	testWriteValue(t, nil)
	testWriteValue(t, true)
	testWriteValue(t, false)
	testWriteValue(t, 3)
	//testWriteValue(t, []int{1, 2, 3})
	testWriteValue(t, "hello")
	testWriteValue(t, "Hello, World!")
	testWriteValue(t, Key{"hello", "world"})
	testWriteValue(t, []interface{}{1, 2, true, 4})
}

func testWriteValue(t *testing.T, val interface{}) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf, nil)
	w.WriteValue(val)
	w.Flush()

	tu.ExpectNil(t, w.Error())

	r := NewReader(buf, nil)
	res, err := r.ReadValue()
	tu.ExpectNil(t, err)
	if !reflect.DeepEqual(val, res) {
		t.Errorf("Expected reflect.DeepEqual(%#v, %#v)", val, res)
	}
}
