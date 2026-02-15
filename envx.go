// Package envx provides a simple way to load environment variables into a struct.
package envx

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	// DotEnvPath is an optional development-only fallback file.
	// Runtime environment variables always take precedence.
	DotEnvPath string
}

func Load(cfg any, opts Options) error {
	var dotEnv map[string]string

	if opts.DotEnvPath != "" {
		parsed, err := parseDotEnv(opts.DotEnvPath)
		if err != nil {
			return err
		}
		dotEnv = parsed
	}

	v := reflect.ValueOf(cfg)
	if !v.IsValid() || v.Kind() != reflect.Pointer || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return errors.New("cfg must be pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if fieldType.PkgPath != "" {
			continue // unexported field
		}

		key := fieldType.Tag.Get("env")
		if key == "" {
			continue
		}

		value, ok := resolveValue(key, fieldType, dotEnv)

		required := strings.EqualFold(fieldType.Tag.Get("required"), "true")
		allowEmpty := strings.EqualFold(fieldType.Tag.Get("allowempty"), "true")

		if !ok && required {
			return fmt.Errorf("missing required env: %s", key)
		}
		if required && value == "" && !allowEmpty {
			return fmt.Errorf("required env is empty: %s", key)
		}

		if !ok {
			continue
		}

		if err := setField(field, value); err != nil {
			return fmt.Errorf("invalid value for %s: %w", key, err)
		}
	}

	return nil
}

func resolveValue(key string, fieldType reflect.StructField, dotEnv map[string]string) (string, bool) {
	if value, ok := os.LookupEnv(key); ok {
		return value, true
	}

	if dotEnv != nil {
		if value, ok := dotEnv[key]; ok {
			return value, true
		}
	}

	if value, ok := fieldType.Tag.Lookup("default"); ok {
		return value, true
	}

	return "", false
}

func setField(field reflect.Value, value string) error {
	if !field.CanSet() {
		return nil
	}

	// Handle pointer types by initializing them
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	// Handle time.Duration
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(d))
		return nil
	}

	// Handle String Slices
	if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.String {
		parts := strings.Split(value, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		field.Set(reflect.ValueOf(parts))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}

	return nil
}
