package protobuf

import "google.golang.org/protobuf/encoding/protowire"

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
