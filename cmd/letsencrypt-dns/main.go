package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/helber/letsencrypt-dns/letsencrypt"
	"github.com/helber/letsencrypt-dns/linode"
)

func main() {
	linode.APIToken = os.Getenv("LINODE_API_KEY")
	domains := flag.String("d", "", "Domains sepered by \",\"")
	flag.Parse()

	domainlist := strings.Split(*domains, ",")
	main := ""
	for _, dom := range domainlist {
		domain := domainutil.Domain(dom)
		// log.Println(domain)
		if domain != main {
			if main == "" {
				main = domain
			} else {
				fmt.Printf("multiple registers must be a same domain %v <> %v\n", main, domain)
				return
			}
		}
	}
	if main == "" {
		fmt.Println("invalid domain")
		return
	}
	// Done Channel
	done := make(chan bool, 1)
	defer close(done)
	go letsencrypt.CallAuto(domainlist, done)
	result := <-done
	if result == true {
		fmt.Print("Congratulations")
	}
}
