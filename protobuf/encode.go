package protobuf

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type protobufFields []*protobufField

type protobufField struct {
	num         int32
	values      []any
	subFields   protobufFields
	subMessages []protobufFields
}

func getSliceType(v reflect.Value) reflect.Kind {
	if v.Kind() != reflect.Slice {
		return reflect.Invalid
	}

	t := reflect.Invalid
	for i := 0; i < v.Len(); i++ {
		// []any not accepted
		currentType := v.Index(i).Kind()
		if currentType != t && t != reflect.Invalid {
			return reflect.Invalid
		}
		t = currentType
	}
	return t
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

func getValues(value any) (values []any) {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Index(i).Interface())
		}
		return
	}

	values = append(values, value)
	return
}

func insertField(fields *protobufFields, value any, nums ...int32) error {
	if len(nums) == 0 {
		return errors.New("nums cannot be empty")
	}

	// Insert value
	if len(nums) == 1 {
		v := reflect.ValueOf(value)

		// value is a struct
		if v.Kind() == reflect.Struct {
			subFields, err := parseFields(value)
			if err != nil {
				return err
			}

			field := &protobufField{
				num:         nums[0],
				subMessages: []protobufFields{subFields},
			}

			*fields = append(*fields, field)
			return nil
		}

		// value is a slice of struct
		if v.Kind() == reflect.Slice && getSliceType(v) == reflect.Struct {
			field := &protobufField{
				num:         nums[0],
				subMessages: []protobufFields{},
			}

			for i := 0; i < v.Len(); i++ {
				subFields, err := parseFields(v.Index(i).Interface())
				if err != nil {
					return err
				}

				field.subMessages = append(field.subMessages, subFields)
			}

			*fields = append(*fields, field)
			return nil
		}

		field := &protobufField{
			num:       nums[0],
			values:    getValues(value),
			subFields: nil,
		}

		*fields = append(*fields, field)
		return nil
	}

	// Create fields
	var nextField *protobufField
	for _, field := range *fields {
		if field.num == nums[0] {
			if field.values != nil {
				return errors.New("values and subFields cannot be used at the same time")
			}

			nextField = field
		}
	}

	if nextField == nil {
		nextField = &protobufField{
			num:       nums[0],
			subFields: make(protobufFields, 0),
		}

		*fields = append(*fields, nextField)
	}

	return insertField(&nextField.subFields, value, nums[1:]...)
}

func parseFields(value any) (protobufFields, error) {
	fields := make(protobufFields, 0)

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	for i := 0; i < t.NumField(); i++ {
		nums, err := getNums(t.Field(i))
		if err != nil {
			return nil, err
		}

		// No tag
		if nums == nil {
			continue
		}

		if err := insertField(&fields, v.Field(i).Interface(), nums...); err != nil {
			return nil, err
		}
	}

	return fields, nil
}

func parse(value any) (protobufFields, error) {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Struct {
		return nil, errors.New("value must be a struct")
	}

	return parseFields(value)
}

func Marshal(value any) ([]byte, error) {
	_, err := parse(value)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
