package protobuf

import (
	"fmt"
	"reflect"
)

func setFieldValue(v reflect.Value, field *Field) error {
	switch v.Kind() {
	case reflect.Int32, reflect.Int64:
		v.SetInt(int64(field.Int()))

	case reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(field.Int()))

	case reflect.Bool:
		v.SetBool(field.Bool())

	case reflect.String:
		v.SetString(field.String())

	case reflect.Float32:
		v.SetFloat(float64(field.Float32()))

	case reflect.Float64:
		v.SetFloat(field.Float64())

	case reflect.Struct:
		decode(field.Bytes(), v.Addr())

	default:
		return fmt.Errorf("Invalid data type: %s", v.Kind())
	}
	return nil
}

func decodeField(b []byte, v reflect.Value, nums ...int32) error {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Slice {
		var fields []Field
		fields = FindAll(b, nums...)
		for _, field := range fields {
			x := reflect.New(v.Type().Elem()).Elem()
			if err := setFieldValue(x, &field); err != nil {
				return err
			}
			v.Set(reflect.Append(v, x))
		}
		return nil
	}

	field := Find(b, nums...)
	if field == nil {
		return nil
	}

	return setFieldValue(v, field)
}

func decodeFields(b []byte, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		nums, err := getNums(t.Field(i))
		if err != nil {
			return err
		}

		// No tag
		if nums == nil {
			continue
		}

		if err := decodeField(b, v.Field(i), nums...); err != nil {
			return err
		}
	}

	return nil
}

func decode(b []byte, v reflect.Value) error {
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("decode: value must be a non-nil pointer")
	}

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("decode: expected struct, got %s", v.Kind())
	}

	return decodeFields(b, v)
}

func Unmarshal(b []byte, value any) error {
	return decode(b, reflect.ValueOf(value))
}
