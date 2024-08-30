package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Alper")

	got := buffer.String()
	want := "Hello, Alper"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
