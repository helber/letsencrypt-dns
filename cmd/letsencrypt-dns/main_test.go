package main

import (
	"os"
)

func Example() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "-d", "lala,lele,lili,lolo,lulu"}
	main()
	// Output:
	// Generating cert for: [lala lele lili lolo lulu] domains
}
