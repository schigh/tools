package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	defaultLength = 8
	charset       = `aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789`
)

var (
	charsetLength = len(charset)
	minLength     = float64(1)
	maxLength     = float64(100)
)

func main() {
	length := defaultLength

	if len(os.Args) > 1 {
		v, err := strconv.Atoi(os.Args[1])
		if err == nil {
			length = int(math.Max(math.Min(float64(v), maxLength), minLength))
		}
	}

	data := make([]byte, length)
	if _, err := rand.Read(data); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	out := make([]byte, length)
	for i, b := range data {
		out[i] = charset[int(b)%charsetLength]
	}

	fmt.Println(string(out))
}
