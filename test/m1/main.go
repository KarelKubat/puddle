package main

import (
	"fmt"

	"github.com/KarelKubat/puddle"
)

// A worker must return something, even if it's empty.
type ret struct{}

// Wrapper for fmt.Printf.
func myPrintf(args puddle.Args) any {
	if len(args) > 1 {
		// Ensure that args[0] can be referenced.
		fmt.Printf(args[0].(string), args[1:]...)
	}
	return ret{}
}

// "puddle" example using "ret" and "myPrintf.
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
