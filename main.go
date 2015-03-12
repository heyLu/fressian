package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// from org.fressian.impl.Codes
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
	MAP                         = 0xC0 // so there *is* special support for maps?
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

type RawReader struct {
	br    *bufio.Reader
	count int
	err   error
}

func newRawReader(r io.Reader) *RawReader {
	return &RawReader{bufio.NewReader(r), 0, nil}
}

func (r *RawReader) Err() error {
	return r.err
}

func (r *RawReader) readRawByte() byte {
	res, err := r.br.ReadByte()
	if err != nil {
		r.err = err
		return 0
	}
	r.count++
	return res
}

func (r *RawReader) readRawInt8() int {
	return int(r.readRawByte())
}

func (r *RawReader) readRawInt16() int {
	return int((r.readRawByte() << 8) + r.readRawByte())
}

func (r *RawReader) readRawInt24() int {
	return int((r.readRawByte() << 16) + (r.readRawByte() << 8) + r.readRawByte())
}

type Reader struct {
	raw *RawReader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{newRawReader(r)}
}

func (r *Reader) Err() error {
	return r.raw.Err()
}

func (r *Reader) readNextCode() byte {
	return r.raw.readRawByte()
}

func (r *Reader) readInt() int {
	var result int

	code := r.readNextCode()
	switch code {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09:
		result = int(code) & 0xff

	default:
		log.Fatal("not implemented")
	}

	return result
}

func (r *Reader) readObject() interface{} {
	return r.read(r.readNextCode())
}

func (r *Reader) read(code byte) interface{} {
	var result interface{}

	switch code {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09:
		// TODO: 0x1A - 0x7F
		result = int(code) & 0xFF

		// TODO: {GET,PUT}_PRIORITY_CACHE, PRIORITY_CACHE_PACKED_START + {0..31}
		// TODO: STRUCT_CACHE_PACKED_START + {0..15}
		// TODO: MAP, SET, UUID, REGEX, URI, BIGINT, BIGDEC, INST, SYM, KEY
		// TODO: {INT,LONG,FLOAT,BOOLEAN,DOUBLE,OBJECT}_ARRAY
		// TODO: BYTES_PACKED_LENGTH_START + {0..7}, BYTES, BYTES_CHUNK
		// TODO: STRING_PACKED_LENGTH_START + {0..7}, STRING, STRING_CHUNK
		// TODO: LIST_PACKED_LENGTH_START + {0..7}, LIST, BEGIN_{CLOSED,OPEN}_LIST

	case TRUE:
		result = true
	case FALSE:
		result = false

		// TODO: DOUBLE, DOUBLE_0, DOUBLE_1, FLOAT, INT

	case NULL:
		result = nil

		// TODO: FOOTER
		// TODO: STRUCTTYPE, STRUCT
		// TODO: RESET_CACHES

	default:
		log.Fatal("not implemented or invalid")
	}

	return result
}

func main() {
	f, err := os.Open("example.fressian")
	if err != nil {
		log.Fatal(err)
	}

	r := NewReader(f)
	obj := r.readObject()
	if r.Err() != nil {
		log.Fatal(r.Err())
	}
	fmt.Printf("%#v\n", obj)
}
