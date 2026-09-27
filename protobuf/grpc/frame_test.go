package grpc

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/etienne-hd/raw-protobuf/protobuf"
)

func generateProtobuf(t *testing.T) []byte {
	t.Helper()

	type Person struct {
		Name string `num:"1"`
		Age  uint32 `num:"2"`
	}

	payload, err := protobuf.Marshal(Person{
		Name: "Etienne",
		Age:  20,
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	return payload
}

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
		compression Compression
	}{
		{
			name:        "no compression",
			compression: CompressionNone,
		},
		{
			name:        "gzip",
			compression: CompressionGzip,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := generateProtobuf(t)

			framed, err := Frame(payload, tt.compression)
			if err != nil {
				t.Fatalf("frame: %v", err)
			}

			compressed := tt.compression != CompressionNone
			checkHeader(t, framed[:5], compressed, uint32(len(framed[5:])))

			if !compressed && !bytes.Equal(framed[5:], payload) {
				t.Fatalf("payload: got %x, want %x", framed[5:], payload)
			}
		})
	}
}
