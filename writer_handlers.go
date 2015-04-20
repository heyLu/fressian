package fressian

import (
	"errors"
	"fmt"
	"reflect"
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
				kvs[i] = k
				kvs[i+1] = val.MapIndex(k)
			}
			w.writeCode(MAP)
			return w.WriteList(kvs)
		}
	}

	return ConversionError(errors.New(fmt.Sprintf("don't know how to convert '%s'", val)))
}
