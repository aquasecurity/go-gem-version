package main

import (
	"fmt"
	"log"

	"github.com/aquasecurity/go-gem-version"
)

func main() {
	v1, err := gem.NewVersion("1.2.a")
	if err != nil {
		log.Fatal(err)
	}

	v2, err := gem.NewVersion("1.2")
	if err != nil {
		log.Fatal(err)
	}

	// Comparison example. There is also GreaterThan, Equal, and just
	// a simple Compare that returns an int allowing easy >=, <=, etc.
	if v1.LessThan(v2) {
		fmt.Printf("%s is less than %s", v1, v2)
	}
}
