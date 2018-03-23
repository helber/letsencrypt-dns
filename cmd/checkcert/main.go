package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/helber/letsencrypt-dns/checkcert"
	mylog "github.com/helber/letsencrypt-dns/log"
	flag "github.com/spf13/pflag"
)

var wg sync.WaitGroup

func main() {
	domains := flag.StringP("domains", "d", "", "Domain host and port (host:port) sepered by \",\"\n\tEx.: www.google.com.br:443,example.com:443,manage.openshift.com:443")
	showHostTime := flag.BoolP("hosttime", "t", false, "Display host and elapsed query time in a table")
	flag.Parse()
	mylog.InitLogs()
	domainlist := strings.Split(*domains, ",")
	results := make(chan checkcert.HostResult)

	for _, dom := range domainlist {
		wg.Add(1)
		go checkcert.CheckHost(dom, results)
	}
	if *showHostTime {
		fmt.Printf("+%s+%s+%s+\n", strings.Repeat("-", 15), strings.Repeat("-", 15), strings.Repeat("-", 62))
		fmt.Printf("| %-14v | %-14v | %-58v |\n", "query time", "expire days", "host:port")
		fmt.Printf("+%s+%s+%s+\n", strings.Repeat("-", 15), strings.Repeat("-", 15), strings.Repeat("-", 62))
	}
	for _, dom := range domainlist {
		res := <-results
		wg.Done()
		log.Printf("DOMAIN=%s days=%d version=%v by=%v", res.Host, res.ExpireDays, res.TLSVersion, res.Issuer)
		log.Printf("Domain=%s", dom)
		if *showHostTime {
			fmt.Printf("| %-14v | %-14v | %-58v |\n", res.ElapsedTime, res.ExpireDays, res.Host)
		} else {
			fmt.Println(res.ExpireDays)
		}
	}
	wg.Wait()
	if *showHostTime {
		fmt.Printf("+%s+%s+%s+\n", strings.Repeat("-", 15), strings.Repeat("-", 15), strings.Repeat("-", 62))
	}
}
