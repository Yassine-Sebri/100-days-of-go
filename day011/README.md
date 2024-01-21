# Dependecy Injection

We want to write a function that greets someone, just like we did in the hello-world chapter but this time we are going to be testing the actual *printing*.

Just to recap, here is what that function could look like

```go
func Greet(name string) {
 fmt.Printf("Hello, %s", name)
}
```

But how can we test this? Calling `fmt.Printf` prints to stdout, which is pretty hard for us to capture using the testing framework.

What we need to do is to be able to **inject** (which is just a fancy word for pass in) the dependency of printing.

**Our function doesn't need to care where or how the printing happens, so we should accept an interface rather than a concrete type.**

If we look at the source code of [fmt.Printf](https://pkg.go.dev/fmt#Printf) we can see a way for us to hook in

```go
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
 return Fprintf(os.Stdout, format, a...)
}
```

Interesting! Under the hood `Printf` just calls `Fprintf` passing in `os.Stdout`.

What exactly is an `os.Stdout`? What does `Fprintf` expect to get passed to it for the 1st argument?

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
 p := newPrinter()
 p.doPrintf(format, a)
 n, err = w.Write(p.buf)
 p.free()
 return
}
```

An `io.Writer`

```go
type Writer interface {
 Write(p []byte) (n int, err error)
}
```

From this we can infer that `os.Stdout` implements `io.Writer`; `Printf` passes `os.Stdout` to `Fprintf` which expects an `io.Writer`.

So we know under the covers we're ultimately using Writer to send our greeting somewhere. Let's use this existing abstraction to make our code testable and more reusable.

## Write the test first

```go
func TestGreet(t *testing.T) {
 buffer := bytes.Buffer{}
 Greet(&buffer, "Chris")

 got := buffer.String()
 want := "Hello, Chris"

 if got != want {
  t.Errorf("got %q want %q", got, want)
 }
}
```

The `Buffer` type from the `bytes` package implements the `Writer` interface, because it has the method `Write(p []byte) (n int, err error)`.

So we'll use it in our test to send in as our `Writer` and then we can check what was written to it after we invoke `Greet`

## Write enough code to make it pass

```go
func Greet(writer *bytes.Buffer, name string) {
 fmt.Fprintf(writer, "Hello, %s", name)
}
```

## Refactor

As discussed earlier `fmt.Fprintf` allows us to pass in an `io.Writer` which we know both `os.Stdout` and `bytes.Buffer` implement.

If we change our code to use the more general purpose interface we can now use it in both tests and in our application.

```go
package main

import (
 "fmt"
 "io"
 "os"
)

func Greet(writer io.Writer, name string) {
 fmt.Fprintf(writer, "Hello, %s", name)
}

func main() {
 Greet(os.Stdout, "Elodie")
}
```
