package protobuf

import (
	"google.golang.org/protobuf/encoding/protowire"
)

type Field struct {
	WireType protowire.Type
	Num      protowire.Number
	Value    any
}

func Find(b []byte, nums ...uint32) *Field {
	fields := FindN(b, 1, nums...)

	if len(fields) == 0 {
		return nil
	}

	return &fields[0]
}

func FindAll(b []byte, nums ...uint32) []Field {
	return FindN(b, 0, nums...)
}

func FindN(b []byte, limit uint, nums ...uint32) []Field {
	var fields []Field

	if len(nums) == 0 {
		return fields
	}

	if len(nums) == 1 {
		for len(b) != 0 {
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
				b = b[n:]
			}

			if num == protowire.Number(nums[0]) {
				fields = append(fields, Field{
					WireType: wireType,
					Num:      num,
					Value:    value,
				})
			}
		}
		return fields
	}

	for len(b) != 0 {
		num, wireType, n := protowire.ConsumeTag(b)
		b = b[n:]

		if num == protowire.Number(nums[0]) {
			if wireType != protowire.BytesType {
				panic("must be a bytes type") // TODO: DEBUG
			}

			value, n := protowire.ConsumeBytes(b)
			b = b[n:]
			fields = append(fields, FindN(value, limit, nums[1:]...)...)
			continue
		}

		switch wireType {
		case protowire.VarintType:
			_, n = protowire.ConsumeVarint(b)
			b = b[n:]
		case protowire.Fixed32Type:
			_, n = protowire.ConsumeFixed32(b)
			b = b[n:]
		case protowire.Fixed64Type:
			_, n = protowire.ConsumeFixed64(b)
			b = b[n:]
		case protowire.BytesType:
			_, n = protowire.ConsumeBytes(b)
			b = b[n:]
		}
	}

	return fields
}

func Unmarshal(b []byte, v any) error {
	return nil
}
