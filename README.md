# puddle: A basic worker pool for Go

<!-- toc -->
- [Pool creation](#pool-creation)
- [Starting work](#starting-work)
- [Waiting for termination](#waiting-for-termination)
  - [Collecting results](#collecting-results)
  - [Discarding results](#discarding-results)
- [Chaining](#chaining)
- [Examples](#examples)
<!-- /toc -->

The package `puddle` is an abstraction of a worker pool that may fit most many cases. At least, it fits mine and I don't have to remember how channels and waitgroups work and what to wrap in a go-routine and what not.

## Pool creation

```go
import "github.com/KarelKubat/puddle"

// When all workers in the pool can run in parallel
p := puddle.New()

// When there may be at most 20 parallel executions at one time
p := puddle.New(puddle.WithSize(20))
```

## Starting work

The pool accepts `puddle.Worker` type functions. Such functions:

- Accept as arguments a list of `any` (this is the type `puddle.Args`)
- Return one value, again `any`.

Usually you'll want to coerce existing functions into the `puddle.Worker` format. This can be easily done by (a) wrapping the returns into one `struct`, (b) wrapping the existing function to accept `puddle.Args` as the argument and to return one value, that struct.

For example, `http.Get()` accepts a string (the URL) and returns a `*http.Response` and an error. A wrapper is:

```go
type outcome struct {
    res  *http.Response
    err  error
}
func httpGet(args puddle.Args) any {
    url := args[0].(string)
    r, e := http.Get(url)
    return outcome{
        res: r,
        err: e,
    }
}
```

A `puddle.Worker` **must** return something, even if it's only a `struct{}`.

Adding workers to the pool and starting them is done using `p.Work()` which has as arguments (a) the function, (b) what the function will receive once it runs, in the format `puddle.Args`. For example:

```go
urls := []string{
    "http://example.com/page/about/this",
    "http://example.com/page/about/that",
    "http://example.com/page/about/something-else",
    // etc.
}
for _, u := range urls {
    p.Add(httpGet, puddle.Args{u})
}
```

## Waiting for termination

There are two ways to wait until the pool's workers finish their work:

1. `p.Wait()` which lets the workers finish but discards the results,
1. `p.Out()` which returns a channel to be consumed in a `range` loop. The loop unblocks once all workers have finished.

### Collecting results

Given the above example about `http.Get()` we would use `p.Out()` and inspect what happened. The outcome is returned as an anonymous `any`, the caller must convert it to whatever the worker returns (in this case an `outcome`).

```go
for v := range p.Out() {
	o := v.(outcome)
	if o.err == nil {
		fmt.Printf("worker returned status %d\n", o.res.StatusCode)
		// Presumably here we'd do something with o.res.Body
	} else {
		fmt.Printf("worker returned error %v\n", o.err)
	}
}
```

### Discarding results

FWIW, `p.Wait()` can be called to wait until all workers have finished. Here is a trivial example:

```go
// Wrapper for fmt.Printf.
func myPrintf(args puddle.Args) any {
	if len(args) > 1 {
		// Ensure that args[0] can be referenced.
		fmt.Printf(args[0].(string), args[1:]...)
	}
	// There must be a return value, even when no one will inspect it.
	return nil
}

// Puddle example using myPrintf. Since we don't want to collect the
// results, we can p.Wait() which just blocks until all workers finish.
func main() {
	p := puddle.New()

	for _, s := range []string{
		"one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten",
	} {
		p.Work(myPrintf, puddle.Args{"%s potato\n", s})
	}
	p.Wait()
}
```

## Chaining

The outcome of one pool can of course start workers in another pool. The below code has a pool `formatter` that spits out strings, and a pool `outputter` that displays them.

- Given that `formatter` workers emit strings, waiting for the `formatter` needs to collect results. Hence `formatter.Out()` is applied.
- Waiting for the `outputter` can be just `outputter.Wait()` since there are no results to collect.
- The code also shows how a lambda function can be a wrapper for `.Work()`.

```go
func main() {
	formatter := puddle.New()
	for _, s := range []string{
		"one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten",
	} {
		formatter.Work(func(a puddle.Args) any {
			return fmt.Sprintf(a[0].(string), a[1:]...)
		}, puddle.Args{"%s potato", s})
	}

	outputter := puddle.New()
	for v := range formatter.Out() {
		s := v.(string)
		outputter.Work(func(a puddle.Args) any {
			fmt.Println(a[0].(string))
			return nil
		}, puddle.Args{s})
	}
	outputter.Wait()
}
```

## Examples

You'll find working examples under `test/` in this distribution.