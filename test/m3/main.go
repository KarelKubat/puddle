package main

import (
	"fmt"
	"net/http"

	"github.com/KarelKubat/puddle"
)

type outcome struct {
	url string         // what we requested
	err error          // why http.Get() failed
	res *http.Response // what http.Get() returned otherwise
}

func httpGet(args puddle.Args) any {
	url := args[0].(string)
	r, e := http.Get(url)
	return outcome{
		url: url,
		res: r,
		err: e,
	}
}

func main() {
	p := puddle.New()
	urls := []string{
		"http://example.com",
		"http://example.com/non-existing-page",
		"https://go.dev/tour/welcome/1",
		"https://pkg.go.dev/net/http",
	}
	for _, u := range urls {
		p.Work(httpGet, puddle.Args{u})
	}
	for v := range p.Out() {
		o := v.(outcome)
		if o.err == nil {
			fmt.Printf("get %q returned status %d\n", o.url, o.res.StatusCode)
		} else {
			fmt.Printf("get %q returned error %v\n", o.url, o.err)
		}
	}
}
