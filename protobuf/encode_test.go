package protobuf

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	type Person struct {
		Name string `num:"1"`
		Age uint32 `num:"2"`
	}

	person := Person{
		Name: "Etienne",
		Age: 20,
	}

	b, err := Marshal(person)
	if err != nil {
		t.Fatalf("unexpected error occured: %v", err)
	}

	f1 := Find(b, 1)
	f2 := Find(b, 2)
	if f1 == nil || f2 == nil {
		t.Fatalf("unable to retreive num. (f1 & f2 must be non nil)")
	}

	if f1.String() != person.Name {
		t.Fatalf("unexpected value got %s, want %s", f1.String(), person.Name)
	}
	if f2.Int() != uint64(person.Age) {
		t.Fatalf("unexpected value got %d, want %d", f1.Int(), person.Age)
	}
}