package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Field struct {
	WireType protowire.Type
	Num      int32
	Value    []byte
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
	// TODO
	return []Field{}
}

func Unmarshal(b []byte, v any) error {
	return nil
}
