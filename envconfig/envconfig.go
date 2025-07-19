package envconfig

import (
	"encoding"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Read reads the environment variables using lookupEnv function and assigns values to the Holder struct.
// If lookupEnv is not provided, os.LookupEnv is used.
// There are the following struct tags:
// - env:"-"		- skips field
// - env:"NAME"		- sets environment name to look for
// - prefix:"NAME"	- sets prefix for embedded structs
// Types can implement the following interfaces to support custom values:
// - encoding.TextUnmarshaler
// - encoding.BinaryUnmarshaler
// - json.Unmarshaler
func Read[T any](holder *T, lookupEnv func(string) (string, bool)) error {
	if lookupEnv == nil {
		lookupEnv = os.LookupEnv
	}

	return read(lookupEnv, "", holder)
}

func read(le func(string) (string, bool), prefix string, holder any) error {
	if i, ok := holder.(encoding.TextUnmarshaler); ok {
		value, _ := le(prefix[:len(prefix)-1])
		if value == "" {
			return nil
		}

		if err := i.UnmarshalText([]byte(value)); err != nil {
			return fmt.Errorf("envconfig: %w", err)
		}
	}
	if i, ok := holder.(encoding.BinaryUnmarshaler); ok {
		value, _ := le(prefix[:len(prefix)-1])
		if value == "" {
			return nil
		}

		if err := i.UnmarshalBinary([]byte(value)); err != nil {
			return fmt.Errorf("envconfig: %w", err)
		}
	}
	if i, ok := holder.(json.Unmarshaler); ok {
		value, _ := le(prefix[:len(prefix)-1])
		if value == "" {
			return nil
		}

		if err := i.UnmarshalJSON([]byte(value)); err != nil {
			return fmt.Errorf("envconfig: %w", err)
		}
	}

	tp := reflect.TypeOf(holder)
	if tp.Kind() != reflect.Ptr {
		panic("envconfig.Read only accepts a pointer, got " + tp.Kind().String())
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		panic("envconfig.Read only accepts a struct, got " + tp.Kind().String())
	}

	retValue := reflect.ValueOf(holder).Elem()
	fields := reflect.VisibleFields(tp)

	for _, field := range fields {
		pref, hasPrefix := field.Tag.Lookup("prefix")
		env, hasEnv := field.Tag.Lookup("env")
		if env == "-" || (!hasEnv && !hasPrefix) {
			continue
		}
		if env == "" && !hasPrefix {
			return fmt.Errorf("envconfig: tag \"env\" can't be empty: %q", field.Name)
		}

		retField := retValue.FieldByName(field.Name)

		if field.Type.Kind() == reflect.Struct {
			if !hasPrefix || pref == "" {
				return fmt.Errorf("envconfig: embedded structs must contain \"prefix\" tag: %q does not have it", field.Name)
			}

			err := read(le, prefix+pref+"_", retField.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		envValue, ok := le(prefix + env)
		if envValue == "" {
			if field.Tag.Get("required") == "true" && !ok {
				return fmt.Errorf("envconfig: required field %q is empty", prefix+env)
			}
			defaultValue := field.Tag.Get("default")
			if defaultValue != "" {
				envValue = defaultValue
			}
			if envValue == "" {
				continue
			}
		}

		if err := setValue(retField, envValue); err != nil {
			return err
		}
	}

	return nil
}

var durationType = reflect.TypeOf(time.Duration(0))

func setValue(inp reflect.Value, value string) error {
	switch inp.Type() {
	case durationType:
		d, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("envconfig: %w", err)
		}
		inp.Set(reflect.ValueOf(d))
		return nil
	}

	switch inp.Kind() {
	case reflect.String:
		inp.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("envconfig: %w", err)
		}
		inp.SetBool(b)
	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("envconfig: %w", err)
		}
		inp.SetInt(int64(i))
	case reflect.Int8:
		return parseInt(inp, value, 10, 8)
	case reflect.Int16:
		return parseInt(inp, value, 10, 16)
	case reflect.Int32:
		return parseInt(inp, value, 10, 32)
	case reflect.Int64:
		return parseInt(inp, value, 10, 64)
	case reflect.Uint:
		return parseUint(inp, value, 10, 0)
	case reflect.Uint8:
		return parseUint(inp, value, 10, 8)
	case reflect.Uint16:
		return parseUint(inp, value, 10, 16)
	case reflect.Uint32:
		return parseUint(inp, value, 10, 32)
	case reflect.Uint64:
		return parseUint(inp, value, 10, 64)
	case reflect.Float32:
		return parseFloat(inp, value, 32)
	case reflect.Float64:
		return parseFloat(inp, value, 64)
	case reflect.Array:
		arr := strings.Split(value, ",")
		for i := 0; i < inp.Len(); i++ {
			err := setValue(inp.Index(i), arr[i])
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		arr := strings.Split(value, ",")
		for i := 0; i < len(arr); i++ {
			elem := reflect.New(inp.Type().Elem()).Elem()
			err := setValue(elem, arr[i])
			if err != nil {
				return err
			}
			inp.Set(reflect.Append(inp, elem))
		}
	case reflect.Map:
		mp := reflect.MakeMap(inp.Type())
		arr := strings.Split(value, ",")
		for i := 0; i < len(arr); i++ {
			kv := strings.Split(arr[i], "=")
			if len(kv) != 2 {
				return fmt.Errorf("envconfig: invalid map value %s", value)
			}
			key := reflect.New(inp.Type().Key()).Elem()
			err := setValue(key, kv[0])
			if err != nil {
				return err
			}
			val := reflect.New(inp.Type().Elem()).Elem()
			err = setValue(val, kv[1])
			if err != nil {
				return err
			}
			mp.SetMapIndex(key, val)
		}
		inp.Set(mp)
	default:
		panic("unhandled default case")
	}

	return nil
}

func parseFloat(inp reflect.Value, value string, bitSize int) error {
	f, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return fmt.Errorf("envconfig: %w", err)
	}
	inp.SetFloat(f)
	return nil
}

func parseUint(inp reflect.Value, value string, base int, bitSize int) error {
	i, err := strconv.ParseUint(value, base, bitSize)
	if err != nil {
		return fmt.Errorf("envconfig: %w", err)
	}
	inp.SetUint(i)
	return nil
}

func parseInt(inp reflect.Value, value string, base, bitSize int) error {
	i, err := strconv.ParseInt(value, base, bitSize)
	if err != nil {
		return fmt.Errorf("envconfig: %w", err)
	}
	inp.SetInt(i)
	return nil
}
