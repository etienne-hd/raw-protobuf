package protobuf

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/encoding/protowire"
)

func Find(b []byte, nums ...int32) *Field {
	fields := FindN(b, 1, nums...)

	if len(fields) == 0 {
		return nil
	}

	return &fields[0]
}

func FindAll(b []byte, nums ...int32) []Field {
	return FindN(b, -1, nums...)
}

func FindN(b []byte, limit int, nums ...int32) []Field {
	var fields []Field

	if len(nums) == 0 {
		return fields
	}

	for len(b) != 0 && (limit == -1 || len(fields) < int(limit)) {
		num, wireType, n := protowire.ConsumeTag(b)
		b = b[n:]

		var value any
		switch wireType {
		case protowire.VarintType:
			value, n = protowire.ConsumeVarint(b)
			b = b[n:]
		case protowire.Fixed32Type:
			value, n = protowire.ConsumeFixed32(b)
			b = b[n:]
		case protowire.Fixed64Type:
			value, n = protowire.ConsumeFixed64(b)
			b = b[n:]
		case protowire.BytesType:
			value, n = protowire.ConsumeBytes(b)

			if num == protowire.Number(nums[0]) && len(nums) > 1 {
				fields = append(fields, FindN(b[:n], limit-len(fields), nums[1:]...)...)
			}

			b = b[n:]
		}

		// deepest level of recursion
		if num == protowire.Number(nums[0]) && len(nums) == 1 {
			fields = append(fields, Field{
				wireType: wireType,
				value:    value,
			})
		}
	}
	return fields
}

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
