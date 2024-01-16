# Maps

To learn about maps, we will build our own dictionary.

## Write the test first

```go
package main

import "testing"

func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    if got != want {
        t.Errorf("Expected %s but got %s", want, &got)
    }
}
```

Declaring a Map is somewhat similar to an array. Except, it starts with the `map` keyword and requires two types. The first is the key type, which is written inside the `[]`. The second is the value type, which goes right after the `[]`.

The key type is special. It can only be a comparable type because without the ability to tell if 2 keys are equal, we have no way to ensure that we are getting the correct value. Comparable types are explained in depth in the [language spec](https://go.dev/ref/spec#Comparison_operators).

The value type, on the other hand, can be any type we want. It can even be another map.

## Write enough code to make it pass

```go
func Search(dictionary map[string]string, word string) string {
    return dictionary[word]
}
```

Getting a value out of a Map is the same as getting a value out of Array `map[key]`.

## Refactor

```go
func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    assertStrings(t, got, want)
}

func assertStrings(t testing.TB, got, want string) {
    t.Helper()

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

We create an `assertStrings` helper to make the implementation more general.

### Using a custom type

We can improve our dictionary's usage by creating a new type around map and making `Search` a method.

In `dictionary_test.go`

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    got := dictionary.Search("test")
    want := "this is just a test"

    assertStrings(t, got, want)
}
```

and in `dictionary.go`

```go
type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
    return d[word]
}
```

Here we created a `Dictionary` type which acts as a thin wrapper around `map`. With the custom type defined, we could create the `Search` method.

The basic search was very easy to implement, but what will happen if we supply a word that's not in our dictionary?

## Write the test

We want the method to return an error so let's write a test for that

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    t.Run("known word", func(t *testing.T) {
        got, _ := dictionary.Search("test")
        want := "this is just a test"

        assertStrings(t, got, want)
    })

    t.Run("unknown word", func(t *testing.T) {
        _, err := dictionary.Search("unknown")
        want := "could not find the word you were looking for"

        if err == nil {
            t.Fatal("expected to get an error.")
        }

        assertStrings(t, err.Error(), want)
    })
}
```

The way to handle this scenario in Go is to return a second argument which is an `Error` type.

`Error`s can be converted to a string with the `.Error()` method, which we do when passing it to the assertion. We are also protecting `assertStrings` with `if` to ensure we don't call `.Error()` on `nil`.

## Write the minimal amount of code

```go
func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", errors.New("could not find the word you were looking for")
    }

    return definition, nil
}
```

In order to make this pass, we are using an interesting property of the map lookup. It can return 2 values. The second value is a boolean which indicates if the key was found successfully.

This property allows us to differentiate between a word that doesn't exist and a word that just doesn't have a definition.

## Refactoring

```go
var ErrNotFound = errors.New("could not find the word you were looking for")

func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", ErrNotFound
    }

    return definition, nil
}
```

We can get rid of the magic error in our Search function by extracting it into a variable.
