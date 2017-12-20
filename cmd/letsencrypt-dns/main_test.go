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
	// main domain required
}

func Example_no_domains() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{""}
	main()
	// Output:
	// at last 1 domain is required 0 given
}

func Example_invalid_token() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "-domain", "ahgora.com.br", "-d", "a1.ahgora.com.br,a2.ahgora.com.br"}
	main()
	// Output:
	// can't call letsencrypt: linode api unathorized
}
