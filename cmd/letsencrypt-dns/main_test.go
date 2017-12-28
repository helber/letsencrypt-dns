package main

import (
	"os"
)

func Example_no_domain() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "-d", "lala,lele,lili,lolo,lulu"}
	main()
	// Output:
	// invalid domain
}

func Example_no_domains() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{""}
	main()
	// Output:
	// invalid domain
}

func Example_invalid_token() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "-d", "a1.example.com.br,a2.example.com"}
	main()
	// Output:
	// multiple registers must be a same domain example.com.br <> example.com
}
