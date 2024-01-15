# Errors

Let's make continue from last day's code and add a `Withdraw` function.

## Write the test

```go
func TestWallet(t *testing.T) {

    t.Run("deposit", func(t *testing.T) {

        wallet := Wallet{}

        wallet.Deposit(Bitcoin(10))

        got := wallet.Balance()

        want := Bitcoin(10)

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    })

    t.Run("withdraw", func(t *testing.T) {

        wallet := Wallet{balance: Bitcoin(20)}

        wallet.Withdraw(Bitcoin(10))

        got := wallet.Balance()

        want := Bitcoin(10)

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    })
}
```

## Write enough code to make it pass

The method is pretty much the opposite of deposit

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
    w.balance -= amount
}
```

## Refactor

There's some duplication in our tests, lets refactor that out.

```go
func TestWallet(t *testing.T) {

    assertBalance := func(t testing.TB, wallet Wallet, want Bitcoin) {
        t.Helper()
        got := wallet.Balance()

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    }

    t.Run("deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("withdraw", func(t *testing.T) {
        wallet := Wallet{balance: Bitcoin(20)}
        wallet.Withdraw(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

}
```

What should happen if you try to `Withdraw` more than is left in the account? For now, our requirement is to assume there is not an overdraft facility.

How do we signal a problem when using `Withdraw`?

In Go, if you want to indicate an error it is idiomatic for your function to return an `err` for the caller to check and act on.

## Write the test first

```go
t.Run("withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, startingBalance)

    if err == nil {
        t.Error("wanted an error but didn't get one")
    }
})
```

We want `Withdraw` to return an error if you try to take out more than you have and the balance should stay the same.

We then check an error has returned by failing the test if it is `nil`.

`nil` is synonymous with `null` from other programming languages. Errors can be `nil` because the return type of `Withdraw` will be `error`, which is an interface. If you see a function that takes arguments or returns values that are interfaces, they can be nillable.

Like `null` if you try to access a value that is `nil` it will throw a **runtime panic**. This is bad! You should make sure that you check for `nil`s.

## Write the minimal amount of code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
    if amount > w.balance {
        return errors.New("you can't withdraw more than your balance")
    }

    w.balance -= amount
    return nil
}
```

## Refactoring

Let's make a quick test helper for our error check to improve the test's readability

```go
assertError := func(t testing.TB, err error) {
    t.Helper()
    if err == nil {
        t.Error("wanted an error but didn't get one")
    }
}
```

And in our test

```go
t.Run("withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertError(t, err)
    assertBalance(t, wallet, startingBalance)
})
```
