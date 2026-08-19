package protobuf

import (
	"errors"
	"fmt"
)

type protobufField struct {
	num      uint
	value    any
	subField map[uint]*protobufField
}


// DEBUG
func showTree(fields map[uint]*protobufField, depth uint) {
	for _, field := range fields {
		for i := uint(0); i < depth; i++ {
			fmt.Print(" ")
		}
		if len(field.subField) > 0 {
			fmt.Printf("%d\\\n", field.num)
			showTree(field.subField, depth+1)
		} else {
			fmt.Printf("%d\n", field.num)
		}
	}
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

	return putField(fields[nums[0]].subField, value, nums[1:]...)
}

func parse(value any) (map[uint]*protobufField, error) {
	fields := make(map[uint]*protobufField)
	return fields, nil
}

func Marshal(value any) ([]byte, error) {
	_, err := parse(value)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
