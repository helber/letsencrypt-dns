package main

import (
	"flag"
	"fmt"
	"log"
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
		log.Panic("at last 1 domain is required ", len(domainlist)-1)
	}
	if *mainDomain == "" {
		log.Panic("main domain required")
	}
	log.Println("Generating cert on (", *mainDomain, "):", domainlist, "domains")
	// Done Channel
	done := make(chan bool)
	err := letsencrypt.Call(*mainDomain, domainlist, done)
	if err != nil {
		log.Panicf("can't call letsencrypt: %s", err)
	}
	result := <-done
	if result == true {
		fmt.Print("Congratulations")
	}

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
	// sldk := []DomainKv{
	// 	{"_acme-challenge.ahgoracloud.com.br", "OP3xHdUaRtLB6ou3UQ878UfuZ0zo_i0wAW7sWdSNqSw"},
	// 	{"_acme-challenge.console.ahgoracloud.com.br", "k-wVRRiv5aBt7RhSFsmq2N6Qa-qU5Ykj5LoIsrk2rM4"},
	// 	{"_acme-challenge.metrics.ahgoracloud.com.br", "IHsh4SiyYrsQF0mh__Lg2hGA0rlot02VguROLn07plM"},
	// 	{"_acme-challenge.logs.ahgoracloud.com.br", "Z-Vn6VfLAHIlZBqmP92SOjY6afB1Lu5yM5110qwEeAE"},
	// 	{"_acme-challenge.node.ahgoracloud.com.br", "VurKbQclijCZwGJl9sXDlrx4LABr2gU2BL_OiPmaUzw"},
	// }

	// preDomains := []string{
	// 	"_acme-challenge.ahgoracloud.com.br",
	// 	"_acme-challenge.console.ahgoracloud.com.br",
	// 	"_acme-challenge.metrics.ahgoracloud.com.br",
	// 	"_acme-challenge.logs.ahgoracloud.com.br",
	// 	"_acme-challenge.node.ahgoracloud.com.br",
	// }

	// notify := make(chan bool)
	// defer close(notify)
	// go dns.WaitForPropagation(preDomains, time.Minute*10, notify)
	// fmt.Println("Wait for publication")
	// val := <-notify
	// fmt.Println("Got value ", val)

}
