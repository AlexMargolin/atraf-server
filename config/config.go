package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	defaultStructTag = "config"
)

// Marshal receives a struct pointer and attempts to parse
// defined environment variables to match its types
//
// It's possible to pass a struct with default values.
// If defined, environment variables will overwrite the struct defaults.
//
// When needed additional type support, modify the Parse func accordingly.
func Marshal(s interface{}) error {
	rt := reflect.TypeOf(s)
	rv := reflect.ValueOf(s)

	// Make the argument is a struct pointer
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct || rv.IsNil() {
		return errors.New("value must be a struct pointer")
	}

	for i := 0; i < rt.Elem().NumField(); i++ {
		rtf := rt.Elem().Field(i) // field value reflection
		rvf := rv.Elem().Field(i) // field type reflection

		tag := rtf.Tag.Get(defaultStructTag)
		if tag == "" {
			return errors.New(fmt.Sprintf("invalid struct config tag [%s]", rtf.Name))
		}

		if value, ok := os.LookupEnv(tag); ok {
			if err := convert(&rvf, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// convert attempts to convert a given string value(v)
// to the reflected struct field type
func convert(rv *reflect.Value, v string) error {
	switch rv.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil || rv.OverflowInt(value) {
			return err
		}
		rv.SetInt(value)
	case reflect.Bool:
		value, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		rv.SetBool(value)
	case reflect.String:
		rv.SetString(v)
	default:
		return errors.New(fmt.Sprintf("unsupported config type [%s]", rv.Type().Kind()))
	}

	return nil
}
