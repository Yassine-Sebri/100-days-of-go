package iteration

import "testing"

func TestRepeat(t *testing.T) {
	t.Run("Repeat the character 5 times", func(t *testing.T) {
		got := Repeat("a", 5)
		want := "aaaaa"
		assertCorrectMessage(t, got, want)
	})

	t.Run("Repeat the character 3 times", func(t *testing.T) {
		got := Repeat("b", 3)
		want := "bbb"
		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("Expected %q but got %q", want, got)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}
}
