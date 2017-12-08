package dns

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// CheckTxt execute query on txt
func CheckTxt(urn string) ([]string, error) {
	resp, err := net.LookupTXT(urn)
	if err != nil {
		err := fmt.Errorf("can't lookup %s error %s", urn, err)
		if err == nil {
			fmt.Println("can't print error ", err)
		}
		return resp, err
	}
	return resp, nil
}

// WaitForPropagation call a list of urls and check DNS propagation
func WaitForPropagation(urls []string, timeout time.Duration, result chan<- string) {
	var wg sync.WaitGroup
	for _, dom := range urls {
		wg.Add(1)
		//
		go func(url string, timeout time.Duration) {
			start := time.Now()
			end := start.Add(timeout)
			for end.After(time.Now()) {
				response, err := CheckTxt(url)
				if err != nil {
					fmt.Println("query error", err)
				} else {
					fmt.Println(response)
					wg.Done()
				}
				time.Sleep(time.Second * 5)
			}
			fmt.Println("Timeout...", url)
			wg.Done()
		}(dom, timeout)
	}
	wg.Wait()
	result <- ""

}
