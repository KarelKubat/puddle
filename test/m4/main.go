package main

import (
	"fmt"

	"github.com/KarelKubat/puddle"
)

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
