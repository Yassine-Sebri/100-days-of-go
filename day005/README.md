# Arrays and slices

Arrays allow us to store multiple elements of the same type in a variable in a particular order.

As a demo, let's make a `Sum` function. `Sum` will take an array of numbers and return the total.

## Write the test

As usual, we start by writing the test

```go
package main

import "testing"

func TestSum(t *testing.T) {
    numbers := [5]int{1, 2, 3, 4, 5}

    got := Sum(numbers)
    want := 15

    if got != want {
        t.Errorf("Expected %d given %v but got %d", want, numbers, got)
    }
}
```

Arrays have a fixed capacity which we define when we declare the variable

We run the test and it fails, as expected.

## Write enough code to make the test pass

```go
package main

func Sum(numbers [5]int) int {
    sum := 0
    for i := 0; i < 5; i++ {
        sum += numbers[i]
    }
    return sum
}
```

We just have to sum all the array values.

## Refactor

Let's introduce range to help clean up our code

```go
package main

func Sum(numbers [5]int) int {
    sum := 0
    for _, number := range numbers {
        sum += number
    }
    return sum
}
```

`range` lets us iterate over an array. On each iteration, `range` returns two values - the index and the value. We are choosing to ignore the index value by using `_` [blank identifier](https://go.dev/doc/effective_go#blank).

## Arrays and their type

An interesting property of arrays is that the size is encoded in its type. If we try to pass an `[4]int` into a function that expects `[5]int`, it won't compile. They are different types so it's just the same as trying to pass a `string` into a function that wants an `int`.
Go has slices which do not encode the size of the collection and instead can have any size.
The next requirement will be to sum collections of varying sizes.

## Write the test first

We will now use the slice type which allows us to have collections of any size. The syntax is very similar to arrays, we just omit the size when declaring them

```go
package main

import "testing"

func TestSum(t *testing.T) {

    t.Run("collection of 5 numbers", func(t *testing.T) {
        numbers := [5]int{1, 2, 3, 4, 5}

        got := Sum(numbers)
        want := 15

        if got != want {
            t.Errorf("Expected %d given %v but got %d", want, numbers, got)
        }
    })

    t.Run("collection of any size", func(t *testing.T) {
        numbers := []int{1, 2, 3}

        got := Sum(numbers)
        want := 6

        if got != want {
            t.Errorf("Expected %d given %v but got %d", want, numbers, got)
        }
    })
}
```

## Write enough code to make it pass

The problem here is we can either

- Break the existing API by changing the argument to Sum to be a slice rather than an array. When we do this, we will potentially ruin someone's day because our other test will no longer compile!
- Create a new function

In our case, no one else is using our function, so rather than having two functions to maintain, let's have just one.

```go
package main

func Sum(numbers []int) int {
    sum := 0
    for _, number := range numbers {
        sum += number
    }
    return sum
}
```

If we try to run the tests they will still not compile, we will have to change the first test to pass in a slice rather than an array.

## Refactoring

It is important to question the value of tests. It should not be a goal to have as many tests as possible, but rather to have as much *confidence* as possible in the code base. Having too many tests can turn in to a real problem and it just adds more overhead in maintenance. **Every test has a cost**.

In our case, you can see that having two tests for this function is redundant. If it works for a slice of one size it's very likely it'll work for a slice of any size (within reason).

Go's built-in testing toolkit features a coverage tool. Whilst striving for 100% coverage should not be the end goal, the coverage tool can help identify areas of your code not covered by tests.

When we run `go test -cover`, we get

```text
PASS
coverage: 100.0% of statements
```

If we delete one of the tests and run the command again, we get the same result.
