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

// Marshal receives a struct pointer and attempts to parse
// defined environment variables to match its types
//
// It's possible to pass a struct with default values.
// If defined, environment variables will overwrite the struct defaults.
//
// When needed additional type support, modify the convert func accordingly.
func (d *Decoder) Marshal(s interface{}) error {
	d.typ = reflect.TypeOf(s)
	d.val = reflect.ValueOf(s)

	// Validate argument is a struct pointer
	if d.val.Kind() != reflect.Ptr || d.val.Elem().Kind() != reflect.Struct || d.val.IsNil() {
		return errors.New("value must be a struct pointer")
	}

	err := d.parse()
	if err != nil {
		return err
	}

	return nil
}

// parse Attempts to locate an environment variable by the struct field Tag name.
// When the environment variable is defined, attempt
// to convert it to its corresponding struct type.
func (d *Decoder) parse() error {
	for i := 0; i < d.val.Elem().NumField(); i++ {
		rt := d.typ.Elem().Field(i)
		rv := d.val.Elem().Field(i)

		if value := os.Getenv(rt.Tag.Get(defaultStructTag)); value != "" {
			if err := d.convert(&rv, value); err != nil {
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

func NewDecoder() *Decoder {
	return &Decoder{}
}
