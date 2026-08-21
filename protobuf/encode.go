package protobuf

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
)

var (
	errInvalidType    = errors.New("invalid data type")
	errNumsEmpty      = errors.New("nums cannot be empty")
	errValueNotStruct = errors.New("value must be a struct")
)

type protobufFields map[int32]*protobufField

type protobufField struct {
	num      int32
	values   []any
	fields   protobufFields
	messages []protobufFields
}

func getNums(t reflect.StructField) ([]int32, error) {
	nums := []int32{}
	tag := t.Tag.Get("num")

	if tag == "" {
		return nil, nil
	}

	for _, rawNum := range strings.Split(tag, ".") {
		num, err := strconv.ParseInt(rawNum, 10, 32)
		if err != nil {
			return nil, err
		}

		nums = append(nums, int32(num))
	}

	return nums, nil
}

func getValues(value any) []any {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Slice {
		return []any{value}
	}

	values := make([]any, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		values = append(values, v.Index(i).Interface())
	}

	return values
}

func makeField(num int32, value any) (*protobufField, error) {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Struct {
		subFields, err := parse(value)
		if err != nil {
			return nil, err
		}

		return &protobufField{
			num:      num,
			messages: []protobufFields{subFields},
		}, nil
	}

	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Struct {
		field := &protobufField{
			num:      num,
			messages: make([]protobufFields, 0, v.Len()),
		}

		for i := 0; i < v.Len(); i++ {
			subFields, err := parse(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}

			field.messages = append(field.messages, subFields)
		}

		return field, nil
	}

	return &protobufField{
		num:    num,
		values: getValues(value),
	}, nil
}

func insertField(fields protobufFields, value any, nums ...int32) error {
	if len(nums) == 0 {
		return errNumsEmpty
	}

	num := nums[0]

	// Insert value
	if len(nums) == 1 {
		if _, ok := fields[num]; ok {
			return errors.New("values and fields cannot be used at the same time") // TODO: refactor
		}

		field, err := makeField(num, value)
		if err != nil {
			return err
		}

		fields[num] = field
		return nil
	}

	// Create field
	var nextField *protobufField
	field, ok := fields[num]
	if ok {
		if field.values != nil {
			return errors.New("values and fields cannot be used at the same time") // TODO: refactor
		}
		if field.fields == nil {
			return errors.New("messages and fields cannot be used at the same time") // TODO: refactor
		}
		nextField = field
	}

	if nextField == nil {
		nextField = &protobufField{
			num:    num,
			fields: make(protobufFields),
		}

		fields[num] = nextField
	}

	return insertField(nextField.fields, value, nums[1:]...)
}

func parse(value any) (protobufFields, error) {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Struct {
		return nil, errValueNotStruct
	}

	fields := make(protobufFields)

	for i := 0; i < t.NumField(); i++ {
		nums, err := getNums(t.Field(i))
		if err != nil {
			return nil, err
		}

		// No tag
		if nums == nil {
			continue
		}

		if err := insertField(fields, v.Field(i).Interface(), nums...); err != nil {
			return nil, err
		}
	}

	return fields, nil
}

func encodeValue(num int32, value any) ([]byte, error) {
	var b []byte

	switch value := value.(type) {
	case int32, int64, uint32, uint64, bool:
		b = protowire.AppendTag(b, protowire.Number(num), protowire.VarintType)

		switch value := value.(type) {
		case int32:
			b = protowire.AppendVarint(b, uint64(value))
		case int64:
			b = protowire.AppendVarint(b, uint64(value))
		case uint32:
			b = protowire.AppendVarint(b, uint64(value))
		case uint64:
			b = protowire.AppendVarint(b, uint64(value))
		case bool:
			if value {
				b = protowire.AppendVarint(b, 1)
			} else {
				b = protowire.AppendVarint(b, 0)
			}
		}

	case string:
		b = protowire.AppendTag(b, protowire.Number(num), protowire.BytesType)
		b = protowire.AppendString(b, value)

	case float32:
		b = protowire.AppendTag(b, protowire.Number(num), protowire.Fixed32Type)
		b = protowire.AppendFixed32(b, math.Float32bits(value))

	case float64:
		b = protowire.AppendTag(b, protowire.Number(num), protowire.Fixed64Type)
		b = protowire.AppendFixed64(b, math.Float64bits(value))

	default:
		return nil, errInvalidType
	}

	return b, nil
}

func encodeField(field *protobufField) ([]byte, error) {
	var b []byte

	if field.values != nil {
		for _, value := range field.values {
			encodedValue, err := encodeValue(field.num, value)
			if err != nil {
				return nil, err
			}

			b = append(b, encodedValue...)
		}
		return b, nil
	}

	if field.fields != nil {
		for _, subField := range field.fields {
			b = protowire.AppendTag(b, protowire.Number(field.num), protowire.BytesType)

			rawField, err := encodeField(subField)
			if err != nil {
				return nil, err
			}

			b = protowire.AppendBytes(b, rawField)
		}
		return b, nil
	}

	for _, message := range field.messages {
		b = protowire.AppendTag(b, protowire.Number(field.num), protowire.BytesType)

		rawMessage, err := encode(message)
		if err != nil {
			return nil, err
		}

		b = protowire.AppendBytes(b, rawMessage)
	}

	return b, nil
}

func encode(fields protobufFields) ([]byte, error) {
	b := []byte{}

	for _, field := range fields {
		encodedField, err := encodeField(field)
		if err != nil {
			return nil, err
		}

		b = append(b, encodedField...)
	}

	return b, nil
}

func Marshal(value any) ([]byte, error) {
	fields, err := parse(value)
	if err != nil {
		return nil, err
	}

	return encode(fields)
}
