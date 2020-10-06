package main

import (
	"fmt"
	"os"

	"github.com/schigh/tools/pkg/cuid"
)

func main() {
	c, err := cuid.New()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(c.String())
}
