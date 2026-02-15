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
	DotEnvPath string
}

func Load(cfg any, opts Options) error {
	if opts.DotEnvPath != "" {
		if err := loadDotEnv(opts.DotEnvPath); err != nil {
			return err
		}
	}

	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
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

		value := os.Getenv(key)

		if value == "" {
			value = fieldType.Tag.Get("default")
		}

		required := strings.EqualFold(fieldType.Tag.Get("required"), "true")

		if value == "" && required {
			return fmt.Errorf("missing required env: %s", key)
		}

		if value == "" {
			continue
		}

		if err := setField(field, value); err != nil {
			return fmt.Errorf("invalid value for %s: %w", key, err)
		}
	}

	return nil
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
