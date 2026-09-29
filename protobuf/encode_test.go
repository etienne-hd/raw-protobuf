package protobuf_test

import (
	"testing"

	"github.com/etienne-hd/raw-protobuf/protobuf"
)

func TestMarshal(t *testing.T) {
	type Person struct {
		Name string `num:"1"`
		Age  uint32 `num:"2"`
	}

	want := Person{
		Name: "Etienne",
		Age:  20,
	}

	payload, err := protobuf.Marshal(want)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	name := protobuf.Find(payload, 1)
	if name == nil {
		t.Fatal("Find() returned nil for field 1")
	}

	if got := name.String(); got != want.Name {
		t.Errorf("name = %q, want %q", got, want.Name)
	}

	age := protobuf.Find(payload, 2)
	if age == nil {
		t.Fatal("Find() returned nil for field 2")
	}

	if got := age.Int(); got != uint64(want.Age) {
		t.Errorf("age = %d, want %d", got, want.Age)
	}
}

func TestMarshalNested(t *testing.T) {
	type Address struct {
		City    string `num:"1"`
		Country string `num:"2"`
	}

	type Person struct {
		Name    string   `num:"1"`
		Age     uint32   `num:"2"`
		Address Address  `num:"3"`
		Aliases []string `num:"4"`
	}

	want := Person{
		Name: "Etienne",
		Age:  20,
		Address: Address{
			City:    "Angouleme",
			Country: "France",
		},
		Aliases: []string{"etienne", "ehode"},
	}

	payload, err := protobuf.Marshal(want)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	t.Run("name", func(t *testing.T) {
		field := protobuf.Find(payload, 1)
		if field == nil {
			t.Fatal("Find() returned nil")
		}

		if got := field.String(); got != want.Name {
			t.Errorf("name = %q, want %q", got, want.Name)
		}
	})

	t.Run("age", func(t *testing.T) {
		field := protobuf.Find(payload, 2)
		if field == nil {
			t.Fatal("Find() returned nil")
		}

		if got := field.Int(); got != uint64(want.Age) {
			t.Errorf("age = %d, want %d", got, want.Age)
		}
	})

	t.Run("address city", func(t *testing.T) {
		field := protobuf.Find(payload, 3, 1)
		if field == nil {
			t.Fatal("Find() returned nil")
		}

		if got := field.String(); got != want.Address.City {
			t.Errorf("city = %q, want %q", got, want.Address.City)
		}
	})

	t.Run("address country", func(t *testing.T) {
		field := protobuf.Find(payload, 3, 2)
		if field == nil {
			t.Fatal("Find() returned nil")
		}

		if got := field.String(); got != want.Address.Country {
			t.Errorf("country = %q, want %q", got, want.Address.Country)
		}
	})
}
