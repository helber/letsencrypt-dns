package main

import (
	"context"
	"log"
	"os"
	"strings"

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
		recordName := "_acme-challenge." + certbotDomain
		filter := cloudflare.DNSRecord{Type: "TXT", Name: recordName, Content: certbotChalenge}
		// records, err := api.GetDNSRecord(context.Background(), &rcontainer, id)
		records, _, err := api.ListDNSRecords(
			context.Background(),
			cloudflare.ZoneIdentifier(id),
			cloudflare.ListDNSRecordsParams{
				Type:    "TXT",
				Name:    recordName,
				Content: certbotChalenge,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		if len(records) == 0 {
			log.Fatalf("record not found filter={Type: \"%v\", Name: \"%v\", Content: \"%v\"}\n", filter.Type, filter.Name, filter.Content)
			os.Exit(1)
		}
		for _, rec := range records {
			log.Println("removing DNS record", rec)
			err := api.DeleteDNSRecord(context.Background(), cloudflare.ZoneIdentifier(certbotDomain), id)
			if err != nil {
				log.Fatal("can't delete record", err)
			}
		}
	}
}
