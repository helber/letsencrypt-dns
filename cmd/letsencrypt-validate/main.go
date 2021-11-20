package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/helber/letsencrypt-dns/dns"
	"github.com/helber/letsencrypt-dns/linode"
	mylog "github.com/helber/letsencrypt-dns/log"
)

func main() {
	certbotDomain := os.Getenv("CERTBOT_DOMAIN")
	certbotChalenge := os.Getenv("CERTBOT_VALIDATION")
	flag.Parse()
	mylog.InitLogs()

	if certbotDomain == "" {
		log.Fatal("domain env (CERTBOT_DOMAIN) var not found")
	}
	if certbotChalenge == "" {
		log.Fatal("validator env (CERTBOT_VALIDATION) var not found")
	}
	mainDomain := dns.GetMainDomain(certbotDomain)
	record := fmt.Sprintf("_acme-challenge.%s", certbotDomain)
	provider, err := dns.GetNsProvider(mainDomain)
	if err != nil {
		log.Fatal("can't get provider", err)
	}

	if provider == "linode.com" {
		linode.APIToken = os.Getenv("LINODE_API_KEY")
		recObj, err := linode.CreateNewTXTRecord(mainDomain, record, certbotChalenge)
		if err != nil {
			log.Fatalln("erro", err)
		}
		log.Printf("Record created ID=%d Obj=%v", recObj.ID, recObj)
	} else {
		api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
		if err != nil {
			log.Fatal(err)
		}
		id, err := api.ZoneIDByName(mainDomain)
		if err != nil {
			log.Fatal(err)
		}
		record, err := api.CreateDNSRecord(context.Background(), id, cloudflare.DNSRecord{Type: "TXT", Name: record, Content: certbotChalenge, TTL: 300})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("New record", record)
	}
	notify := make(chan bool)
	defer close(notify)
	go dns.WaitForPropagation([]string{record}, time.Minute*60, notify)
	log.Println("Wait for publication")
	val := <-notify
	log.Println("Got value ", val)
}
