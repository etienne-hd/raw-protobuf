package protobuf

import (
	"fmt"
	"reflect"
)

func decode(b []byte, value any) error {
	v := reflect.ValueOf(value)
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

		var field *Field
		switch v.Field(i).Interface().(type) {
		case int32, int64:
			field = Find(b, nums...)
			if field != nil {
				v.Field(i).SetInt(int64(field.Int()))
			}
		case uint32, uint64:
			field = Find(b, nums...)
			if field != nil {
				v.Field(i).SetUint(uint64(field.Int()))
			}
		case bool:
			field = Find(b, nums...)
			if field != nil {
				v.Field(i).SetBool(field.Bool())
			}
		case string:
			field = Find(b, nums...)
			if field != nil {
				v.Field(i).SetString(field.String())
			}
		case float32:
			field = Find(b, nums...)
			if field != nil {
				v.Field(i).SetFloat(float64(field.Float32()))
			}
		case float64:
			field = Find(b, nums...)
			if field != nil {
				v.Field(i).SetFloat(field.Float64())
			}
		}

	}

	return nil
}

func Unmarshal(b []byte, value any) error {
	return decode(b, value)
}
