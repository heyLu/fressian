package fressian

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
)

type rawWriter struct {
	bw    *bufio.Writer
	count int
	err   error
}

func newRawWriter(w io.Writer) *rawWriter {
	return &rawWriter{bufio.NewWriter(w), 0, nil}
}

func (w *rawWriter) writeRawByte(b byte) error {
	err := w.bw.WriteByte(b)
	if err != nil {
		w.err = err
		return err
	}

	w.count++
	return nil
}

func (w *rawWriter) writeRawInt16(i int) error {
	n, err := w.bw.Write([]byte{
		byte((i >> 8) & 0xff),
		byte(i & 0xff)})
	if err != nil {
		w.err = err
		return err
	}

	w.count += n
	return nil
}

func (w *rawWriter) writeRawInt24(i int) error {
	n, err := w.bw.Write([]byte{
		byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff),
		byte(i & 0xff),
	})
	if err != nil {
		w.err = err
		return err
	}

	w.count += n
	return nil
}

func (w *rawWriter) writeRawInt32(i int) error {
	n, err := w.bw.Write([]byte{
		byte((i >> 24) & 0xff),
		byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff),
		byte(i & 0xff),
	})
	if err != nil {
		w.err = err
		return err
	}

	w.count += n
	return nil
}

func (w *rawWriter) writeRawInt40(i int) error {
	n, err := w.bw.Write([]byte{
		byte((i >> 32) & 0xff),
		byte((i >> 24) & 0xff),
		byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff),
		byte(i & 0xff),
	})
	if err != nil {
		w.err = err
		return err
	}

	w.count += n
	return nil
}

func (w *rawWriter) writeRawInt48(i int) error {
	n, err := w.bw.Write([]byte{
		byte((i >> 40) & 0xff),
		byte((i >> 32) & 0xff),
		byte((i >> 24) & 0xff),
		byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff),
		byte(i & 0xff),
	})
	if err != nil {
		w.err = err
		return err
	}

	w.count += n
	return nil
}

func (w *rawWriter) writeRawInt64(i int) error {
	n, err := w.bw.Write([]byte{
		byte((i >> 56) & 0xff),
		byte((i >> 48) & 0xff),
		byte((i >> 40) & 0xff),
		byte((i >> 32) & 0xff),
		byte((i >> 24) & 0xff),
		byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff),
		byte(i & 0xff),
	})
	if err != nil {
		w.err = err
		return err
	}

	w.count += n
	return nil
}

func (w *rawWriter) writeRawFloat32(f float32) error {
	err := binary.Write(w.bw, binary.BigEndian, f)
	if err != nil {
		w.err = err
		return err
	}

	return nil
}

func (w *rawWriter) writeRawFloat64(f float64) error {
	err := binary.Write(w.bw, binary.BigEndian, f)
	if err != nil {
		w.err = err
		return err
	}

	return nil
}

func (w *rawWriter) writeRawBytes(bytes []byte, offset int, length int) error {
	n, err := w.bw.Write(bytes[offset : offset+length])
	if err != nil {
		w.err = err
		return err
	}

	w.count += n
	return nil
}

func (w *rawWriter) reset() {
	w.count = 0
	//w.err = nil
	// TODO: reset checksum
}

type Writer struct {
	raw              *rawWriter
	priorityCache    map[interface{}]int
	priorityCacheIdx int
	structCache      map[interface{}]int
	structCacheIdx   int
	handler          WriteHandler
}

func NewWriter(w io.Writer, handler WriteHandler) *Writer {
	if handler == nil {
		handler = DefaultHandler
	}
	return &Writer{
		newRawWriter(w),
		make(map[interface{}]int, 16), // TODO: fix these numbers
		0,
		make(map[interface{}]int, 16),
		0,
		handler,
	}
}

func (w *Writer) WriteNil() error {
	return w.writeCode(NULL)
}

func (w *Writer) writeCode(code int) error {
	return w.raw.writeRawByte(byte(code))
}

func (w *Writer) WriteBool(b bool) error {
	if b {
		return w.writeCode(TRUE)
	} else {
		return w.writeCode(FALSE)
	}
}

func (w *Writer) WriteInt(i int) error {
	return w.internalWriteInt(i)
}

func (w *Writer) WriteFloat32(f float32) error {
	return w.raw.writeRawFloat32(f)
}

func (w *Writer) WriteFloat64(f float64) error {
	return w.raw.writeRawFloat64(f)
}

func (w *Writer) WriteString(s string) error {
	log.Fatal("WriteString: not implemented")
	return nil
}

func (w *Writer) WriteList(l []interface{}) error {
	if l == nil {
		return w.WriteNil()
	}

	length := len(l)
	if length < LIST_PACKED_LENGTH_END {
		w.raw.writeRawByte(byte(LIST_PACKED_LENGTH_START + length))
	} else {
		w.writeCode(LIST)
		w.writeCount(length)
	}
	for _, o := range l {
		w.WriteObject(o)
	}
	return w.raw.err
}

func (w *Writer) WriteBytes(bytes []byte) error {
	if bytes == nil {
		return w.WriteNil()
	}

	return w.WriteBytes_(bytes, 0, len(bytes))
}

func (w *Writer) WriteBytes_(bytes []byte, offset int, length int) error {
	if length < BYTES_PACKED_LENGTH_END {
		w.raw.writeRawByte(byte(BYTES_PACKED_LENGTH_START + length))
		w.raw.writeRawBytes(bytes, offset, length)
	} else {
		for length > BYTE_CHUNK_SIZE {
			w.writeCode(BYTES_CHUNK)
			w.writeCount(BYTE_CHUNK_SIZE)
			w.raw.writeRawBytes(bytes, offset, BYTE_CHUNK_SIZE)
			offset += BYTE_CHUNK_SIZE
			length -= BYTE_CHUNK_SIZE
		}
		w.writeCode(BYTES)
		w.writeCount(length)
		w.raw.writeRawBytes(bytes, offset, length)
	}

	return w.raw.err
}

func (w *Writer) writeCount(c int) error {
	return w.WriteInt(c)
}

// c.f. java.lang.Long#numberOfLeadingZeros
func numberOfLeadingZeros(i int) int {
	if i == 0 {
		return 64
	}

	n := 1
	x := i >> 32
	if x == 0 {
		n += 32
		x = i
	}
	if x>>16 == 0 {
		n += 16
		x <<= 16
	}
	if x>>24 == 0 {
		n += 8
		x <<= 8
	}
	if x>>28 == 0 {
		n += 4
		x <<= 4
	}
	if x>>30 == 0 {
		n += 2
		x <<= 2
	}
	n -= x >> 31

	return n
}

func bitSwitch(i int) int {
	if i < 0 {
		i = ^i
	}

	return numberOfLeadingZeros(i)
}

func (w *Writer) internalWriteInt(i int) error {
	switch bitSwitch(i) {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14:
		w.writeCode(INT)
		return w.raw.writeRawInt64(i)
	case 15, 16, 17, 18, 19, 20, 21, 22:
		w.raw.writeRawByte(INT_PACKED_7_ZERO + byte(i>>48))
		return w.raw.writeRawInt48(i)
	case 23, 24, 25, 26, 27, 28, 29, 30:
		w.raw.writeRawByte(INT_PACKED_6_ZERO + byte(i>>40))
		return w.raw.writeRawInt40(i)
	case 31, 32, 33, 34, 35, 36, 37, 38:
		w.raw.writeRawByte(INT_PACKED_5_ZERO + byte(i>>32))
		return w.raw.writeRawInt32(i)
	case 39, 40, 41, 42, 43, 44:
		w.raw.writeRawByte(INT_PACKED_4_ZERO + byte(i>>24))
		return w.raw.writeRawInt24(i)
	case 45, 46, 47, 48, 49, 50, 51:
		w.raw.writeRawByte(INT_PACKED_3_ZERO + byte(i>>16))
		return w.raw.writeRawInt16(i)
	case 52, 53, 54, 55, 56, 57:
		w.raw.writeRawByte(INT_PACKED_2_ZERO + byte(i>>8))
		return w.raw.writeRawByte(byte(i))
	case 58, 59, 60, 61, 62, 63, 64:
		if i < -1 {
			w.raw.writeRawByte(INT_PACKED_2_ZERO + byte(i>>8))
		}
		return w.raw.writeRawByte(byte(i))
	default:
		log.Fatal("int too big: ", i)
		return nil
	}
}

func (w *Writer) clearCaches() {
	w.priorityCache = make(map[interface{}]int, 16)
	w.priorityCacheIdx = 0
	w.structCache = make(map[interface{}]int, 16)
	w.structCacheIdx = 0
}

func (w *Writer) ResetCaches() error {
	w.clearCaches()
	return w.writeCode(RESET_CACHES)
}

var tagToCode = map[interface{}]int{
/*"map":       MAP,
"set":       SET,
"uuid":      UUID,
"regex":     REGEX,
"uri":       URI,
"bigint":    BIGINT,
"bigdec":    BIGDEC,
"inst":      INST,
"sym":       SYM,
"key":       KEY,
"int[]":     INT_ARRAY,
"float[]":   FLOAT_ARRAY,
"double[]":  DOUBLE_ARRAY,
"long[]":    LONG_ARRAY,
"boolean[]": BOOLEAN_ARRAY,
"Object[]":  OBJECT_ARRAY,*/
}

func (w *Writer) writeTag(tag interface{}, componentCount int) error {
	shortcutCode, ok := tagToCode[tag]
	if ok {
		return w.writeCode(shortcutCode)
	} else {
		idx, ok := w.structCache[tag]
		if !ok {
			w.structCache[tag] = w.structCacheIdx
			w.structCacheIdx += 1
			w.writeCode(STRUCTTYPE)
			w.WriteObject(tag)
			return w.WriteInt(componentCount)
		} else if idx < STRUCT_CACHE_PACKED_END {
			return w.writeCode(STRUCT_CACHE_PACKED_START + idx)
		} else {
			w.writeCode(STRUCT)
			return w.WriteInt(idx)
		}
	}
}

func (w *Writer) WriteExt(tag interface{}, fields ...interface{}) error {
	w.writeTag(tag, len(fields))
	for _, field := range fields {
		w.WriteObject(field)
	}
	return w.raw.err
}

func shouldSkipCache(val interface{}) bool {
	if val == nil {
		return true
	}

	switch val := val.(type) {
	case bool:
		return true
	case int:
		switch bitSwitch(val) {
		case 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64:
			return true
		default:
			return false
		}
	case string:
		return len(val) == 0
	case float64:
		return val == 0.0 || val == 1.0
	default:
		return false
	}
}

func (w *Writer) doWrite(tag string, val interface{}, wh WriteHandler, cache bool) error {
	if cache {
		if shouldSkipCache(val) {
			return w.doWrite(tag, val, wh, false)
		} else {
			idx, ok := w.priorityCache[val]
			if !ok {
				w.priorityCache[tag] = w.priorityCacheIdx
				w.priorityCacheIdx += 1
				w.writeCode(PUT_PRIORITY_CACHE)
				return w.doWrite(tag, val, wh, false)
			} else if idx < PRIORITY_CACHE_PACKED_END {
				return w.writeCode(PRIORITY_CACHE_PACKED_START + idx)
			} else {
				w.writeCode(GET_PRIORITY_CACHE)
				return w.WriteInt(idx)
			}
		}
	} else {
		return wh(w, val)
	}
}

func (w *Writer) WriteAs(tag string, val interface{}, cache bool) error {
	// TODO: lookup handler
	return w.doWrite(tag, val, w.handler, cache)
}

// WriteAny or even Write?
func (w *Writer) WriteObject(val interface{}) error {
	// TODO: "" should be nil
	return w.WriteAs("", val, false)
}

func (w *Writer) BeginClosedList() error {
	return w.writeCode(BEGIN_CLOSED_LIST)
}

func (w *Writer) EndList() error {
	return w.writeCode(END_COLLECTION)
}

func (w *Writer) BeginOpenList() error {
	if w.raw.count != 0 {
		return errors.New("open list must be called from the top level, outside any footer context")
	}

	err := w.writeCode(BEGIN_OPEN_LIST)
	w.raw.reset()
	return err
}

func (w *Writer) Error() error {
	return w.raw.err
}

func (w *Writer) Flush() error {
	return w.raw.bw.Flush()
}
