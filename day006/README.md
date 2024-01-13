# Structs, methods & interfaces

Suppose that we need some geometry code to calculate the perimeter of a rectangle given a height and width. We can write a `Perimeter(width float64, height float64)` function, where `float64` is for floating-point numbers like `123.45`.

## Write the test

```go
package main

import "testing"

func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    want := 40.0

    if got != want {
        t.Errorf("Expected %.2f but got %.2f", got, want)
    }
}
```

The `f` is for our `float64` and the `.2` means print 2 decimal places.

## Write enough code to make it pass

```go
package main

func Perimeter(width, height float64) float64 {
    return 2 * (width + height)
}
```

So far, so easy. Now let's create a function called Area(width, height float64) which returns the area of a rectangle. We start by writing the test

```go
func TestArea(t *testing.T) {
    got := Area(12.0, 6.0)
    want := 72.0

    if got != want {
        t.Errorf("Expected %.2f but got %.2f", got, want)
    }
}
```

then we write the code to pass the test

```go
func Area(width, height float64) float64 {
    return width * height
}
```

## Refactor

Our code does the job, but it doesn't contain anything explicit about rectangles. An unwary developer might try to supply the width and height of a triangle to these functions without realising they will return the wrong answer.

We could just give the functions more specific names like `RectangleArea`. A neater solution is to define our own type called `Rectangle` which encapsulates this concept for us.

We can create a simple type using a [struct](https://go.dev/ref/spec#Struct_types).

We declare a struct like this

```go
type Rectangle struct {
    Width  float64
    Height float64
}
```

Now let's refactor the tests to use `Rectangle` instead of plain `float64`

```go
func TestPerimeter(t *testing.T) {
    rectangle := Rectangle{10.0, 10.0}
    got := Perimeter(rectangle)
    want := 40.0

    if got != want {
        t.Errorf("Expected %.2f but got %.2f", got, want)
    }
}

func TestArea(t *testing.T) {
    rectangle := Rectangle{12.0, 6.0}
    got := Area(rectangle)
    want := 72.0

    if got != want {
        t.Errorf("Expected %.2f but got %.2f", got, want)
    }
}
```

We can access the fields of a struct with the syntax of `myStruct.field`.

```go
func Perimeter(rectangle Rectangle) float64 {
    return 2 * (rectangle.Width + rectangle.Height)
}

func Area(rectangle Rectangle) float64 {
    return rectangle.Width * rectangle.Height
}
```

The tests are green again.

Our next requirement is to write an Area function for circles.

## Write the test first

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12.0, 6.0}
        got := Area(rectangle)
        want := 72.0

        if got != want {
            t.Errorf("Expected %g but got %g", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10.0}
        got := Area(circle)
        want := 314.1592653589793

        if got != want {
            t.Errorf("Expected %g but got %g", got, want)
        }
    })
}
```

## Write the minimal amount of code for the test to pass

We need to define our `Circle` type.

```go
type Circle struct {
    Radius float64
}
```

Because the `Area` function takes a `Rectangle` as input, we cannot use it to calculate the area of `circle`. To solve this, we can use [methods](https://go.dev/ref/spec#Method_declarations).

### What are methods?

So far we have only been writing *functions* but we have been using some methods. When we call `t.Errorf` we are calling the method `Errorf` on the instance of our `t` (`testing.T`).

A method is a function with a receiver. A method declaration binds an identifier, the method name, to a method, and associates the method with the receiver's base type.

Methods are very similar to functions but they are called by invoking them on an instance of a particular type. Where you can just call functions wherever you like, such as `Area(rectangle)` you can only call methods on "things".

Let's rewrite the tests to use methods

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := rectangle.Area()
        want := 72.0

        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := circle.Area()
        want := 314.1592653589793

        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    })
}
```

Now let's define the methods

```go
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}
```

The syntax for declaring methods is almost the same as functions and that's because they're so similar. The only difference is the syntax of the method receiver `func (receiverName ReceiverType) MethodName(args)`.

When the method is called on a variable of that type, you get your reference to its data via the `receiverName` variable. In many other programming languages this is done implicitly and you access the receiver via `this`.

It is a convention in Go to have the receiver variable be the first letter of the type.

```text
r Rectangle
```

We run the tests and they pass.

## Refactoring

There is some duplication in our tests.

All we want to do is take a collection of shapes, call the `Area()` method on them and then check the result.

We want to be able to write some kind of `checkArea` function that we can pass both `Rectangle`s and `Circle`s to, but fail to compile if we try to pass in something that isn't a shape.

With Go, we can codify this intent with **interfaces**.

[Interfaces](https://go.dev/ref/spec#Interface_types) are a very powerful concept in statically typed languages like Go because they allow you to make functions that can be used with different types and create highly-decoupled code whilst still maintaining type-safety.

Let's introduce this by refactoring our tests.

```go
func TestArea(t *testing.T) {

    checkArea := func(t testing.TB, shape Shape, want float64) {
        t.Helper()
        got := shape.Area()

        if got != want {
            t.Errorf("Expected %g but got %g", want, got)
        }
    }

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        checkArea(t, rectangle, 72.0)
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        checkArea(t, circle, 314.1592653589793)
    })
}
```

How does something become a shape? We just tell Go what a `Shape` is using an interface declaration

```go
type Shape interface {
    Area() float64
}
```

We're creating a new type just like we did with `Rectangle` and `Circle` but this time it is an `interface` rather than a `struct`.

Once we add this to the code, the tests will pass.

## Wait, what?

This is quite different to interfaces in most other programming languages. Normally you have to write code to say `My type Foo implements interface Bar`.

But in our case

- `Rectangle` has a method called `Area` that returns a `float64` so it satisfies the `Shape` interface
- `Circle` has a method called `Area` that returns a `float64` so it satisfies the `Shape` interface
- `string` does not have such a method, so it doesn't satisfy the interface

In Go **interface resolution is implicit**. If the type you pass in matches what the interface is asking for, it will compile.

## Further refactoring

Now that we have some understanding of structs we can introduce "table driven tests".

[Table driven tests](https://go.dev/wiki/TableDrivenTests) are useful when you want to build a list of test cases that can be tested in the same manner.

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        want  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.want {
            t.Errorf("got %g want %g", got, tt.want)
        }
    }

}
```

You can see how it would be very easy for a developer to introduce a new shape, implement `Area` and then add it to the test cases. In addition, if a bug is found with `Area` it is very easy to add a new test case to exercise it before fixing it.

Table driven tests are a great fit when we wish to test various implementations of an interface, or if the data being passed in to a function has lots of different requirements that need testing.
