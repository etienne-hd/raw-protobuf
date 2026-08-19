package protobuf

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
)

type protobufFields map[int32]*protobufField

type protobufField struct {
	num       int32
	value     any
	subFields protobufFields
}

func putField(fields protobufFields, value any, nums ...int32) error {
	if len(nums) == 0 {
		return errors.New("nums cannot be empty")
	}

	num := nums[0]
	field := fields[num]

	// Last level: insert the value
	if len(nums) == 1 {
		if field != nil {
			return errors.New("a field already exist")
		}

		fields[num] = &protobufField{
			num:   num,
			value: value,
		}
		return nil
	}

	// Intermediate level: create the parent if necessary
	if field == nil {
		field = &protobufField{
			num:       num,
			subFields: make(protobufFields),
		}

		fields[num] = field
	}

	// A field cannot contain both a value and sub-fields
	if field.subFields == nil {
		return errors.New("a field already exist")
	}

	return putField(field.subFields, value, nums[1:]...)
}

func getNums(t reflect.StructField) ([]int32, error) {
	rawNums := t.Tag.Get("nums")
	if rawNums == "" {
		return nil, nil
	}

	var nums []int32
	for _, rawNum := range strings.Split(rawNums, ".") {
		num, err := strconv.ParseInt(rawNum, 10, 64)
		if err != nil {
			return nil, err
		}

		nums = append(nums, int32(num))
	}

	return nums, nil
}

func parseFields(fields protobufFields, value any, nums ...int32) error {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	for i := 0; i < t.NumField(); i++ {
		fieldNums, err := getNums(t.Field(i))
		if err != nil {
			return err
		}
		// tag "nums" not found (skip)
		if fieldNums == nil {
			continue
		}

		// append nums + fieldNums
		currentNums := append([]int32{}, nums...)
		currentNums = append(currentNums, fieldNums...)

		// sub struct implementation
		if t.Field(i).Type.Kind() == reflect.Struct {
			if err := parseFields(fields, v.Field(i).Interface(), currentNums...); err != nil {
				return err
			}
			continue
		}

		// insert field into fields
		err = putField(fields, v.Field(i).Interface(), currentNums...)
		if err != nil {
			return err
		}
	}
	return nil
}

func parse(value any) (protobufFields, error) {
	fields := make(protobufFields)

	if err := parseFields(fields, value); err != nil {
		return nil, err
	}

	return fields, nil
}

func encodeValue(value any) []byte {
	switch value := value.(type) {
	case int32:
		return protowire.AppendVarint(nil, uint64(value))
	case int64:
		return protowire.AppendVarint(nil, uint64(value))
	case uint32:
		return protowire.AppendVarint(nil, uint64(value))
	case uint64:
		return protowire.AppendVarint(nil, uint64(value))
	case string:
		return protowire.AppendString(nil, value)
	case bool:
		if value {
			return protowire.AppendVarint(nil, uint64(1))
		}
		return protowire.AppendVarint(nil, uint64(0))
	case float32:
		return protowire.AppendFixed32(nil, uint32(value))
	case float64:
		return protowire.AppendFixed64(nil, uint64(value))
	default:
		return []byte{}
	}
}

func getWireType(value any) protowire.Type {
	switch value.(type) {
	case int32, int64, uint32, uint64, bool:
		return protowire.VarintType
	case string:
		return protowire.BytesType
	case float32:
		return protowire.Fixed32Type
	case float64:
		return protowire.Fixed64Type
	default:
		return 0
	}
}

// TODO: Implement slices
func encodeField(field *protobufField) (buffer []byte) {

	// sub struct
	if field.subFields != nil {
		buffer = protowire.AppendTag(buffer, protowire.Number(field.num), protowire.BytesType)
		buffer = protowire.AppendBytes(buffer, encode(field.subFields))
		return
	}

	buffer = protowire.AppendTag(buffer, protowire.Number(field.num), getWireType(field.value))
	buffer = append(buffer, encodeValue(field.value)...)
	return
}

func encode(fields protobufFields) (buffer []byte) {
	for _, field := range fields {
		buffer = append(buffer, encodeField(field)...)
	}

	return
}

func Marshal(value any) ([]byte, error) {
	fields, err := parse(value)
	if err != nil {
		return nil, err
	}

	return encode(fields), nil
}
