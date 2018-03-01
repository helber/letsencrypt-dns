package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/helber/letsencrypt-dns/checkcert"
	mylog "github.com/helber/letsencrypt-dns/log"
)

func main() {
	domains := flag.String("d", "", "Domain host and port (host:port) sepered by \",\"\nEx.: www.google.com.br:443,example.com:443,manage.openshift.com:443")
	flag.Parse()
	mylog.InitLogs()
	domainlist := strings.Split(*domains, ",")

	for _, dom := range domainlist {
		res, err := checkcert.CheckHost(dom)
		if err != nil {
			log.Printf("error %v", err)
		}
		// log.Printf("DOMAIN=%s days=%d", res.Host, res.ExpireDays)
		fmt.Println(res.ExpireDays)
	}
}
