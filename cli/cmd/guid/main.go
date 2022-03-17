package main

import (
	"flag"
	"fmt"

	"github.com/ntwrk1/guid"
)

var (
	prefix string
)

func main() {
	flag.StringVar(&prefix, "prefix", "", "byte prefix")
	flag.Parse()

	if len(prefix) >= 2 {
		guid.SetGlobalPrefixBytes(prefix[0], prefix[1])
	}

	g := guid.New()

	fmt.Println(g.String())
}
