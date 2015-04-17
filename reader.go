// Package fressian supports reading values in the fressian format in
// Go.
package fressian

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"
)

// Byte codes for the types fressian supports
const (
	PRIORITY_CACHE_PACKED_START = 0x80
	PRIORITY_CACHE_PACKED_END   = 0xA0
	STRUCT_CACHE_PACKED_START   = 0xA0
	STRUCT_CACHE_PACKED_END     = 0xB0
	LONG_ARRAY                  = 0xB0
	DOUBLE_ARRAY                = 0xB1
	BOOLEAN_ARRAY               = 0xB2
	INT_ARRAY                   = 0xB3
	FLOAT_ARRAY                 = 0xB4
	OBJECT_ARRAY                = 0xB5
	MAP                         = 0xC0
	SET                         = 0xC1
	UUID                        = 0xC3
	REGEX                       = 0xC4
	URI                         = 0xC5
	BIGINT                      = 0xC6
	BIGDEC                      = 0xC7
	INST                        = 0xC8
	SYM                         = 0xC9
	KEY                         = 0xCA
	GET_PRIORITY_CACHE          = 0xCC
	PUT_PRIORITY_CACHE          = 0xCD
	PRECACHE                    = 0xCE
	FOOTER                      = 0xCF
	FOOTER_MAGIC                = 0xCFCFCFCF
	BYTES_PACKED_LENGTH_START   = 0xD0
	BYTES_PACKED_LENGTH_END     = 0xD8
	BYTES_CHUNK                 = 0xD8
	BYTES                       = 0xD9
	STRING_PACKED_LENGTH_START  = 0xDA
	STRING_PACKED_LENGTH_END    = 0xE2
	STRING_CHUNK                = 0xE2
	STRING                      = 0xE3
	LIST_PACKED_LENGTH_START    = 0xE4
	LIST_PACKED_LENGTH_END      = 0xEC
	LIST                        = 0xEC
	BEGIN_CLOSED_LIST           = 0xED
	BEGIN_OPEN_LIST             = 0xEE
	STRUCTTYPE                  = 0xEF
	STRUCT                      = 0xF0
	META                        = 0xF1
	ANY                         = 0xF4
	TRUE                        = 0xF5
	FALSE                       = 0xF6
	NULL                        = 0xF7
	INT                         = 0xF8
	FLOAT                       = 0xF9
	DOUBLE                      = 0xFA
	DOUBLE_0                    = 0xFB
	DOUBLE_1                    = 0xFC
	END_COLLECTION              = 0xFD
	RESET_CACHES                = 0xFE
	INT_PACKED_1_START          = 0xFF
	INT_PACKED_1_END            = 0x40
	INT_PACKED_2_START          = 0x40
	INT_PACKED_2_ZERO           = 0x50
	INT_PACKED_2_END            = 0x60
	INT_PACKED_3_START          = 0x60
	INT_PACKED_3_ZERO           = 0x68
	INT_PACKED_3_END            = 0x70
	INT_PACKED_4_START          = 0x70
	INT_PACKED_4_ZERO           = 0x72
	INT_PACKED_4_END            = 0x74
	INT_PACKED_5_START          = 0x74
	INT_PACKED_5_ZERO           = 0x76
	INT_PACKED_5_END            = 0x78
	INT_PACKED_6_START          = 0x78
	INT_PACKED_6_ZERO           = 0x7A
	INT_PACKED_6_END            = 0x7C
	INT_PACKED_7_START          = 0x7C
	INT_PACKED_7_ZERO           = 0x7E
	INT_PACKED_7_END            = 0x80
)

const (
	BYTE_CHUNK_SIZE = 65535
)

// Tagged is a generic interface for tagged data.
type Tagged interface {
	Key() string
	Value() []interface{}
}

// Key represends a fressian keyword, consisting of a namespace (which
// may be empty) and a name.
type Key struct {
	Namespace string
	Name      string
}

func (k Key) Key() string          { return "key" }
func (k Key) Value() []interface{} { return []interface{}{k.Namespace, k.Name} }

type structType struct {
	tag    string
	fields int
}

type StructAny struct {
	Tag    string
	Values []interface{}
}

type rawReader struct {
	br    *bufio.Reader
	count int
	err   error
}

func newRawReader(r io.Reader) *rawReader {
	return &rawReader{bufio.NewReader(r), 0, nil}
}

func (r *rawReader) readRawByte() byte {
	res, err := r.br.ReadByte()
	if err != nil {
		r.err = err
		return 0
	}
	r.count++
	return res
}

func (r *rawReader) readRawInt8() int {
	return int(r.readRawByte())
}

func (r *rawReader) readRawInt16() int {
	return (int(r.readRawByte()) << 8) + int(r.readRawByte())
}

func (r *rawReader) readRawInt24() int {
	return (int(r.readRawByte()) << 16) + (int(r.readRawByte()) << 8) + int(r.readRawByte())
}

func (r *rawReader) readRawInt32() int {
	return (int(r.readRawByte()) << 24) + (int(r.readRawByte()) << 16) + (int(r.readRawByte()) << 8) + int(r.readRawByte())
}

func (r *rawReader) readRawInt40() int {
	return (r.readRawInt8() << 32) | r.readRawInt32()
}

func (r *rawReader) readRawInt48() int {
	return (r.readRawInt16() << 32) | r.readRawInt32()
}

func (r *rawReader) readRawInt64() int {
	return ((int(r.readRawByte()) & 0xff) << 56) |
		((int(r.readRawByte()) & 0xff) << 48) |
		((int(r.readRawByte()) & 0xff) << 40) |
		((int(r.readRawByte()) & 0xff) << 32) |
		((int(r.readRawByte()) & 0xff) << 24) |
		((int(r.readRawByte()) & 0xff) << 16) |
		((int(r.readRawByte()) & 0xff) << 8) |
		(int(r.readRawByte()) & 0xff)
}

// ReadHandler is an alias for custom handlers of tagged data.
//
// A handler MUST read fieldCount values when called.
type ReadHandler func(r *Reader, tag string, fieldCount int) interface{}

// Reader reads fressian values from another io.Reader
type Reader struct {
	raw           *rawReader
	priorityCache []interface{}
	structCache   []interface{}
	handlers      map[string]ReadHandler
}

type markerObject struct{}

var underConstruction = markerObject{}

// NewReader creates a new Reader.
func NewReader(r io.Reader, handlers map[string]ReadHandler) *Reader {
	return &Reader{newRawReader(r), make([]interface{}, 0, 32), make([]interface{}, 0, 16), handlers}
}

func (r *Reader) err() error {
	return r.raw.err
}

// ReadObject reads the next object from the Reader.
func (r *Reader) ReadObject() (interface{}, error) {
	return r.readObject(), r.err()
}

func (r *Reader) readNextCode() byte {
	return r.raw.readRawByte()
}

func (r *Reader) readInt() int {
	var result int

	code := r.readNextCode()
	switch code {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F:
		result = int(code) & 0xFF

	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F:
		result = ((int(code) - INT_PACKED_2_ZERO) << 8) | r.raw.readRawInt8()

	case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F:
		result = ((int(code) - INT_PACKED_3_ZERO) << 16) | r.raw.readRawInt16()

	case 0x70, 0x71, 0x72, 0x73:
		result = ((int(code) - INT_PACKED_4_ZERO) << 24) | r.raw.readRawInt24()

	case 0x74, 0x75, 0x76, 0x77:
		result = ((int(code) - INT_PACKED_5_ZERO) << 32) | r.raw.readRawInt32()

	case 0x78, 0x79, 0x7A, 0x7B:
		result = ((int(code) - INT_PACKED_6_ZERO) << 40) | r.raw.readRawInt40()

	case 0x7C, 0x7D, 0x7E, 0x7F:
		result = ((int(code) - INT_PACKED_7_ZERO) << 48) | r.raw.readRawInt48()

	case INT:
		result = r.raw.readRawInt64()

	default:
		obj := r.read(code)
		i, ok := obj.(int)
		if ok {
			return i
		} else {
			log.Fatalf("not an int: 0x%x, %#v\n", code, obj)
		}
	}

	return result
}

func (r *Reader) readObject() interface{} {
	return r.read(r.readNextCode())
}

func (r *Reader) read(code byte) interface{} {
	var result interface{}

	switch code {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F:
		result = int(code) & 0xFF

	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F:
		result = ((int(code) - INT_PACKED_2_ZERO) << 8) | r.raw.readRawInt8()

	case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F:
		result = ((int(code) - INT_PACKED_3_ZERO) << 16) | r.raw.readRawInt16()

	case 0x70, 0x71, 0x72, 0x73:
		result = ((int(code) - INT_PACKED_4_ZERO) << 24) | r.raw.readRawInt24()

	case 0x74, 0x75, 0x76, 0x77:
		result = ((int(code) - INT_PACKED_5_ZERO) << 32) | r.raw.readRawInt32()

	case 0x78, 0x79, 0x7A, 0x7B:
		result = ((int(code) - INT_PACKED_6_ZERO) << 40) | r.raw.readRawInt40()

	case 0x7C, 0x7D, 0x7E, 0x7F:
		result = ((int(code) - INT_PACKED_7_ZERO) << 48) | r.raw.readRawInt48()

	case PUT_PRIORITY_CACHE:
		idx := len(r.priorityCache)
		r.priorityCache = append(r.priorityCache, underConstruction)
		r.priorityCache[idx] = r.readObject()
		result = r.priorityCache[idx]

	case GET_PRIORITY_CACHE:
		result = lookupCache(r.priorityCache, r.readInt())

	case PRIORITY_CACHE_PACKED_START + 0, PRIORITY_CACHE_PACKED_START + 1,
		PRIORITY_CACHE_PACKED_START + 2, PRIORITY_CACHE_PACKED_START + 3,
		PRIORITY_CACHE_PACKED_START + 4, PRIORITY_CACHE_PACKED_START + 5,
		PRIORITY_CACHE_PACKED_START + 6, PRIORITY_CACHE_PACKED_START + 7,
		PRIORITY_CACHE_PACKED_START + 8, PRIORITY_CACHE_PACKED_START + 9,
		PRIORITY_CACHE_PACKED_START + 10, PRIORITY_CACHE_PACKED_START + 11,
		PRIORITY_CACHE_PACKED_START + 12, PRIORITY_CACHE_PACKED_START + 13,
		PRIORITY_CACHE_PACKED_START + 14, PRIORITY_CACHE_PACKED_START + 15,
		PRIORITY_CACHE_PACKED_START + 16, PRIORITY_CACHE_PACKED_START + 17,
		PRIORITY_CACHE_PACKED_START + 18, PRIORITY_CACHE_PACKED_START + 19,
		PRIORITY_CACHE_PACKED_START + 20, PRIORITY_CACHE_PACKED_START + 21,
		PRIORITY_CACHE_PACKED_START + 22, PRIORITY_CACHE_PACKED_START + 23,
		PRIORITY_CACHE_PACKED_START + 24, PRIORITY_CACHE_PACKED_START + 25,
		PRIORITY_CACHE_PACKED_START + 26, PRIORITY_CACHE_PACKED_START + 27,
		PRIORITY_CACHE_PACKED_START + 28, PRIORITY_CACHE_PACKED_START + 29,
		PRIORITY_CACHE_PACKED_START + 30, PRIORITY_CACHE_PACKED_START + 31:
		result = lookupCache(r.priorityCache, int(code-PRIORITY_CACHE_PACKED_START))

	case STRUCT_CACHE_PACKED_START + 0, STRUCT_CACHE_PACKED_START + 1,
		STRUCT_CACHE_PACKED_START + 2, STRUCT_CACHE_PACKED_START + 3,
		STRUCT_CACHE_PACKED_START + 4, STRUCT_CACHE_PACKED_START + 5,
		STRUCT_CACHE_PACKED_START + 6, STRUCT_CACHE_PACKED_START + 7,
		STRUCT_CACHE_PACKED_START + 8, STRUCT_CACHE_PACKED_START + 9,
		STRUCT_CACHE_PACKED_START + 10, STRUCT_CACHE_PACKED_START + 11,
		STRUCT_CACHE_PACKED_START + 12, STRUCT_CACHE_PACKED_START + 13,
		STRUCT_CACHE_PACKED_START + 14, STRUCT_CACHE_PACKED_START + 15:
		st := lookupCache(r.structCache, int(code-STRUCT_CACHE_PACKED_START)).(structType)
		result = r.handleStruct(st.tag, st.fields)

	case MAP:
		kvs := r.readObject().([]interface{})
		m := make(map[interface{}]interface{}, len(kvs)/2)
		for i := 0; i < len(kvs); i += 2 {
			m[kvs[i]] = kvs[i+1]
		}
		result = m

		// TODO: SET

	case UUID:
		result = r.handleStruct("uuid", 2)

		// TODO: REGEX

	case URI:
		result = r.handleStruct("uri", 1)

		// TODO: BIGINT, BIGDEC

	case INST:
		milliseconds := int64(r.readInt())
		result = time.Unix(milliseconds/1000, (milliseconds%1000)*10e6)

		// TODO: SYM

	case KEY:
		result = r.handleStruct("key", 2)

	case INT_ARRAY, LONG_ARRAY:
		length := r.readCount()
		nums := make([]int, length)
		for i := 0; i < length; i++ {
			nums[i] = r.readInt()
		}
		result = nums

	case FLOAT_ARRAY:
		length := r.readCount()
		floats := make([]float32, length)
		for i := 0; i < length; i++ {
			floats[i] = r.readObject().(float32)
		}
		result = floats

	case BOOLEAN_ARRAY:
		length := r.readCount()
		bools := make([]bool, length)
		for i := 0; i < length; i++ {
			bools[i] = r.readObject().(bool)
		}
		result = bools

	case DOUBLE_ARRAY:
		length := r.readCount()
		doubles := make([]float64, length)
		for i := 0; i < length; i++ {
			doubles[i] = r.readObject().(float64)
		}
		result = doubles

	case OBJECT_ARRAY:
		result = r.readObjects(r.readCount())

	case BYTES_PACKED_LENGTH_START + 0, BYTES_PACKED_LENGTH_START + 1,
		BYTES_PACKED_LENGTH_START + 2, BYTES_PACKED_LENGTH_START + 3,
		BYTES_PACKED_LENGTH_START + 4, BYTES_PACKED_LENGTH_START + 5,
		BYTES_PACKED_LENGTH_START + 6, BYTES_PACKED_LENGTH_START + 7:
		result = r.internalReadBytes(int(code - BYTES_PACKED_LENGTH_START))

	case BYTES:
		length := r.readCount()
		result = r.internalReadBytes(length)

	case BYTES_CHUNK:
		result = r.internalReadChunkedBytes()

	case STRING_PACKED_LENGTH_START + 0,
		STRING_PACKED_LENGTH_START + 1,
		STRING_PACKED_LENGTH_START + 2,
		STRING_PACKED_LENGTH_START + 3,
		STRING_PACKED_LENGTH_START + 4,
		STRING_PACKED_LENGTH_START + 5,
		STRING_PACKED_LENGTH_START + 6,
		STRING_PACKED_LENGTH_START + 7:
		result = r.internalReadString(int(code - STRING_PACKED_LENGTH_START))

	case STRING:
		result = r.internalReadString(r.readCount())

		// TODO: STRING_CHUNK

	case LIST_PACKED_LENGTH_START + 0,
		LIST_PACKED_LENGTH_START + 1,
		LIST_PACKED_LENGTH_START + 2,
		LIST_PACKED_LENGTH_START + 3,
		LIST_PACKED_LENGTH_START + 4,
		LIST_PACKED_LENGTH_START + 5,
		LIST_PACKED_LENGTH_START + 6,
		LIST_PACKED_LENGTH_START + 7:
		length := int(code - LIST_PACKED_LENGTH_START)
		list := make([]interface{}, length)
		for i := 0; i < length; i++ {
			list[i] = r.readObject()
		}
		result = list

	case LIST:
		length := r.readCount()
		result = r.readObjects(length)

	case BEGIN_CLOSED_LIST:
		result = r.readClosedList()

	case BEGIN_OPEN_LIST:
		result = r.readOpenList()

	case TRUE:
		result = true
	case FALSE:
		result = false

	case DOUBLE:
		bs := r.internalReadBytes(8)
		var double float64
		err := binary.Read(bytes.NewBuffer(bs), binary.BigEndian, &double)
		if err != nil {
			log.Fatal("invalid double")
		}
		result = double

	case DOUBLE_0:
		result = float64(0.0)

	case DOUBLE_1:
		result = float64(1.0)

	case FLOAT:
		result = -1.0

	case INT:
		result = r.raw.readRawInt64()

	case NULL:
		result = nil

		// TODO: FOOTER

	case STRUCTTYPE:
		tag := r.readObject().(string)
		fields := r.readInt()
		r.structCache = append(r.structCache, structType{tag, fields})
		result = r.handleStruct(tag, fields)

	case STRUCT:
		st := lookupCache(r.structCache, r.readInt()).(structType)
		result = r.handleStruct(st.tag, st.fields)

	case RESET_CACHES:
		r.priorityCache = make([]interface{}, 0, 32)
		r.structCache = make([]interface{}, 0, 16)
		result = r.readObject()

	default:
		log.Fatalf("not implemented or invalid: 0x%x\n", code)
	}

	return result
}

func (r *Reader) readCount() int {
	return r.readInt()
}

func (r *Reader) readObjects(length int) []interface{} {
	list := make([]interface{}, length)
	for i := 0; i < length; i++ {
		list[i] = r.readObject()
	}
	return list
}

func (r *Reader) internalReadBytes(length int) []byte {
	bs := make([]byte, length)
	for i := 0; i < length; i++ {
		bs[i] = r.raw.readRawByte()
	}
	return bs
}

func (r *Reader) internalReadChunkedBytes() []byte {
	bs := make([]byte, 0)
	code := byte(BYTES_CHUNK)
	for code == BYTES_CHUNK {
		bs = append(bs, r.internalReadBytes(r.readCount())...)
		code = r.readNextCode()
	}
	if code != BYTES {
		log.Fatal("invalid byte chunk")
	}
	bs = append(bs, r.internalReadBytes(r.readCount())...)
	return bs
}

func (r *Reader) internalReadString(length int) string {
	bs := make([]byte, length)
	for i := 0; i < length; i++ {
		bs[i] = r.raw.readRawByte()
	}
	return string(bs)
}

func (r *Reader) readClosedList() []interface{} {
	list := make([]interface{}, 0)
	for {
		code := r.readNextCode()
		if code == END_COLLECTION {
			return list
		}
		list = append(list, r.read(code))
	}
}

func (r *Reader) readOpenList() []interface{} {
	list := make([]interface{}, 0)
	for {
		code := r.readNextCode()
		if r.err() == io.EOF {
			code = END_COLLECTION
		}
		if code == END_COLLECTION {
			return list
		}
		list = append(list, r.read(code))
	}
}

func (r *Reader) handleStruct(key string, fieldCount int) interface{} {
	switch key {
	case "key":
		namespace := r.readObject()
		if namespace == nil {
			namespace = ""
		}
		name := r.readObject()
		return Key{
			Namespace: namespace.(string),
			Name:      name.(string),
		}

	case "uuid":
		obj := r.readObject()
		bs, ok := obj.([]byte)
		if !ok || len(bs) != 16 {
			log.Fatal("invalid uuid")
		}
		return fmt.Sprintf("%x-%x-%x-%x-%x", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:16])

	case "uri":
		rawURL := r.readObject().(string)
		u, err := url.Parse(rawURL)
		if err != nil {
			log.Fatal(err)
		}
		return u

	default:
		if handler, ok := r.handlers[key]; ok {
			return handler(r, key, fieldCount)
		}

		vals := r.readObjects(fieldCount)
		return StructAny{key, vals}
	}

	return nil
}

func lookupCache(cache []interface{}, idx int) interface{} {
	if idx < len(cache) {
		obj := cache[idx]
		if obj == underConstruction {
			log.Fatal("circular reference in cache!")
		} else {
			return obj
		}
	} else {
		log.Fatal("cache index out of range ", idx)
	}

	log.Fatal("unreachable")
	return nil
}
