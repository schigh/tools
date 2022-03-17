package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	prefix string
	secret string
	suffix string
	cost   string
)

func main() {
	flag.StringVar(&secret, "secret", "", "the secret to be encrypted")
	flag.StringVar(&prefix, "prefix", "", "prefix applied to secret before encryption")
	flag.StringVar(&suffix, "suffix", "", "suffix applied to secret before encryption")
	flag.StringVar(&cost, "cost", "default", "encryption cost (min|max|default)")
	flag.Parse()

	if secret == "" {
		log.Fatalln("secret is required")
	}

	var bCost int
	switch cost {
	case "min":
		bCost = bcrypt.MinCost
	case "max":
		bCost = bcrypt.MaxCost
	case "default":
		bCost = bcrypt.DefaultCost
	default:
		log.Fatalf("'%s' is an invalid cost. Use 'min', 'max' or omit the flag to use default", cost)
	}

	raw := []byte(strings.TrimSpace(prefix + secret + suffix))
	hash, err := bcrypt.GenerateFromPassword(raw, bCost)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(hash))
}
