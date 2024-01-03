# Hello, Golang

## Introduction

As per tradition, the first program in a new lanuage should be a [Hello World](https://en.wikipedia.org/wiki/%22Hello,_World!%22_program) program.

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Golang")
}
```

## How it works

The code starts with a package declaration, stating that this file belongs to the main package. In Go, the main package is special and serves as the entry point for executable programs.

Next, there's an import statement. It brings in the "fmt" package, which stands for "format." The "fmt" package contains the "Println" function that we use to print text.

The func keyword is how we define a function. In Go, the main function is where the execution of the program begins.

When this program is run, it will output "Hello, Golang" to the console. To run it type `go run hello.go`.

## Domain code and side effects

Domain code refers to the core logic or functionality of a program that directly relates to the problem domain or business logic it aims to address. A side effect is an interaction with the outside world. These interactions include:

- Reading a file, and/or writing to a file
- Making a network request (calling an API, downloading a file...)
- Reading from a global state (e.g. global variable, a parameter from the parent's closure...)
- Throwing/intercepting an exception
- For web applications, using DOM objects and methods
- Logging messages to the console

In functional programs, functions have to be pure, i.e. side-effect-free because side effects are not deterministic, and they are hard to reason about.

In order to be able to properly write tests, we separate side effects from domain code

```go
package main

import "fmt"

func Hello() string {
    return "Hello, Golang"
}

func main() {
    fmt.Println(Hello())
}
```

We have created a new function again with `func` but this time we've added another keyword `string` in the definition. This means this function returns a string.

Now create a new file called `hello_test.go` where we are going to write a test for our `Hello` function

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello()
    want := "Hello, Golang"
    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

## Go Modules

In order to be able to test our code, we need to first initialize and create a new module for our Go project. This is done using the `mod init SOMENAME` command where `SOMENAME` is the name of the project.

That will create a new file with the following contents:

```go
module hello

go 1.20
```

The `go.mod` file defines the module name and version for the project. The module name is a unique identifier for the project, and it helps in specifying dependencies in a consistent manner.

## Back to Testing

Go provides built-in support for testing through its testing package, allowing developers to write and run tests for their code easily. The testing package is part of the Go standard library, and it includes testing-related functions and conventions that make it straightforward to create tests and ensure code quality.

### Writing tests

Writing a test is just like writing a function, with a few rules

- It needs to be in a file with a name like `xxx_test.go`
- The test function must start with the word `Test`
- The test function takes one argument only `t *testing.T`
- In order to use the `*testing.T` type, you need to import "testing", like we did with "fmt" in the other file

For now, it's enough to know that your t of type `*testing.T` is your "hook" into the testing framework so you can do things like `t.Fail()` when you want to fail.

Other new topics covered:

`if`

If statements in Go are very much like other programming languages.

`Declaring variables`

We're declaring some variables with the syntax varName := value, which lets us re-use some values in our test for readability.

`t.Errorf`

We are calling the Errorf method on our t which will print out a message and fail the test. The f stands for format which allows us to build a string with values inserted into the placeholder values %q.
