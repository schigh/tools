package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func main() {
	times := 1
	if len(os.Args) > 1 {
		t, err := strconv.Atoi(os.Args[1])
		if err == nil && t > 0 {
			times = t
		}
	}

	for i := 0; i < times; i++ {
		out()
	}
}

func out() {
	u := uuid.New()
	fmt.Println(u.String())
}

