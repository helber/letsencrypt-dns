package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/helber/letsencrypt-dns/dns"
	"github.com/helber/letsencrypt-dns/linode"
	mylog "github.com/helber/letsencrypt-dns/log"
)

func main() {
	linode.APIToken = os.Getenv("LINODE_API_KEY")
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
	err := linode.RemoveRecordByName(sub, mainDomain)
	if err != nil {
		log.Fatal(err)
	}

}
