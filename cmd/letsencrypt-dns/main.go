package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/helber/letsencrypt-dns/letsencrypt"
	"github.com/helber/letsencrypt-dns/linode"
)

func main() {
	linode.APIToken = os.Getenv("LINODE_API_KEY")
	domains := flag.String("d", "", "Domains sepered by \",\"")
	mainDomain := flag.String("domain", "", "Main domain")
	flag.Parse()
	domainlist := strings.Split(*domains, ",")

	if len(domainlist) == 1 {
		fmt.Printf("at last 1 domain is required %v given", len(domainlist)-1)
		return
	}
	if *mainDomain == "" {
		fmt.Printf("main domain required")
		return
	}
	// Done Channel
	done := make(chan bool)
	defer close(done)
	err := letsencrypt.CallAuto(domainlist, done)
	if err != nil {
		fmt.Printf("can't call letsencrypt: %s", err)
		return
	}
	result := <-done
	if result == true {
		fmt.Print("Congratulations")
	}
}
