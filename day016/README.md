# Sync

We want to make a counter which is safe to use concurrently.

We'll start with an unsafe counter and verify its behaviour works in a single-threaded environment.

Then we'll exercise it's unsafeness, with multiple goroutines trying to use the counter via a test, and fix it.

## Write the first test first

We want our API to give us a method to increment the counter and then retrieve its value.

```go
func TestCounter(t *testing.T) {
 t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
  counter := Counter{}
  counter.Inc()
  counter.Inc()
  counter.Inc()

  if counter.Value() != 3 {
   t.Errorf("got %d, want %d", counter.Value(), 3)
  }
 })
}
```

## Write enough code to make it pass

```go
type Counter struct {
 value int
}

func (c *Counter) Inc() {
 c.value++
}

func (c *Counter) Value() int {
 return c.value
}
```

## Refactor

```go
t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
 counter := Counter{}
 counter.Inc()
 counter.Inc()
 counter.Inc()

 assertCounter(t, counter, 3)
})
```

```go
func assertCounter(t testing.TB, got Counter, want int) {
 t.Helper()
 if got.Value() != want {
  t.Errorf("got %d, want %d", got.Value(), want)
 }
}
```

## Next steps

That was easy enough but now we have a requirement that it must be safe to use in a concurrent environment.

## Write the test first

```go
t.Run("it runs safely concurrently", func(t *testing.T) {
 wantedCount := 1000
 counter := Counter{}

 var wg sync.WaitGroup
 wg.Add(wantedCount)

 for i := 0; i < wantedCount; i++ {
  go func() {
   counter.Inc()
   wg.Done()
  }()
 }
 wg.Wait()

 assertCounter(t, counter, wantedCount)
})
```

This will loop through our `wantedCount` and fire a goroutine to call `counter.Inc()`.

We are using `sync.WaitGroup` which is a convenient way of synchronizing concurrent processes.

By waiting for `wg.Wait()` to finish before making our assertions we can be sure all of our goroutines have attempted to `Inc` the `Counter`.

## Try to run the test

The test fails

```plaintext
got 941, want 1000
```

This demonstrates it does not work when multiple goroutines are trying to mutate the value of the counter at the same time.

## Make it pass

A simple solution is to add a lock to our `Counter`, ensuring only one goroutine can increment the counter at a time. Go's Mutex provides such a lock.

```go
type Counter struct {
 mu    sync.Mutex
 value int
}

func (c *Counter) Inc() {
 c.mu.Lock()
 defer c.mu.Unlock()
 c.value++
}
```

What this means is any goroutine calling `Inc` will acquire the lock on `Counter` if they are first. All the other goroutines will have to wait for it to be Unlocked before getting access.

## Copying mutexes

Our test passes but our code is still a bit dangerous.

We can run `go vet` to check why.

To solve this we should pass in a pointer to our `Counter` instead, so change the signature of `assertCounter`.

```go
func assertCounter(t testing.TB, got *Counter, want int)
```

Now we have to create a function to initialize the type.

```go
func NewCounter() *Counter {
 return &Counter{}
}
```
