package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/helber/letsencrypt-dns/dns"
	"github.com/helber/letsencrypt-dns/linode"
)

func main() {
	linode.APIToken = os.Getenv("LINODE_API_KEY")
	certbotDomain := os.Getenv("CERTBOT_DOMAIN")
	certbotChalenge := os.Getenv("CERTBOT_VALIDATION")

	if certbotDomain == "" {
		log.Fatal("domain env (CERTBOT_DOMAIN) var not found")
	}
	if certbotChalenge == "" {
		log.Fatal("validator env (CERTBOT_VALIDATION) var not found")
	}
	mainDomain := dns.GetMainDomain(certbotDomain)
	record := fmt.Sprintf("_acme-challenge.%s", certbotDomain)

	recObj, err := linode.CreateNewTXTRecord(mainDomain, record, certbotChalenge)
	if err != nil {
		log.Fatalln("erro", err)
	}
	log.Printf("Record created ID=%d Obj=%v", recObj.ID, recObj)

	notify := make(chan bool)
	defer close(notify)
	go dns.WaitForPropagation([]string{record}, time.Minute*60, notify)
	fmt.Println("Wait for publication")
	val := <-notify
	fmt.Println("Got value ", val)
}
