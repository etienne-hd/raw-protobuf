package grpc

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
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

func (algorithm Compression) uncompress(payload []byte) ([]byte, error) {
	switch algorithm {
	case CompressionNone:
		return payload, nil

	case CompressionGzip:
		reader, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("create gzip reader: %w", err)
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("uncompress gzip payload: %w", err)
		}

		return data, nil

	default:
		return nil, fmt.Errorf("unsupported compression algorithm: %d", algorithm)
	}
}
