package fressian

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWriteBackAndForth(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, nil)
	writer.WriteInt(124124214211)
	//writer.WriteInt(1023)
	//writer.WriteInt(583)
	writer.Flush()
	fmt.Println(writer.Error())
	fmt.Printf("%#v\n", buf.Bytes())

	reader := NewReader(buf, nil)
	val, err := reader.ReadObject()
	fmt.Println(err)
	fmt.Println(val)
}
