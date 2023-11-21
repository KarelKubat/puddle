package main

import (
	"errors"
	"fmt"

	"github.com/KarelKubat/puddle"
)

type ret struct {
	n   int
	err error
}

func myPrintf(args puddle.Args) any {
	if len(args) < 1 {
		return ret{n: 0, err: errors.New("myPrintf: no fomat string")}
	}
	n, err := fmt.Printf(args[0].(string), args[1:]...)
	return ret{n: n, err: err}
}

func main() {
	p := puddle.New()
	p.Work(myPrintf, puddle.Args{"%v %v\n", "hello", "world"})
	p.Work(myPrintf, puddle.Args{"%v %v %v\n", "here", "I", "am"})
	for v := range p.Wait() {
		r := v.(ret)
		fmt.Printf("got n=%v, err=%v\n", r.n, r.err)
	}
}
