package protobuf

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
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
		return errors.New("nums cannot be empty")
	}

	num := nums[0]

	// Insert value
	if len(nums) == 1 {
		if _, ok := fields[num]; ok {
			return errors.New("values and fields cannot be used at the same time")
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
			return errors.New("values and fields cannot be used at the same time")
		}
		if field.fields == nil {
			return errors.New("messages and fields cannot be used at the same time")
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
		return nil, errors.New("value must be a struct")
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

func Marshal(value any) ([]byte, error) {
	_, err := parse(value)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
