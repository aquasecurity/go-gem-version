package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/aquasecurity/go-gem-version"
)

func main() {
	versionsRaw := []string{"1.1", "0.7.1", "1.4.a", "1.4.a.1", "1.4", "1.4.0.1"}
	versions := make([]gem.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, err := gem.NewVersion(raw)
		if err != nil {
			log.Fatal(err)
		}
		versions[i] = v
	}

	// After this, the versions are properly sorted
	sort.Sort(gem.Collection(versions))

	for _, v := range versions {
		fmt.Println(v)
	}
}
