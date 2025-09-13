package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"strings"

	"github.com/helber/letsencrypt-dns/checkcert"

	mylog "github.com/helber/letsencrypt-dns/log"
	flag "github.com/spf13/pflag"
)

func main() {
	domains := flag.StringP("domains", "d", "", "Domain host and port (host:port:sni) sepered by \",\"\n\tEx.: www.google.com.br:443,example.com:443,manage.openshift.com:443")
	traceenable := flag.Bool("trace", false, "Trace to stderr")
	flag.Parse()
	mylog.InitLogs()
	if *traceenable {
		trace.Start(os.Stderr)
		defer trace.Stop()
	}
	domainlist := strings.Split(*domains, ",")
	results := checkcert.CheckHostsParallel(domainlist...)
	for _, result := range results {
		if result.Err == nil {
			fmt.Println(result.ExpireDays)
		}
	}
}
