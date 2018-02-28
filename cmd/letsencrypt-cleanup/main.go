package main

import (
	"flag"
	"log"
	"os"
	"strings"

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
	log.Printf("domain=[%s] chalenge=[%s]\n", certbotDomain, certbotChalenge)

	if certbotDomain == "" {
		log.Fatal("domain env (CERTBOT_DOMAIN) var not found")
	}
	if certbotChalenge == "" {
		log.Fatal("validator env (CERTBOT_VALIDATION) var not found")
	}
	mainDomain := dns.GetMainDomain(certbotDomain)
	log.Printf("main domain=%s\n", mainDomain)
	sub := strings.TrimSuffix(certbotDomain, "."+mainDomain)
	provider, err := dns.GetNsProvider(mainDomain)
	if err != nil {
		log.Fatal("can't get provider", err)
	}
	if provider == "linode.com" {
		linode.APIToken = os.Getenv("LINODE_API_KEY")
		err := linode.RemoveRecordByName(sub, mainDomain)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
		if err != nil {
			log.Fatal(err)
		}
		id, err := api.ZoneIDByName(mainDomain)
		if err != nil {
			log.Fatal(err)
		}
		// Fetch records
		records, err := api.DNSRecords(id, cloudflare.DNSRecord{Type: "TXT", Name: certbotDomain})
		if err != nil {
			log.Fatal(err)
		}
		for _, rec := range records {
			log.Println("REMOVING DNS RECORD", rec)
			err := api.DeleteDNSRecord(id, rec.ID)
			if err != nil {
				log.Fatal("can't delete record", err)
			}
		}
	}
}
