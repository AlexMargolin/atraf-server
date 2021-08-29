package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	defaultConfFile = "conf.json"
	defaultEnvTag   = "env"
)

type Decoder struct{}

// Decode receives struct pointer which represents our config structure.
// Assignment Order (High Priority -> Low Priority):
// 1. Environment Variable -> Configuration File -> Struct Defaults
func (decoder *Decoder) Decode(v interface{}) error {
	decoder.ApplyJson(v)

	return decoder.unmarshal(v)
}

// ApplyJson attempts to unmarshal the json object onto the config struct
func (decoder *Decoder) ApplyJson(v interface{}) {
	file, err := os.Open(defaultConfFile)
	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(v)
	if err != nil {
		return
	}
}

func (decoder *Decoder) unmarshal(v interface{}) error {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	// Validate argument is a struct pointer
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct || val.IsNil() {
		return errors.New("value must be a struct pointer")
	}

	for i := 0; i < val.Elem().NumField(); i++ {
		rt := typ.Elem().Field(i) // Field Reflect Type
		rv := val.Elem().Field(i) // Field Reflect Value

		tag := rt.Tag.Get(defaultEnvTag) // Field Reflect Tag
		if value := os.Getenv(tag); value != "" {
			if err := decoder.covert(value, &rv); err != nil {
				return err
			}
		}
	}

	return nil
}

func (decoder *Decoder) parseLine() {

}

// Convert attempts to convert a given string value(v)
// to the reflected struct field type(*rv)
func (Decoder) covert(s string, rv *reflect.Value) error {
	switch rv.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(s, 10, 64)
		if err != nil || rv.OverflowInt(value) {
			return err
		}
		rv.SetInt(value)

	case reflect.Bool:
		value, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		rv.SetBool(value)

	case reflect.String:
		rv.SetString(s)

	default:
		return errors.New(fmt.Sprintf("unsupported config type [%s]", rv.Type().Kind()))
	}

	return nil
}

func NewConfig() *Decoder {
	return &Decoder{}
}
