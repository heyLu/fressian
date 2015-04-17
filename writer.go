package fressian

import (
	"bufio"
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
		return err
	}

	w.count += n
	return nil
}

func (w *rawWriter) writeRawBytes(bytes []byte, offset int, length int) error {
	n, err := w.bw.Write(bytes[offset : offset+length])
	if err != nil {
		return err
	}

	w.count += n
	return nil
}

type Writer struct {
	raw *rawWriter
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

// WriteFloat (WriteFloat32?)
// WriteDouble (WriteFloat64?)

// WriteString

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

func bitSwitch(i int) int {
	log.Fatal("not implemented")
	return 0
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

// WriteAs

// WriteAny or even Write?
func (w *Writer) WriteObject(o interface{}) error {
	log.Fatal("not implemented")
	return nil
}
