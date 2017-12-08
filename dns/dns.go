package dns

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type urnCheck struct {
	urn  string
	done bool
}

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

// WaitForPublication call a list of urns and
func WaitForPublication(urns []string, timeout time.Duration, result chan<- string) {
	var wg sync.WaitGroup
	for _, dom := range urns {
		wg.Add(1)
		go func(url urnCheck, timeout time.Duration) {
			start := time.Now()
			end := start.Add(timeout)
			for end.After(time.Now()) {
				if url.done == false {
					response, err := CheckTxt(url.urn)
					if err != nil {
						fmt.Println("query error", err)
					} else {
						fmt.Println(response)
						url.done = true
						wg.Done()
					}
				} else {
					fmt.Println("Checked", url.urn)
				}
				time.Sleep(time.Second * 5)
			}
			fmt.Println("Timeout...", url.urn)
			wg.Done()
		}(urnCheck{dom, false}, timeout)
	}
	wg.Wait()
	result <- ""

}
