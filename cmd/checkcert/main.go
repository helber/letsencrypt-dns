package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"strings"

	"github.com/helber/letsencrypt-dns/checkcert"
	"github.com/olekukonko/tablewriter"

	mylog "github.com/helber/letsencrypt-dns/log"
	flag "github.com/spf13/pflag"
)

func main() {
	domains := flag.StringP("domains", "d", "", "Domain host and port (host:port) sepered by \",\"\n\tEx.: www.google.com.br:443,example.com:443,manage.openshift.com:443")
	displayTable := flag.BoolP("displaytable", "t", false, "Display host and elapsed query time in a table")
	traceenable := flag.Bool("trace", false, "Trace to stderr")
	flag.Parse()
	mylog.InitLogs()
	if *traceenable {
		trace.Start(os.Stderr)
		defer trace.Stop()
	}
	domainlist := strings.Split(*domains, ",")
	results := checkcert.CheckHostsParallel(domainlist...)
	if *displayTable {
		OutputTable(results)
	} else {
		for _, result := range results {
			if result.Err == nil {
				fmt.Println(result.ExpireDays)
			}
		}
	}
}

// OutputTable set output to ascii table
func OutputTable(results []checkcert.HostResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Query Time", "Expire Days", "Host:port", "Error"})
	for _, res := range results {
		data := []string{}
		data = append(data, fmt.Sprintf("%v", res.ElapsedTime))
		if res.Err == nil {
			data = append(data, fmt.Sprintf("%d", res.ExpireDays))
		} else {
			data = append(data, "")
		}
		data = append(data, res.Host)
		e := res.Err
		if e == nil {
			data = append(data, "")
		} else {
			data = append(data, e.Error())
		}
		table.Append(data)
	}
	table.Render()
}
