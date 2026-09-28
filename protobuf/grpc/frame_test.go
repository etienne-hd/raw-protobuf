package grpc_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/etienne-hd/raw-protobuf/protobuf/grpc"
)

func checkHeader(t *testing.T, header []byte, compressed bool, payloadLen uint32) {
	t.Helper()

	var wantFlag byte
	if compressed {
		wantFlag = 1
	}

	if header[0] != wantFlag {
		t.Fatalf("compression flag: got %d, want %d", header[0], wantFlag)
	}

	gotLen := binary.BigEndian.Uint32(header[1:5])
	if gotLen != payloadLen {
		t.Fatalf("payload length: got %d, want %d", gotLen, payloadLen)
	}
}

func TestFrame(t *testing.T) {
	tests := []struct {
		name        string
		compression grpc.Compression
	}{
		{
			name:        "no compression",
			compression: grpc.CompressionNone,
		},
		{
			name:        "gzip",
			compression: grpc.CompressionGzip,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := generateProtobuf(t)

			framed, err := grpc.Frame(payload, tt.compression)
			if err != nil {
				t.Fatalf("frame: %v", err)
			}

			compressed := tt.compression != grpc.CompressionNone
			checkHeader(t, framed[:5], compressed, uint32(len(framed[5:])))

			if !compressed && !bytes.Equal(framed[5:], payload) {
				t.Fatalf("payload: got %x, want %x", framed[5:], payload)
			}
		})
	}
}
