package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	var (
		times int
		seed string
	)

	flag.IntVar(&times, "n", 1, "number of ids to generate")
	flag.StringVar(&seed, "seed", time.Now().Format("2006-01-02"), "seed for uuid")

	for i := 0; i < times; i++ {
		out(seed)
	}
}

func out(seed string) {
	uuid.SetNodeID([]byte(seed))
	u, err := uuid.NewUUID()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(u.String())
}

