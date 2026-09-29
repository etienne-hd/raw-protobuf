package protobuf

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type Person struct {
		Name string `num:"1"`
		Age  uint32 `num:"2"`
	}

	want := Person{
		Name: "Etienne",
		Age:  20,
	}

	payload, err := Marshal(want)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var got Person

	if err := Unmarshal(payload, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal() = %+v, want %+v", got, want)
	}
}

func TestUnmarshalNested(t *testing.T) {
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
		Aliases: []string{
			"etienne",
			"ehode",
		},
	}

	payload, err := Marshal(want)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var got Person

	if err := Unmarshal(payload, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal() = %+v, want %+v", got, want)
	}
}
