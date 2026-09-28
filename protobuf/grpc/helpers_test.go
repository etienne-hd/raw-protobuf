package grpc_test

import (
	"testing"

	"github.com/etienne-hd/raw-protobuf/protobuf"
)

func generateProtobuf(t *testing.T) []byte {
	t.Helper()

	type Person struct {
		Name    string `num:"1"`
		Age     uint32 `num:"2"`
		Friends []Person
	}

	payload, err := protobuf.Marshal(Person{
		Name: "Etienne",
		Age:  20,
		Friends: []Person{
			{
				Name:    "A",
				Age:     1,
				Friends: []Person{},
			},
			{
				Name: "Z",
				Age:  99,
				Friends: []Person{
					{
						Name: "B",
						Age:  123,
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	return payload
}
