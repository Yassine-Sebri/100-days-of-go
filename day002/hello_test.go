package main

import "testing"

func TestHello(t *testing.T) {
	t.Run("Saying hello to people", func(t *testing.T) {
		got := Hello("Yassine")
		want := "Hello, Yassine"
		AssertCorrectMessage(t, got, want)
	})
	t.Run("Say 'Hello, Golang' when an empty string is supplied", func(t *testing.T) {
		got := Hello("")
		want := "Hello, Golang"
		AssertCorrectMessage(t, got, want)
	})
}

func AssertCorrectMessage(t testing.TB, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
