package protobuf

type protobufFields []*protobufField

type protobufField struct {
	num       int32
	values    []any
	subFields protobufFields
}

func Marshal(value any) ([]byte, error) {
	return nil, nil
}
