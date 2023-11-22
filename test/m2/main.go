package main

import (
	"errors"
	"fmt"

	"github.com/KarelKubat/puddle"
)

// What we are interested as the return value, which is what fmt.Printf returns.
// We need to wrap the two values into one struct.
type ret struct {
	n   int
	err error
}

// Wrapper for fmt.Printf.
func myPrintf(args puddle.Args) any {
	if len(args) < 1 {
		return ret{n: 0, err: errors.New("myPrintf: no fomat string")}
	}
	n, err := fmt.Printf(args[0].(string), args[1:]...)
	return ret{n: n, err: err}
}

// Wrapper for fmt.Pritnln.
func myPrintln(args puddle.Args) any {
	n, err := fmt.Println(args...)
	return ret{n: n, err: err}
}

// "puddle" example using "ret" and "nyPrintf / myPrintln".
func main() {
	p := puddle.New()

	for _, s := range []string{
		"one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten",
	} {
		p.Work(myPrintf, puddle.Args{"%s potato\n", s})
	}
	for v := range p.Out() {
		r := v.(ret)
		fmt.Printf("myPrintf returned: n=%v, err=%v\n", r.n, r.err)
	}

	p.Work(myPrintln, puddle.Args{"And now I am", "done"})
	for v := range p.Out() {
		r := v.(ret)
		fmt.Printf("myPrintln returned: n=%v, err=%v\n", r.n, r.err)
	}
}
