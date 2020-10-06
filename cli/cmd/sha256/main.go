package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/schigh/str"
)

func main() {
	if len(os.Args) < 2 {
		randomSHA256()
		return
	}

	s := str.SHA256(os.Args[1])
	fmt.Println(s)
}

func randomSHA256() {
	b := make([]byte, 1024)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	s := str.SHA256(string(b))
	fmt.Println(s)
}

