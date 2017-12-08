package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/helber/letsencrypt-dns/dns"
	"github.com/helber/letsencrypt-dns/letsencrypt"
	"github.com/helber/letsencrypt-dns/linode"
)

func main() {
	linode.APIToken = os.Getenv("LINODE_API_KEY")
	domains := flag.String("d", "", "Domains sepered by \",\"")
	flag.Parse()
	domainlist := strings.Split(*domains, ",")
	fmt.Println("Generating cert for:", domainlist, "domains")
	letsencrypt.Call(domainlist)

	// resolv, err := dns.CheckTxt("_acme-challenge.ah-notifications-ahgora.ahgoracloud.com.br")
	// if err != nil {
	// 	fmt.Println("got error", err)
	// } else {
	// 	fmt.Printf("\nresolv response:%s\n\n", resolv)
	// }

	// err := linode.CreateNewTXTRecord("ahgoracloud.com.br", "_lalal_challenge.ahgoracloud.com.br", "My Name is Helber Maciel Guerra")
	// if err != nil {
	// 	fmt.Println("Errro", err)
	// }

	notify := make(chan string, 2)
	defer close(notify)
	// // _acme-challenge.sales-analytics
	testDom := []string{"A_lalal_challenge.ahgoracloud.com.br", "_acme-challenge.sales-analytics.ahgoracloud.com.br", "_acme-challenge-fall.sales-analytics.ahgoracloud.com.br"}
	go dns.WaitForPropagation(testDom, time.Second*30, notify)
	fmt.Println("Wait for publication")
	val := <-notify
	fmt.Println("Got value ", val)
}
