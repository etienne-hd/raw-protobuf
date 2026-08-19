package protobuf

import (
	"testing"
)

func assertField(t *testing.T, fields map[uint]*protobufField, value any, nums ...uint) {
	t.Helper()

	field := fields

	for i, num := range nums {
		if field[num] == nil {
			t.Fatalf("field %d does not exist", i)
		}

		if i == len(nums) - 1 {
			if field[num].value != value {
				t.Fatalf("field value = %v, want %v", field[num].value, value)
			}
			continue
		}

		if field[num].subField == nil {
			t.Fatalf("field %d has no subField", i)
		}
		field = field[num].subField
	}
}

func TestPutField(t *testing.T) {
	fields := make(map[uint]*protobufField)

	nums := []uint{1, 2, 3}
	value := 123

	if err := putField(fields, value, nums...); err != nil {
		t.Fatalf("putField() returned an unexpected error: %v", err)
	}

	assertField(t, fields, value, nums...)

	if err := putField(fields, value, nums...); err == nil {
		t.Fatal("putField() allowed overwriting an existing field")
	}

	nums = []uint{1, 2, 4, 1}
	value = 1234

	if err := putField(fields, value, nums...); err != nil {
		t.Fatalf("putField() returned an unexpected error: %v", err)
	}

	assertField(t, fields, value, nums...)

	nums = []uint{1, 2, 4, 1, 2}
	if err := putField(fields, value, nums...); err == nil {
		t.Fatalf("putField() allowed subField and value at the same time")
	}
}
