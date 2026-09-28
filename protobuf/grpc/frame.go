package grpc

import (
	"encoding/binary"
)

func Frame(payload []byte, algorithm Compression) ([]byte, error) {
	payload, err := algorithm.compress(payload)
	if err != nil {
		return nil, err
	}

	header := make([]byte, 5)

	if algorithm != CompressionNone {
		header[0] = 1
	}

	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))

	return append(header, payload...), nil
}
