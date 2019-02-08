package main

import (
	"fmt"
	"log"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/helber/letsencrypt-dns/letsencrypt"
	mylog "github.com/helber/letsencrypt-dns/log"
)

func main() {
	domains := flag.StringP("domains", "d", "", "Domains sepered by \",\"")
	flag.Parse()
	mylog.InitLogs()

	domainlist := strings.Split(*domains, ",")
	main := ""
	for _, dom := range domainlist {
		domain := domainutil.Domain(dom)
		if domain == "" {
			fmt.Println("invalid domain")
			return
		}
		// log.Printf("Domain=%s", domain)
		if domain != main {
			if domain != "" {
				main = domain
			} else {
				fmt.Printf("multiple registers must be a same domain %v <> %v\n", main, domain)
				return
			}
		}
	}
	if main == "" {
		log.Printf("invalid domain")
		return
	}

	// Done Channel
	done := make(chan bool, 1)
	defer close(done)
	go letsencrypt.CallAuto(domainlist, done)
	result := <-done
	if result == true {
		log.Printf("Congratulations")
	}
}
