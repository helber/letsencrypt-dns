package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/helber/letsencrypt-dns/linode"
)

func main() {
	linode.APIToken = os.Getenv("LINODE_API_KEY")
	domains := flag.String("d", "", "Domains sepered by \",\"")
	mainDomain := flag.String("domain", "", "Main domain")
	flag.Parse()
	domainlist := strings.Split(*domains, ",")
	log.Println("Generating cert on (", *mainDomain, "):", domainlist, "domains")
	// Comunication Channels
	// propagation := make(chan bool)
	// done := make(chan bool)
	// txtRecords := make(chan letsencrypt.TXTRecord)

	// go dns.WaitForPropagation(txtDomain, 10*time.Minute, propagation)
	// letsencrypt.Call(*mainDomain, domainlist, txtRecords, propagation, done)

	domainObj, err := linode.GetDomainObject(*mainDomain)
	if err != nil {
		log.Panic(err)
	}

	records := []linode.Record{}
	for i, kv := range domainlist {
		value := fmt.Sprintf("__%v_challenge.%v", i, kv)
		rec, err := linode.CreateNewTXTRecord(*mainDomain, kv, value)
		records = append(records, rec)
		if err != nil {
			log.Fatal("can't create record")
			os.Exit(1)
		}
		log.Println("New record created", rec)
	}
	time.Sleep(time.Minute * 3)
	// Clean
	for _, record := range records {
		err := linode.RemoveRecord(record, domainObj)
		if err != nil {
			log.Println("error removing", record)
		}
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
