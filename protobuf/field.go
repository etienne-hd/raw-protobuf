package protobuf

import (
	"math"

	"google.golang.org/protobuf/encoding/protowire"
)

type Field struct {
	wireType protowire.Type
	value    any
}

func (f *Field) Interface() any {
	return f.value
}

func (f *Field) Int() uint64 {
	if f.wireType != protowire.VarintType {
		panic("Field.Int() called on non-int field")
	}

	return f.value.(uint64)
}

func (f *Field) Bool() bool {
	if f.wireType != protowire.VarintType {
		panic("Field.Bool() called on non-bool field")
	}

	if v := f.value.(uint64); v == 1 {
		return true
	} else {
		return false
	}
}

func (f *Field) Float32() float32 {
	if f.wireType != protowire.Fixed32Type {
		panic("Field.Float32() called on non-float32 field")
	}

	return math.Float32frombits(f.value.(uint32))
}

func (f *Field) Float64() float64 {
	if f.wireType != protowire.Fixed64Type {
		panic("Field.Float64() called on non-float64 field")
	}

	return math.Float64frombits(f.value.(uint64))
}

func (f *Field) String() string {
	if f.wireType != protowire.BytesType {
		panic("Field.String() called on non-string field")
	}

	switch value := f.value.(type) {
	case []byte:
		return string(value)
	default:
		panic("Field.String() called on invalid field value")
	}
}

func (f *Field) Bytes() []byte {
	if f.wireType != protowire.BytesType {
		panic("Field.Bytes() called on non-bytes field")
	}

	return f.value.([]byte)
}
