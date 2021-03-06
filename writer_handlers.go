package fressian

import (
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"time"
)

type WriteHandler func(w *Writer, val interface{}) error

type ConversionError error

func IsConversionError(e error) bool {
	_, ok := e.(ConversionError)
	return ok
}

func DefaultHandler(w *Writer, val interface{}) error {
	if val == nil {
		return w.WriteNil()
	}

	switch val := val.(type) {
	case bool:
		return w.WriteBool(val)
	case int:
		return w.WriteInt(val)
	case float32:
		return w.WriteFloat32(val)
	case float64:
		return w.WriteFloat64(val)
	case string:
		return w.WriteString(val)
	case Keyword:
		w.writeCode(KEY)
		w.WriteValue(val.Namespace)
		return w.WriteValue(val.Name)
	case UUID:
		w.writeCode(CODE_UUID)
		return w.WriteBytes(val.Bytes())
	case time.Time:
		w.writeCode(INST)
		return w.WriteInt(int(val.Unix() * 1000))
	case *url.URL:
		w.writeCode(URI)
		return w.WriteString(val.String())
	case big.Int:
		w.writeCode(BIGINT)
		bs := val.Bytes()
		if val.Sign() < 0 {
			for i, _ := range bs {
				bs[i] ^= 0xff
			}
			bs[len(bs)-1] -= 1
		}
		return w.WriteBytes_(bs, 0, len(bs))
	case []bool:
		w.writeCode(BOOLEAN_ARRAY)
		w.writeCount(len(val))
		for _, b := range val {
			w.WriteBool(b)
		}
		return w.Error()
	case []int:
		w.writeCode(INT_ARRAY)
		w.writeCount(len(val))
		for _, i := range val {
			w.WriteInt(i)
		}
		return w.Error()
	case []byte:
		return w.WriteBytes_(val, 0, len(val))
	case []interface{}:
		return w.WriteList(val)
	default:
		switch reflect.TypeOf(val).Kind() {
		case reflect.Slice:
			// TODO: don't copy, write directly
			val := reflect.ValueOf(val)
			vals := make([]interface{}, val.Len())
			for i := 0; i < val.Len(); i++ {
				vals[i] = val.Index(i).Interface()
			}
			return w.WriteList(vals)
		case reflect.Map:
			val := reflect.ValueOf(val)
			ks := val.MapKeys()
			kvs := make([]interface{}, len(ks)*2)
			for i, k := range ks {
				kvs[i*2] = k.Interface()
				kvs[i*2+1] = val.MapIndex(k).Interface()
			}
			w.writeCode(MAP)
			return w.WriteList(kvs)
		default:
			return ConversionError(errors.New(fmt.Sprintf("don't know how to convert '%s'", val)))
		}
	}
}
