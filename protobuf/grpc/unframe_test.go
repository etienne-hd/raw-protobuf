package grpc_test

import (
	"bytes"
	"testing"

	"github.com/etienne-hd/raw-protobuf/protobuf/grpc"
)

func TestUnframe(t *testing.T) {
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := generateProtobuf(t)

			framed, err := grpc.Frame(payload, test.compression)
			if err != nil {
				t.Fatalf("frame: %v", err)
			}

			unframed, err := grpc.Unframe(framed, test.compression)
			if err != nil {
				t.Fatalf("unframed: %v", err)
			}

			if !bytes.Equal(payload, unframed) {
				t.Fatalf("payload and unframed must be equal: got %x, want %x", unframed, payload)
			}
		})
	}
}
