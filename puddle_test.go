package puddle

import (
	"fmt"
	"sync/atomic"
	"testing"
)

type ret struct {
	s string
}

const MAX = 100

func workerForOut(a Args) any {
	if len(a) != 1 {
		panic(fmt.Sprintf("worker got %d args, need 1", len(a)))
	}
	msg := fmt.Sprintf("i saw a %d", a[0].(int))
	return ret{
		s: msg,
	}
}

var count atomic.Int32

func workerForWait(a Args) any {
	fmt.Printf("counter increased to %v", count.Add(1))
	return struct{}{}
}

func TestOut(t *testing.T) {
	p := New()
	for i := 0; i < MAX; i++ {
		p.Work(workerForOut, Args{i})
	}
	seen := map[string]struct{}{}
	for v := range p.Out() {
		r := v.(ret)
		seen[r.s] = struct{}{}
	}
	for i := 0; i < MAX; i++ {
		want := fmt.Sprintf("i saw a %d", i)
		if _, ok := seen[want]; !ok {
			t.Errorf("during TestOut: %q not seen", want)
		}
	}
}

func TestWait(t *testing.T) {
	p := New(Opts{Size: 1})
	for i := 0; i < MAX; i++ {
		p.Work(workerForWait, Args{})
	}
	p.Wait()
	if count.Load() != MAX {
		t.Errorf("during TestWait: counter reached %d, want %d", count.Load(), MAX)
	}
}
