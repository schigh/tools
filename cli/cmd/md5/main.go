package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/schigh/str"
)

func main() {
	if len(os.Args) < 2 {
		randomMD5()
		return
	}

	s := str.MD5(os.Args[1])
	fmt.Println(s)
}

func randomMD5() {
	b := make([]byte, 1024)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	s := str.MD5(string(b))
	fmt.Println(s)
}

