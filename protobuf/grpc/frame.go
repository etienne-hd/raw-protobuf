package grpc

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
)

type Compression uint8

const (
	CompressionNone Compression = iota
	CompressionGzip
)

func (algorithm Compression) compress(payload []byte) ([]byte, error) {
	switch algorithm {
	case CompressionNone:
		return payload, nil

	case CompressionGzip:
		var buf bytes.Buffer

		writer := gzip.NewWriter(&buf)
		if _, err := writer.Write(payload); err != nil {
			return nil, err
		}

		if err := writer.Close(); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil

	default:
		return nil, fmt.Errorf("unsupported compression algorithm: %d", algorithm)
	}
}

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
