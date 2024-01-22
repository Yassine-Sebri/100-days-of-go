# Mocking

Let's write a program which counts down from 3, printing each number on a new line (with a 1-second pause) and when it reaches zero it will print "Go!" and exit.

```plaintext
3
2
1
Go!
```

We'll tackle this by writing a function called `Countdown` which we will then put inside a `main` program so it looks something like this:

```go
package main

func main() {
 Countdown()
}
```

## Write the test first

Our software needs to print to stdout and we saw how we could use Dependency Injection (DI) to facilitate testing this in the DI section.

```go
func TestCountdown(t *testing.T) {
 buffer := &bytes.Buffer{}

 Countdown(buffer)

 got := buffer.String()
 want := "3"

 if got != want {
  t.Errorf("got %q want %q", got, want)
 }
}
```

## Write enough code to make it pass

```go
func Countdown(out *bytes.Buffer) {
 fmt.Fprint(out, "3")
}
```

We're using `fmt.Fprint` which takes an `io.Writer` (like `*bytes.Buffer`) and sends a `string` to it. The test should pass.

## Refactor

We know that while `*bytes.Buffer` works, it would be better to use a general purpose interface instead.

```go
package main

import (
 "fmt"
 "io"
 "os"
)

func Countdown(out io.Writer) {
 fmt.Fprint(out, "3")
}

func main() {
 Countdown(os.Stdout)
}
```

Next we can make it print 2,1 and then "Go!".

## Write the test

```go
func TestCountdown(t *testing.T) {
 buffer := &bytes.Buffer{}

 Countdown(buffer)

 got := buffer.String()
 want := `3
2
1
Go!`

 if got != want {
  t.Errorf("got %q want %q", got, want)
 }
}
```

The backtick syntax is another way of creating a `string` but lets you include things like newlines, which is perfect for our test.

## Write code to make it pass

```go
func Countdown(out io.Writer) {
 for i := 3; i > 0; i-- {
  fmt.Fprintln(out, i)
 }
 fmt.Fprint(out, "Go!")
}
```

We use a `for` loop counting backwards with `i--` and use `fmt.Fprintln` to print to out with our number followed by a newline character. Finally use `fmt.Fprint` to send "Go!" aftward.

## Refactoring

There's not much to refactor other than refactoring some magic values into named constants.

```go
const (
 finalWord      = "Go!"
 countdownStart = 3
)

func Countdown(out io.Writer) {
 for i := countdownStart; i > 0; i-- {
  fmt.Fprintln(out, i)
 }
 fmt.Fprint(out, finalWord)
}
```

If we run the program now, we should get the desired output but we don't have it as a dramatic countdown with the 1-second pauses.

Go lets us achieve this with `time.Sleep`.

```go
func Countdown(out io.Writer) {
 for i := countdownStart; i > 0; i-- {
  fmt.Fprintln(out, i)
  time.Sleep(1 * time.Second)
 }

 fmt.Fprint(out, finalWord)
}
```
