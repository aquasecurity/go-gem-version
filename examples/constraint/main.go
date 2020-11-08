package main

import (
	"fmt"
	"log"

	"github.com/aquasecurity/go-gem-version"
)

func main() {
	v, err := gem.NewVersion("2.1")
	if err != nil {
		log.Fatal(err)
	}

	c, err := gem.NewConstraints(">= 1.0, < 1.4 || > 2.0")
	if err != nil {
		log.Fatal(err)
	}

	if c.Check(v) {
		fmt.Printf("%s satisfies constraints '%s'", v, c)
	}
}
