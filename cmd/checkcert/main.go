package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/helber/letsencrypt-dns/checkcert"
	mylog "github.com/helber/letsencrypt-dns/log"
	flag "github.com/spf13/pflag"
)

func main() {
	domains := flag.StringP("domains", "d", "", "Domain host and port (host:port) sepered by \",\"\n\tEx.: www.google.com.br:443,example.com:443,manage.openshift.com:443")
	flag.Parse()
	mylog.InitLogs()
	domainlist := strings.Split(*domains, ",")

	for _, dom := range domainlist {
		res, err := checkcert.CheckHost(dom)
		if err != nil {
			log.Printf("error %v", err)
			log.Fatal("check return error")
		}
		log.Printf("DOMAIN=%s days=%d", res.Host, res.ExpireDays)
		fmt.Println(res.ExpireDays)
	}
}
