# Iteration

To do stuff repeatedly in Go, we'll need `for`. In Go there are no `while`, `do`, `until` keywords, we can only use `for`.
Let's write a test for a function that repeats a character 5 times.

```go
package repeat

import "testing"

func test_repeat(t *testing.T) {
    got := repeat("a")
    want := "aaaaa"

    if got != want {
        t.Errorf("Expected %q but got %q", want, got)
    }
}
```

When we run the test, it fails because repeat is not defined.

## Write enough code to make it pass

First, we define the `Repeat` function to fix the error.

```go
package iteration

func Repeat(character string) string {
    return ""
}
```

We get the following error

```text
repeat_test.go:10: Expected "aaaaa" but got ""
```

Time to implement the logic needed for the function to work as expected

```go
package iteration

func Repeat(character string) string {
    repeated := ""
    for i := 0; i < 5; i++ {
        repeated = repeated + character
    }
    return repeated
}
```

The for syntax is very unremarkable and follows most C-like languages. We run the test and it passes.

## Refactor

Now that the test has passed, we refactor the code.

```go
package iteration

const repeatCount = 5

func Repeat(character string) string {
    repeated := ""
    for i := 0; i < repeatCount; i++ {
        repeated += character
    }
    return repeated
}
```

`+=` called "the Add AND assignment operator", adds the right operand to the left operand and assigns the result to left operand. It works with other types like integers.
We also remove magic numbers, which refers to a numeric constants or values that are used directly in the code without any explanation or symbolic representation. In our case, it refers to the number 5 which is the number of repeatitions.

## Benchmarking

Writing benchmarks in Go is another first-class feature of the language and it is very similar to writing tests.

```go
func BenchmarkRepeat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Repeat("a")
    }
}
```

When the benchmark code is executed, it runs b.N times and measures how long it takes.
The amount of times the code is run shouldn't matter to you, the framework will determine what is a "good" value for that to let you have some decent results.
To run the benchmarks we use `go test -bench=.` (or `go test -bench="."` on Windows Powershell).

## More iteration

This time, we want to specify the number of times the character should be repeated. We rewrite our tests to reflect this.

```go
func TestRepeat(t *testing.T) {
    t.Run("Repeat the character 5 times", func(t *testing.T) {
        got := Repeat("a", 5)
        want := "aaaaa"

        if got != want {
            t.Errorf("Expected %q but got %q", want, got)
        }
    })

    t.Run("Repeat the character 3 times", func(t *testing.T) {
        got := Repeat("b", 3)
        want := "bbb"

        if got != want {
            t.Errorf("Expected %q but got %q", want, got)
        }
    })
}
```

We modify the `Repeat` function to pass the test

```go
package iteration

func Repeat(character string, repeatCount int) string {
    repeated := ""
    for i := 0; i < repeatCount; i++ {
        repeated += character
    }
    return repeated
}
```

and finally, we Refactor

```go
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
```
