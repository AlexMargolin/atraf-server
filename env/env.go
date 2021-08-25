package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	defaultStructTag = "env"
)

type Decoder struct {
	typ reflect.Type
	val reflect.Value
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

// Marshal receives a struct pointer and attempts to parse
// defined environment variables to match its types
//
// It's possible to pass a struct with default values.
// If defined, environment variables will overwrite the struct defaults.
//
// When needed additional type support, modify the convert func accordingly.
func (dec *Decoder) Marshal(s interface{}) error {
	dec.typ = reflect.TypeOf(s)
	dec.val = reflect.ValueOf(s)

	// Validate argument is a struct pointer
	if dec.val.Kind() != reflect.Ptr || dec.val.Elem().Kind() != reflect.Struct || dec.val.IsNil() {
		return errors.New("value must be a struct pointer")
	}

	err := dec.parse()
	if err != nil {
		return err
	}

	return nil
}

// parse Attempts to locate an environment variable by the struct field Tag name.
// When the environment variable is defined, attempt
// to convert it to its corresponding struct type.
func (dec *Decoder) parse() error {
	for i := 0; i < dec.val.Elem().NumField(); i++ {
		rt := dec.typ.Elem().Field(i)
		rv := dec.val.Elem().Field(i)

		if value := os.Getenv(rt.Tag.Get(defaultStructTag)); value != "" {
			if err := dec.convert(&rv, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// convert attempts to convert a given string value(v)
// to the reflected struct field type(*rv)
func (Decoder) convert(rv *reflect.Value, v string) error {
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
