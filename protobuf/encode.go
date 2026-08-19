package protobuf

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type protobufField struct {
	num      uint
	value    any
	subField map[uint]*protobufField
}

func putField(fields map[uint]*protobufField, value any, nums ...uint) error {
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
			num:      num,
			subField: make(map[uint]*protobufField),
		}

		fields[num] = field
	}

	// A field cannot contain both a value and sub-fields
	if field.subField == nil {
		return errors.New("a field already exist")
	}

	return putField(field.subField, value, nums[1:]...)
}

func getNums(t reflect.StructField) ([]uint, error) {
	rawNums := t.Tag.Get("nums")
	if rawNums == "" {
		return nil, nil
	}

	var nums []uint
	for _, rawNum := range strings.Split(rawNums, ".") {
		num, err := strconv.ParseInt(rawNum, 10, 64)
		if err != nil {
			return nil, err
		}

		nums = append(nums, uint(num))
	}

	return nums, nil
}

func parseFields(fields map[uint]*protobufField, value any, nums ...uint) error {
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
		currentNums := append([]uint{}, nums...)
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

func parse(value any) (map[uint]*protobufField, error) {
	fields := make(map[uint]*protobufField)

	if err := parseFields(fields, value); err != nil {
		return nil, err
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
