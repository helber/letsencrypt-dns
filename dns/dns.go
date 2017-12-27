package dns

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/bobesa/go-domain-util/domainutil"
)

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

// WaitForPropagation call a list of urls and check DNS propagation
func WaitForPropagation(urls []string, timeout time.Duration, result chan<- bool) {
	var wg sync.WaitGroup
	for _, dom := range urls {
		wg.Add(1)
		go func(url string, timeout time.Duration) {
			start := time.Now()
			end := start.Add(timeout)
			for end.After(time.Now()) {
				response, err := CheckTxt(url)
				if err == nil {
					log.Println("DNS propagation done", url, response)
					wg.Done()
					return
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
