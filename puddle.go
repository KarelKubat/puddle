// Package puddle is a somewhat simplified worker pool for Go, but it should serve lots of use cases.
package puddle

import (
	"sync"
)

// Worker is the type of workers to attach to a puddle: they must accept Args, and return an
// anonymous interface type.
type Worker func(Args) any

// Args are the arguments to any Worker to attach: a list of anonymous interfaces.
type Args []any

// Pool is the shallow puddle.
type Pool struct {
	wg   sync.WaitGroup
	ch   chan any
	size int
	open bool
}

// Opts are options for New.
type Opts struct {
	Size int
}

/*
New returns a new Pool. Examples:

	// Unlimited parallel workers.
	pl := puddle.New()

	// Limit the workers to 10 at a time.
	pl := puddle.New(puddle.Opts{Size: 10})
*/
func New(opts ...Opts) *Pool {
	wp := &Pool{}
	for _, o := range opts {
		wp.size = o.Size
	}
	wp.openChannel()
	return wp
}

// openChannel is a helper.
func (wp *Pool) openChannel() {
	var mu sync.Mutex

	mu.Lock()
	defer mu.Unlock()
	if wp.open {
		return
	}
	if wp.size > 0 {
		wp.ch = make(chan any, wp.size)
	} else {
		wp.ch = make(chan any)
	}
	wp.open = true
}

/*
Work adds a worker to the puddle. The worker must (a) accept anything of type Args, which is a list
of anonymous interfaces, and (b) it must return one anonymous any interface. Existing functions can
be easily wrapped. Below is an example for fmt.Printf(), which returns the number of written bytes
and an error.

	// fmt.Printf's return values in one struct.
	type ret struct {
		n   int
		err error
	}
	// Wrapper for fmt.Printf() to emit its return values as a ret type.
	func myPrintf(args puddle.Args) any {
		// Ensure that the below args[0] doesn't crash.
		if len(args) == 0 {
			return 0, errors.New("myPrintf called without arguments")
		}
		// Call fmt.Printf, return outcome as one value.
		n, err := fmt.Printf(args[0].(string), args[1:]...)
		return ret{n: n, err: err}
	}
*/
func (wp *Pool) Work(fn Worker, args Args) {
	wp.openChannel()
	wp.wg.Add(1)
	go func(args Args) {
		defer wp.wg.Done()
		wp.ch <- fn(args)
	}(args)
}

/*
Out returns a channel of an "any" interface, which is also the return type of a Worker.
The caller may consume the channel in a "for range" statement which unblocks when
all workers have finished.

Note that in situations where the workers' results are not relevant, Wait can be used.

The type that the channel returns is an anonymous "any". The caller must ensure that this
is converted to whatever a worker function (type "Worker") returns. Example:

	// Hypothetical return from "myFunc"
	type ret struct {
		err error  // incase it went wrong
		timing time.Duraton  // how long it took
	}

	// Hypothetical "myFunc" conforming to puddle.Worker.
	func myFunc(args puddle.Args) any {
		// fill me in
		return ret{err: nil, timing: time.Second}
	}

	func main() {
		p := puddle.New(puddle.Opts{Size: 200})  // no more than 200 parallel workers
		for i := 0; i < 100000 {
			p.Work(myFunc, puddle.Args{i})
		}
		// Collect the outcomes.
		for v := range p.Out() {
			r := v.(ret)
			if r.err != nil {
				fmt.Fprintln(os.Stderr, r.err)
			} else {
				fmt.Printf("processed in %v\n", r.timing)
			}
		}
	}
*/
func (wp *Pool) Out() chan any {
	go func() {
		wp.wg.Wait()
		close(wp.ch)
		wp.open = false
	}()
	return wp.ch
}

// Wait blocks until all workers in the puddle have finished and discards the outcomes. It can be
// used instead of Out in situations where the workers' returns are not relevant.
func (wp *Pool) Wait() {
	for range wp.Out() {
	}
}
