package dns

import (
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/bobesa/go-domain-util/domainutil"
)

// GetMainDomain return main domain of FQDN
func GetMainDomain(hostname string) string {
	return domainutil.Domain(hostname)
}

// CheckTxt execute query on txt
func CheckTxt(urn string) ([]string, error) {
	resp, err := net.LookupTXT(urn)
	if err != nil {
		log.Printf("wait DNS txt %s", urn)
		return resp, err
	}
	return resp, nil
}

// GetNsProvider query ns servers and return provider
func GetNsProvider(domain string) (prov string, err error) {
	ns, err := net.LookupNS(domain)
	if err != nil {
		log.Println("lookup error", err)
		return
	}
	for i := len(ns) - 1; i >= 0; i-- {
		prov = domainutil.Domain(strings.TrimSuffix(ns[i].Host, "."))
		log.Println(ns[i].Host, "-->", prov)
	}
	return
}

// WaitForPropagation call a list of urls and check DNS propagation
func WaitForPropagation(urls []string, timeout time.Duration, result chan<- bool) {
	var wg sync.WaitGroup
	sucess := 1
	for _, dom := range urls {
		wg.Add(1)
		go func(url string, timeout time.Duration) {
			start := time.Now()
			end := start.Add(timeout)
			for end.After(time.Now()) {
				response, err := CheckTxt(url)
				if err == nil {
					log.Printf("DNS OK %s %v %d\n", url, response, sucess)
					if sucess > 2 {
						log.Println("DNS propagation done", url, response)
						wg.Done()
						return
					}
					sucess++
				} else {
					sucess = 0
				}
				time.Sleep(time.Second * 30)
			}
			log.Println("Timeout...", url)
			wg.Done()
			result <- false
		}(dom, timeout)
	}
	wg.Wait()
	result <- true
}
