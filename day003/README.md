# Adder

## Introduction

Today, We will just write a simple function that adds two integers.

## Write the test

We start by writing the test. We create a file called `adder_test.go`

```go
package adder

import "testing"

func TestAdder(t *testing.T) {
    got := Add(2, 2)
    want := 4
    if got != want {
        t.Errorf("Expected '%d' but got '%d'", want, got)
    }
}
```

The Add function is still not defined so if we try to run `go test` it returns an error.

## Write logic

We create a new file `adder.go` and we define the functionality we need.

```go
package adder

func Add(a, b int) int {
    return a + b
}
```

We run the tests again and this time we get an ok.
