package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	sub := strings.TrimSuffix(record, "."+mainDomain)
	err := linode.RemoveRecordByName(sub, mainDomain)
	if err != nil {
		log.Fatal(err)
	}

}
