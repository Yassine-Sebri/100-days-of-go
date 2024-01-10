package adder

import "testing"

func TestAdder(t *testing.T) {
	got := Add(2, 2)
	want := 4
	if got != want {
		t.Errorf("Expected '%d' but got '%d'", want, got)
	}
}
