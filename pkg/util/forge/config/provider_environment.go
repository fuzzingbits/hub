package config

import (
	"os"
	"reflect"
)

// ProviderEnvironment is foobar
type ProviderEnvironment struct{}

// Unmarshal is foobar
func (p ProviderEnvironment) Unmarshal(target interface{}) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidValue
	}

	rv = rv.Elem()

	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		valueField := rv.Field(i)
		switch valueField.Kind() {
		case reflect.Struct:
			iface := valueField.Addr().Interface()
			err := p.Unmarshal(iface)
			if err != nil {
				return err
			}
		}

		typeField := t.Field(i)
		tag := typeField.Tag.Get("env")
		if tag == "" {
			continue
		}

		if !valueField.CanSet() {
			return ErrUnexportedField
		}

		envVar, ok := os.LookupEnv(tag)
		if !ok {
			continue
		}

		err := set(typeField.Type, valueField, envVar)
		if err != nil {
			return err
		}
	}

	return nil
}
