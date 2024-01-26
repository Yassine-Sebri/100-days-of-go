# Concurrency

Here's the setup: a colleague has written a function, `CheckWebsites`, that checks the status of a list of URLs.

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
 results := make(map[string]bool)

 for _, url := range urls {
  results[url] = wc(url)
 }

 return results
}
```

It returns a map of each URL checked to a boolean value: `true` for a good response; `false` for a bad response.

You also have to pass in a `WebsiteChecker` which takes a single URL and returns a boolean. This is used by the function to check all the websites.

Using dependency injection has allowed them to test the function without making real HTTP calls, making it reliable and fast.

Here's the test they've written:

```go
package concurrency

import (
 "reflect"
 "testing"
)

func mockWebsiteChecker(url string) bool {
 if url == "waat://furhurterwe.geds" {
  return false
 }
 return true
}

func TestCheckWebsites(t *testing.T) {
 websites := []string{
  "http://google.com",
  "http://blog.gypsydave5.com",
  "waat://furhurterwe.geds",
 }

 want := map[string]bool{
  "http://google.com":          true,
  "http://blog.gypsydave5.com": true,
  "waat://furhurterwe.geds":    false,
 }

 got := CheckWebsites(mockWebsiteChecker, websites)

 if !reflect.DeepEqual(want, got) {
  t.Fatalf("wanted %v, got %v", want, got)
 }
}
```

The function is in production and being used to check hundreds of websites. But your colleague has started to get complaints that it's slow, so they've asked you to help speed it up.

## Write a test

Let's use a benchmark to test the speed of `CheckWebsites` so that we can see the effect of our changes.

```go
package concurrency

import (
 "testing"
 "time"
)

func slowStubWebsiteChecker(_ string) bool {
 time.Sleep(20 * time.Millisecond)
 return true
}

func BenchmarkCheckWebsites(b *testing.B) {
 urls := make([]string, 100)
 for i := 0; i < len(urls); i++ {
  urls[i] = "a url"
 }
 b.ResetTimer()
 for i := 0; i < b.N; i++ {
  CheckWebsites(slowStubWebsiteChecker, urls)
 }
}
```

The benchmark tests `CheckWebsites` using a slice of one hundred urls and uses a new fake implementation of `WebsiteChecker`. `slowStubWebsiteChecker` is deliberately slow. It uses `time.Sleep` to wait exactly twenty milliseconds and then it returns true. We use `b.ResetTimer()` in this test to reset the time of our test before it actually runs.

When we run the benchmark using `go test -bench=`.

```plaintext
BenchmarkCheckWebsites-8               1        2048364674 ns/op
PASS
ok      day014  2.052s
```

`CheckWebsites` has been benchmarked at 2048364674 nanoseconds - about two seconds.

Let's try to make it faster.

## Write enough code to make it pass

Now we can finally talk about concurrency which, for the purposes of the following, means "having more than one thing in progress." This is something that we do naturally everyday.

Normally in Go when we call a function `doSomething()` we wait for it to return (even if it has no value to return, we still wait for it to finish). We say that this operation is blocking - it makes us wait for it to finish. An operation that does not block in Go will run in a separate process called a goroutine. Think of a process as reading down the page of Go code from top to bottom, going 'inside' each function when it gets called to read what it does. When a separate process starts, it's like another reader begins reading inside the function, leaving the original reader to carry on going down the page.

To tell Go to start a new goroutine we turn a function call into a go statement by putting the keyword go in front of it: `go doSomething()`.

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
 results := make(map[string]bool)

 for _, url := range urls {
  go func() {
   results[url] = wc(url)
  }()
 }

 return results
}
```

Because the only way to start a goroutine is to put `go` in front of a function call, we often use anonymous functions when we want to start a goroutine. An anonymous function literal looks just the same as a normal function declaration, but without a name (unsurprisingly). You can see one above in the body of the `for` loop.

Anonymous functions have a number of features which make them useful, two of which we're using above. Firstly, they can be executed at the same time that they're declared - this is what the `()` at the end of the anonymous function is doing. Secondly they maintain access to the lexical scope in which they are defined - all the variables that are available at the point when you declare the anonymous function are also available in the body of the function.

The body of the anonymous function above is just the same as the loop body was before. The only difference is that each iteration of the loop will start a new goroutine, concurrent with the current process (the `WebsiteChecker` function). Each goroutine will add its result to the results map.

Now we run `go test`

```plaintext
--- FAIL: TestCheckWebsites (0.00s)
    CheckWebsites_test.go:32: wanted map[http://blog.gypsydave5.com:true http://google.com:true waat://furhurterwe.geds:false], got map[]
FAIL
exit status 1
FAIL    day014  0.001s
```

## A quick aside into a parallel(ism) universe

We are caught by the original test `CheckWebsites`, it's now returning an empty map. What went wrong?

None of the goroutines that our `for` loop started had enough time to add their `result` to the results map; the `WebsiteChecker` function is too fast for them, and it returns the still empty map.

## Channels

Channels are a Go data structure that can both receive and send values. These operations, along with their details, allow communication between different processes.

In this case we want to think about the communication between the parent process and each of the goroutines that it makes to do the work of running the `WebsiteChecker` function with the url.

```go
package concurrency

type WebsiteChecker func(string) bool
type result struct {
 string
 bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
 results := make(map[string]bool)
 resultChannel := make(chan result)

 for _, url := range urls {
  go func(u string) {
   resultChannel <- result{u, wc(u)}
  }(url)
 }

 for i := 0; i < len(urls); i++ {
  r := <-resultChannel
  results[r.string] = r.bool
 }

 return results
}
```

Alongside the `results` map we now have a `resultChannel`, which we make in the same way. `chan result` is the type of the channel - a channel of `result`. The new type, `result` has been made to associate the return value of the `WebsiteChecker` with the url being checked - it's a `struct` of `string` and `bool`. As we don't need either value to be named, each of them is anonymous within the `struct`; this can be useful in when it's hard to know what to name a value.

Now when we iterate over the urls, instead of writing to the `map` directly we're sending a `result` struct for each call to `wc` to the `resultChannel` with a send statement. This uses the `<-` operator, taking a channel on the left and a value on the right:

```go
// Send statement
resultChannel <- result{u, wc(u)}
```

The next `for` loop iterates once for each of the urls. Inside we're using a receive expression, which assigns a value received from a channel to a variable. This also uses the `<-` operator, but with the two operands now reversed: the channel is now on the right and the variable that we're assigning to is on the left:

```go
// Receive expression
r := <-resultChannel
```

We then use the `result` received to update the map.

By sending the results into a channel, we can control the timing of each write into the results map, ensuring that it happens one at a time. Although each of the calls of `wc`, and each send to the result channel, is happening in parallel inside its own process, each of the results is being dealt with one at a time as we take values out of the result channel with the receive expression.

We have parallelized the part of the code that we wanted to make faster, while making sure that the part that cannot happen in parallel still happens linearly. And we have communicated across the multiple processes involved by using channels.

When we run the benchmark:

```plaintext
BenchmarkCheckWebsites-8              56          20835050 ns/op
PASS
ok      day014  1.194
```

20835050 nanoseconds, or 0.02 seconds, about one hundred times faster than the original function.
