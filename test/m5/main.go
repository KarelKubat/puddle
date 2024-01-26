package main

import (
	"fmt"
	"os/exec"

	"github.com/KarelKubat/puddle"
)

func main() {
	hostnames := []string{
		"google.com",
		"example.com",
		"non-existent.whatever",
	}

	worker := func(args puddle.Args) any {
		hostname := args[0].(string)
		cmd := exec.Command("ping", hostname)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%q: %v", hostname, err)
		}
		return nil
	}

	p := puddle.New()
	for _, h := range hostnames {
		p.Work(worker, puddle.Args{h})
	}
	for v := range p.Out() {
		if v != nil {
			fmt.Println(v)
		}
	}
}
