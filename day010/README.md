# Maps (Continued)

Let's create a way to add new words to our dictionary.

## Write the test first

```go
func TestAdd(t *testing.T) {
 dictionary := Dictionary{}
 dictionary.Add("test", "this is just a test")

 want := "this is just a test"
 got, err := dictionary.Search("test")
 if err != nil {
 t.Fatal("should find added word")
 }

assertStrings(t, got, want)
}
```

In this test, we are utilizing our `Search` function to make the validation of the dictionary a little easier.

## Write enough code to make it pass

```go
func (d Dictionary) Add(word, definition string) {
 d[word] = definition
}
```

Adding to a map is also similar to an array. We just need to specify a key and set it equal to a value.

## Maps and Pointers

An interesting property of maps is that you can modify them without passing as an address to it (e.g `&myMap`).

```text
A map value is a pointer to a runtime.hmap structure.
```

So when you pass a map to a function/method, you are indeed copying it, but just the pointer part, not the underlying data structure that contains the data.

A gotcha with maps is that they can be a `nil` value. A `nil` map behaves like an empty map when reading, but attempts to write to a `nil` map will cause a runtime panic.

Therefore, we should never initialize an empty map variable:

```go
var m map[string]string
```

Instead, we can initialize an empty map like we were doing above, or use the `make` keyword to create a map for us:

```go
var dictionary = map[string]string{}

// OR

var dictionary = make(map[string]string)
```

Both approaches create an empty `hash map` and point `dictionary` at it. Which ensures that we will never get a runtime panic.

## Refactor

We made variables for word and definition, and moved the definition assertion into its own helper function.

Our `Add` is looking good. Except, we didn't consider what happens when the value we are trying to add already exists!

Map will not throw an error if the value already exists. Instead, they will go ahead and overwrite the value with the newly provided value. This can be convenient in practice, but makes our function name less than accurate. `Add` should not modify existing values. It should only add new words to our dictionary.

## Write the test

```go
func TestAdd(t *testing.T) {
 t.Run("new word", func(t *testing.T) {
  dictionary := Dictionary{}
  word := "test"
  definition := "this is just a test"

  err := dictionary.Add(word, definition)

  assertError(t, err, nil)
  assertDefinition(t, dictionary, word, definition)
 })

 t.Run("existing word", func(t *testing.T) {
  word := "test"
  definition := "this is just a test"
  dictionary := Dictionary{word: definition}
  err := dictionary.Add(word, "new test")

  assertError(t, err, ErrWordExists)
  assertDefinition(t, dictionary, word, definition)
 })
}
```

## Write the minimal amount of code to make it pass

```go
var (
 ErrNotFound   = errors.New("could not find the word you were looking for")
 ErrWordExists = errors.New("cannot add word because it already exists")
)

func (d Dictionary) Add(word, definition string) error {
 _, err := d.Search(word)

 switch err {
 case ErrNotFound:
  d[word] = definition
 case nil:
  return ErrWordExists
 default:
  return err
 }

 return nil
}
```

Here we are using a `switch` statement to match on the error. Having a `switch` like this provides an extra safety net, in case `Search` returns an error other than `ErrNotFound`.

Next, let's create a function to Delete a word in the dictionary.

## Write the tests

```go
func TestDelete(t *testing.T) {
 word := "test"
 dictionary := Dictionary{word: "test definition"}

 dictionary.Delete(word)

 _, err := dictionary.Search(word)
 if err != ErrNotFound {
  t.Errorf("Expected %q to be deleted", word)
 }
}
```

## Write code to make it pass

```go
func (d Dictionary) Delete(word string) {
 delete(d, word)
}
```

Go has a built-in function `delete` that works on maps. It takes two arguments. The first is the map and the second is the key to be removed.

The `delete` function returns nothing, and we based our `Delete` method on the same notion. Since deleting a value that's not there has no effect, unlike our and `Add` method, we don't need to complicate the API with errors.
