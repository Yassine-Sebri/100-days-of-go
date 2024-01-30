# Context (Continued)

The approach we used yesterday is ok, but is it idiomatic?

Does it make sense for our web server to be concerned with manually cancelling `Store`? What if `Store` also happens to depend on other slow-running processes? We'll have to make sure that `Store.Cancel` correctly propagates the cancellation to all of its dependants.

Let's try and pass through the `context` to our `Store` and let it be responsible. That way it can also pass the `context` through to its dependants and they too can be responsible for stopping themselves.

## Write the test first

We'll have to change our existing tests as their responsibilities are changing. The only thing our handler is responsible for now is making sure it sends a context through to the downstream `Store` and that it handles the error that will come from the `Store` when it is cancelled.

Let's update our `Store` interface to show the new responsibilities.

```go
type Store interface {
 Fetch(ctx context.Context) (string, error)
}
```

Delete the code inside our handler for now

```go
func Server(store Store) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
 }
}
```

Update our `SpyStore`

```go
type SpyStore struct {
 response string
 t        *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
 data := make(chan string, 1)

 go func() {
  var result string
  for _, c := range s.response {
   select {
   case <-ctx.Done():
    log.Println("spy store got cancelled")
    return
   default:
    time.Sleep(10 * time.Millisecond)
    result += string(c)
   }
  }
  data <- result
 }()

 select {
 case <-ctx.Done():
  return "", ctx.Err()
 case res := <-data:
  return res, nil
 }
}
```

We have to make our spy act like a real method that works with context.

We are simulating a slow process where we build the result slowly by appending the string, character by character in a goroutine. When the goroutine finishes its work it writes the string to the `data` channel. The goroutine listens for the `ctx.Done` and will stop the work if a signal is sent in that channel.

Finally the code uses another `select` to wait for that goroutine to finish its work or for the cancellation to occur.

Finally we can update our tests. We comment out our cancellation test so we can fix the happy path test first.

```go
t.Run("returns data from store", func(t *testing.T) {
 data := "hello, world"
 store := &SpyStore{response: data, t: t}
 svr := Server(store)

 request := httptest.NewRequest(http.MethodGet, "/", nil)
 response := httptest.NewRecorder()

 svr.ServeHTTP(response, request)

 if response.Body.String() != data {
  t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
 }
})
```

## Write enough code to make it pass

```go
func Server(store Store) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  data, _ := store.Fetch(r.Context())
  fmt.Fprint(w, data)
 }
}
```

The test now passes.

## Write the test

We need to test that we do not write any kind of response on the error case. Sadly `httptest.ResponseRecorder` doesn't have a way of figuring this out so we'll have to roll our own spy to test for this.

```go
type SpyResponseWriter struct {
 written bool
}

func (s *SpyResponseWriter) Header() http.Header {
 s.written = true
 return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
 s.written = true
 return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
 s.written = true
}
```

Our `SpyResponseWriter` implements `http.ResponseWriter` so we can use it in the test.

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
 data := "hello, world"
 store := &SpyStore{response: data, t: t}
 svr := Server(store)

 request := httptest.NewRequest(http.MethodGet, "/", nil)

 cancellingCtx, cancel := context.WithCancel(request.Context())
 time.AfterFunc(5*time.Millisecond, cancel)
 request = request.WithContext(cancellingCtx)

 response := &SpyResponseWriter{}

 svr.ServeHTTP(response, request)

 if response.written {
  t.Error("a response should not have been written")
 }
})
```

## Write code to make it pass

```go
func Server(store Store) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  data, err := store.Fetch(r.Context())

  if err != nil {
   return // todo: log error however you like
  }

  fmt.Fprint(w, data)
 }
}
```

We can see after this that the server code has become simplified as it's no longer explicitly responsible for cancellation, it simply passes through `context` and relies on the downstream functions to respect any cancellations that may occur.
