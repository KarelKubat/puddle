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
	if len(a) != 0 {
		panic(fmt.Sprintf("worker got %d args, need 0", len(a)))
	}
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
		p.Work(workerForWait, nil)
	}
	p.Wait()
	if count.Load() != MAX {
		t.Errorf("during TestWait: counter reached %d, want %d", count.Load(), MAX)
	}
}

func TestNewSizing(t *testing.T) {
	for _, test := range []struct {
		desc   string
		p      *Pool
		wantSz int
	}{
		{
			desc:   "New()",
			p:      New(),
			wantSz: 0,
		},
		{
			desc:   "New(Opts{Size:20})",
			p:      New(Opts{Size: 20}),
			wantSz: 20,
		},
		{
			desc:   "New(WithSize(30))",
			p:      New(WithSize(30)),
			wantSz: 30,
		},
	} {
		if gotSz := test.p.size; gotSz != test.wantSz {
			t.Errorf("%v: size=%v, want %v", test.desc, gotSz, test.wantSz)
		}
	}
}
