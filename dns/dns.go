package dns

import (
	"fmt"
	"net"
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

// WaitForPublication
func WaitForPublication(urn string, result chan<- string) {
	for i := 0; i <= 10; i++ {
		response, err := CheckTxt(urn)
		if err != nil {
			fmt.Println("query error", err)
		} else {
			fmt.Println(response)
			if response[0] != "" {
				result <- response[0]
			}
		}
		time.Sleep(time.Second * 1)
	}
	result <- ""

}
