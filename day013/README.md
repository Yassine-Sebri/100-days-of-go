# Mocking (continued)

The tests still pass and the software works as intended but we have some problems:

- Our tests take 3 seconds to run.
  - Every forward-thinking post about software development emphasises the importance of quick feedback loops.
  - **Slow tests ruin developer productivity**.
  - Imagine if the requirements get more sophisticated warranting more tests. Are we happy with 3s added to the test run for every new test of `Countdown`?
- We have not tested an important property of our function.  

We have a dependency on `Sleep`ing which we need to extract so we can then control it in our tests.

If we can mock `time.Sleep` we can use dependency injection to use it instead of a "real" `time.Sleep` and then we can **spy on the calls** to make assertions on them.

## Write the test first

Let's define our dependency as an interface. This lets us then use a *real* `Sleeper` in `main` and a *spy* `Sleeper` in our tests. By using an interface our `Countdown` function is oblivious to this and adds some flexibility for the caller.

```go
type Sleeper interface {
 Sleep()
}
```

We made a design decision that our `Countdown` function would not be responsible for how long the sleep is. This simplifies our code a little for now at least and means a user of our function can configure that sleepiness however they like.

Now we need to make a mock of it for our tests to use.

```go
type SpySleeper struct {
 Calls int
}

func (s *SpySleeper) Sleep() {
 s.Calls++
}
```

Spies are a kind of mock which can record how a dependency is used. They can record the arguments sent in, how many times it has been called, etc. In our case, we're keeping track of how many times `Sleep()` is called so we can check it in our test.

Update the tests to inject a dependency on our Spy and assert that the sleep has been called 3 times.

```go
func TestCountdown(t *testing.T) {
 buffer := &bytes.Buffer{}
 spySleeper := &SpySleeper{}

 Countdown(buffer, spySleeper)

 got := buffer.String()
 want := `3
2
1
Go!`

 if got != want {
  t.Errorf("got %q want %q", got, want)
 }

 if spySleeper.Calls != 3 {
  t.Errorf("not enough calls to sleeper, want 3 got %d", spySleeper.Calls)
 }
}
```

## Write enough code to make it pass

```go
type DefaultSleeper struct{}

func (d *DefaultSleeper) Sleep() {
 time.Sleep(1 * time.Second)
}

func Countdown(out io.Writer, sleeper Sleeper) {
 for i := countdownStart; i > 0; i-- {
  fmt.Fprintln(out, i)
  sleeper.Sleep()
 }

 fmt.Fprint(out, finalWord)
}

func main() {
 sleeper := &DefaultSleeper{}
 Countdown(os.Stdout, sleeper)
}
```

The test should pass and no longer take 3 seconds.

## Still some problems

There's still another important property we haven't tested.

Our latest change only asserts that it has slept 3 times, but those sleeps could occur out of sequence.

Let's use spying again with a new test to check the order of operations is correct.

We have two different dependencies and we want to record all of their operations into one list. So we'll create one spy for them both.

```go
const (
 write = "write"
 sleep = "sleep"
)

type SpyCountdownOperations struct {
 Calls []string
}

func (s *SpyCountdownOperations) Sleep() {
 s.Calls = append(s.Calls, sleep)
}

func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
 s.Calls = append(s.Calls, write)
 return
}
```

Our `SpyCountdownOperations` implements both `io.Writer` and `Sleeper`, recording every call into one slice. In this test we're only concerned about the order of operations, so just recording them as list of named operations is sufficient.

We can now add a sub-test into our test suite which verifies our sleeps and prints operate in the order we hope

```go
func TestCountdown(t *testing.T) {
 t.Run("prints 3 to Go!", func(t *testing.T) {
  buffer := &bytes.Buffer{}
  Countdown(buffer, &SpyCountdownOperations{})

  got := buffer.String()
  want := `3
2
1
Go!`

  if got != want {
   t.Errorf("got %q want %q", got, want)
  }
 })

 t.Run("sleep before every print", func(t *testing.T) {
  spySleepPrinter := &SpyCountdownOperations{}
  Countdown(spySleepPrinter, spySleepPrinter)

  want := []string{
   write,
   sleep,
   write,
   sleep,
   write,
   sleep,
   write,
  }

  if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
   t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
  }
 })
}
```

We now have our function and its 2 important properties properly tested.
