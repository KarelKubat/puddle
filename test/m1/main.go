package main

import (
	"fmt"

	"github.com/KarelKubat/puddle"
)

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
	p := puddle.New(puddle.WithSize(4))

	for _, s := range []string{
		"one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten",
	} {
		p.Work(myPrintf, puddle.Args{"%s potato\n", s})
	}
	p.Wait()
}
