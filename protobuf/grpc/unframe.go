package grpc

import (
	"encoding/binary"
	"fmt"
)

func Unframe(payload []byte, algorithm Compression) ([]byte, error) {
	if payload[0] == 1 && algorithm == CompressionNone {
		return nil, fmt.Errorf("payload is compressed but compression algorithm is none")
	} else if payload[0] == 0 && algorithm != CompressionNone {
		return nil, fmt.Errorf("payload is not compressed but compression algorithm is not none")
	}

	len := binary.BigEndian.Uint32(payload[1:5])

	return algorithm.uncompress(payload[5 : len+5])
}
