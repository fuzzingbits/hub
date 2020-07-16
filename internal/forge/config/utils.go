package config

import (
	"errors"
	"reflect"
	"strconv"
)

var (
	// ErrInvalidValue returned when the value passed to Unmarshal is nil or not a
	// pointer to a struct.
	ErrInvalidValue = errors.New("value must be a non-nil pointer to a struct")

	// ErrUnsupportedType returned when a field with tag "env" is unsupported.
	ErrUnsupportedType = errors.New("field is an unsupported type")

	// ErrUnexportedField returned when a field with tag "env" is not exported.
	ErrUnexportedField = errors.New("field must be exported")
)

func set(t reflect.Type, f reflect.Value, value string) error {
	switch t.Kind() {
	case reflect.Ptr:
		ptr := reflect.New(t.Elem())
		err := set(t.Elem(), ptr.Elem(), value)
		if err != nil {
			return err
		}
		f.Set(ptr)
	case reflect.String:
		f.SetString(value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		f.SetBool(v)
	case reflect.Int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		f.SetInt(int64(v))
	default:
		return ErrUnsupportedType
	}

	return nil
}
