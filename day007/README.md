# Pointers

Let's make a `Wallet` struct which lets us deposit Bitcoin.

## Write the test

```go
func TestWallet(t *testing.T) {
    wallet := Wallet{}

    wallet.Deposit(10)

    got := wallet.Balance()
    want := 10

    if got != want {
        t.Errorf("Expected %d but got %d", want, got)
    }
}
```

In the [previous day](day006) we accessed fields directly with the field name, however in our very *secure* wallet we don't want to expose our inner state to the rest of the world. We want to control access via methods.

## Write enough code to make it pass

We start by defining the struct and the methods.

```go
type Wallet struct{}

func (w Wallet) Deposit(amount int) {

}

func (w Wallet) Balance() int {
    return 0
}
```

We get

```text
wallet_test.go:14: Expected 10 but got 0
```

We will need some kind of `balance` variable in our struct to store the state

```go
type Wallet struct {
    balance int
}
```

In Go if a symbol (variables, types, functions et al) starts with a lowercase symbol then it is private outside the package it's defined in.

In our case we want our methods to be able to manipulate this value, but no one else.

```go
func (w Wallet) Deposit(amount int) {
    w.balance += amount
}

func (w Wallet) Balance() int {
    return w.balance
}
```

We now run the test but it still fails

```text
wallet_test.go:14: Expected 10 but got 0
```

## That's not quite right

Well this is confusing, our code looks like it should work. We add the new amount onto our balance and then the balance method should return the current state of it.

In Go, **when you call a function or a method the arguments are copied**.

When calling `func (w Wallet) Deposit(amount int)` the `w` is a copy of whatever we called the method from.

This means that when we change the value of the balance inside the code, we are working on a copy of what came from the test. Therefore the balance in the test is unchanged.

We can fix this with pointers. [Pointers](https://gobyexample.com/pointers) let us point to some values and then let us change them. So rather than taking a copy of the whole Wallet, we instead take a pointer to that wallet so that we can change the original values within it.

```go
func (w *Wallet) Deposit(amount int) {
    w.balance += amount
}

func (w *Wallet) Balance() int {
    return w.balance
}
```

The difference is the receiver type is `*Wallet` rather than `Wallet` which you can read as "a pointer to a wallet".

Now you might wonder, why did they pass? We didn't dereference the pointer in the function, like so:

```go
func (w *Wallet) Balance() int {
    return (*w).balance
}
```

and seemingly addressed the object directly. In fact, the code above using `(*w)` is absolutely valid. However, the makers of Go deemed this notation cumbersome, so the language permits us to write `w.balance`, without an explicit dereference. These pointers to structs even have their own name: *struct pointers* and they are [automatically dereferenced](https://golang.org/ref/spec#Method_values).

Technically we do not need to change `Balance` to use a pointer receiver as taking a copy of the balance is fine. However, by convention we should keep your method receiver types the same for consistency.

## Refactor

We said we were making a Bitcoin wallet but we have not mentioned them so far. We've been using `int` because they're a good type for counting things!

It seems a bit overkill to create a `struct` for this. `int` is fine in terms of the way it works but it's not descriptive.

Go lets us create new types from existing ones.

The syntax is `type MyName OriginalType`

```go
type Bitcoin int

type Wallet struct {
    balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
    w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
    return w.balance
}
```

To make `Bitcoin` you just use the syntax `Bitcoin(999)`.

By doing this we're making a new type and we can declare methods on them. This can be very useful when you want to add some domain specific functionality on top of existing types.

Let's implement [Stringer](https://pkg.go.dev/fmt#Stringer) on Bitcoin

```go
type Stringer interface {
    String() string
}
```

This interface is defined in the `fmt` package and lets us define how our type is printed when used with the `%s` format string in prints.

```go
func (b Bitcoin) String() string {
    return fmt.Sprintf("%d BTC", b)
}
```

Next we need to update our test format strings so they will use String() instead.

```go
if got != want {
    t.Errorf("got %s want %s", got, want)
}
```
