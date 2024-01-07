# Hello, Golang... again

## Introduction

In the previous day, the test was written after the code just so we could get an example of how to write a test and declare a function. From this point on tests will be written first. This is known as test driven development (TDD).

## Hello, You

Our next requirement is to let us specify the recipient of the greeting. Let's start by capturing these requirements in a test.

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello("Yassine")
    want := "Hello, Yassine"
    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

When we run `go test`, we get the following error

```text
./hello_test.go:6:15: too many arguments in call to Hello
        have (string)
        want ()
```

We add an argument to the the hello function

```go
func Hello(name string) string {
    return "Hello, Golang"
}
```

We run `go test` again and we get this error

```text
hello_test.go:9: got "Hello, Golang" want "Hello, Yassine"
```

We finally have a compiling program but it is not meeting our requirements according to the test. Let's make the test pass by using the name argument and concatenate it with `Hello,`

```go
func Hello(name string) string {
    return "Hello, " + name
}
```

## Hello, Golang

The next requirement is when our function is called with an empty string it defaults to printing "Hello, World", rather than "Hello, ".

Start by writing a new failing test

```go
func TestHello(t *testing.T) {
    t.Run("Saying hello to people", func(t *testing.T) {
        got := Hello("Yassine")
        want := "Hello, Yassine"
        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })
    t.Run("Say 'Hello, Golang' when an empty string is supplied", func(t *testing.T) {
        got := Hello("")
        want := "Hello, Golang"
        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })
}
```

Now, let's fix the code

```go
func Hello(name string) string {
    if name == "" {
        name = "Golang"
    }
    return "Hello, " + name
}
```

Now that the tests are passing, we can and should refactor our tests.

```go
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
```

`t.Helper()` is needed to tell the test suite that this method is a helper. By doing this when it fails the line number reported will be in our function call rather than inside our test helper. This will help other developers track down problems easier.
