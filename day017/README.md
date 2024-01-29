# Context

In this chapter we'll use the package context to help us manage long-running processes.

We're going to start with a classic example of a web server that when hit kicks off a potentially long-running process to fetch some data for it to return in the response.

We will exercise a scenario where a user cancels the request before the data can be retrieved and we'll make sure the process is told to give up.

```go
type Store interface {
 Fetch() string
}

func Server(store Store) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, store.Fetch())
 }
}
```

```go
type StubStore struct {
 response string
}

func (s *StubStore) Fetch() string {
 return s.response
}

func TestServer(t *testing.T) {
 data := "hello, world"
 svr := Server(&StubStore{data})

 request := httptest.NewRequest(http.MethodGet, "/", nil)
 response := httptest.NewRecorder()

 svr.ServeHTTP(response, request)

 if response.Body.String() != data {
  t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
 }
}
```

## Write the test first

Our handler will need a way of telling the `Store` to cancel the work

```go
type Store interface {
 Fetch() string
 Cancel()
}
```

We will need to adjust our spy so it takes some time to return data and a way of knowing it has been told to cancel. We'll also rename it to `SpyStore` as we are now observing the way it is called. It'll have to add `Cancel` as a method to implement the `Store` interface.

```go
type SpyStore struct {
 response  string
 cancelled bool
}

func (s *SpyStore) Fetch() string {
 time.Sleep(100 * time.Millisecond)
 return s.response
}

func (s *SpyStore) Cancel() {
 s.cancelled = true
}
```

Let's add a new test where we cancel the request before 100 milliseconds and check the store to see if it gets cancelled

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
 data := "hello, world"
 store := &SpyStore{response: data}
 svr := Server(store)

 request := httptest.NewRequest(http.MethodGet, "/", nil)

 cancellingCtx, cancel := context.WithCancel(request.Context())
 time.AfterFunc(5*time.Millisecond, cancel)
 request = request.WithContext(cancellingCtx)

 response := httptest.NewRecorder()

 svr.ServeHTTP(response, request)

 if !store.cancelled {
  t.Error("store was not told to cancel")
 }
})
```

The context package provides functions to derive new Context values from existing ones. These values form a tree: when a Context is canceled, all Contexts derived from it are also canceled.

What we do is derive a new `cancellingCtx` from our request which returns us a `cancel` function. We then schedule that function to be called in 5 milliseconds by using `time.AfterFunc`. Finally we use this new context in our request by calling `request.WithContext`.

## Write enough code to make it pass

```go
func Server(store Store) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  store.Cancel()
  fmt.Fprint(w, store.Fetch())
 }
}
```

This code passes the test but it doesn't achieve our goal so let's update the tests.

```go
t.Run("returns data from store", func(t *testing.T) {
 data := "hello, world"
 store := &SpyStore{response: data}
 svr := Server(store)

 request := httptest.NewRequest(http.MethodGet, "/", nil)
 response := httptest.NewRecorder()

 svr.ServeHTTP(response, request)

 if response.Body.String() != data {
  t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
 }

 if store.cancelled {
  t.Error("it should not have cancelled the store")
 }
})
```

Now we're forced to write a more sensible implementation

```go
func Server(store Store) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  ctx := r.Context()

  data := make(chan string, 1)

  go func() {
   data <- store.Fetch()
  }()

  select {
  case d := <-data:
   fmt.Fprint(w, d)
  case <-ctx.Done():
   store.Cancel()
  }
 }
}
```

`context` has a method `Done()` which returns a channel which gets sent a signal when the context is "done" or "cancelled". We want to listen to that signal and call `store.Cancel` if we get it but we want to ignore it if our `Store` manages to `Fetch` before it.

To manage this we run `Fetch` in a goroutine and it will write the result into a new channel `data`. We then use `select` to effectively race to the two asynchronous processes and then we either write a response or `Cancel`.

## Refactor

We can refactor our test code a bit by making assertion methods on our spy

```go
type SpyStore struct {
 response  string
 cancelled bool
 t         *testing.T
}

func (s *SpyStore) assertWasCancelled() {
 s.t.Helper()
 if !s.cancelled {
  s.t.Error("store was not told to cancel")
 }
}

func (s *SpyStore) assertWasNotCancelled() {
 s.t.Helper()
 if s.cancelled {
  s.t.Error("store was told to cancel")
 }
}
```

```go
func TestServer(t *testing.T) {
 data := "hello, world"

 t.Run("returns data from store", func(t *testing.T) {
  store := &SpyStore{response: data, t: t}
  svr := Server(store)

  request := httptest.NewRequest(http.MethodGet, "/", nil)
  response := httptest.NewRecorder()

  svr.ServeHTTP(response, request)

  if response.Body.String() != data {
   t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
  }

  store.assertWasNotCancelled()
 })

 t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
  store := &SpyStore{response: data, t: t}
  svr := Server(store)

  request := httptest.NewRequest(http.MethodGet, "/", nil)

  cancellingCtx, cancel := context.WithCancel(request.Context())
  time.AfterFunc(5*time.Millisecond, cancel)
  request = request.WithContext(cancellingCtx)

  response := httptest.NewRecorder()

  svr.ServeHTTP(response, request)

  store.assertWasCancelled()
 })
}
```
