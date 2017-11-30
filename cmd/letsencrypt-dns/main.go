package main

import (
	"fmt"
	"os"

	"github.com/helber/letsencrypt-dns/dns"
	"github.com/helber/letsencrypt-dns/linode"
)

func main() {
	linode.APIToken = os.Getenv("LINODE_API_KEY")
	// resolv, err := dns.CheckTxt("_acme-challenge.ah-notifications-ahgora.ahgoracloud.com.br")
	// if err != nil {
	// 	fmt.Println("got error", err)
	// } else {
	// 	fmt.Printf("\nresolv response:%s\n\n", resolv)
	// }

	// err := linode.CreateNewTXTRecord("ahgoracloud.com.br", "_lalal_challenge.ahgoracloud.com.br", "My Name is Helber Maciel Guerra")
	// if err != nil {
	// 	fmt.Println("Errro", err)
	// }
	notify := make(chan string, 2)
	notify1 := make(chan string, 2)
	// _acme-challenge.sales-analytics
	go dns.WaitForPublication("_lalal_challenge.ahgoracloud.com.br", notify)
	go dns.WaitForPublication("_acme-challenge.sales-analytics.ahgoracloud.com.br", notify1)
	val := <-notify
	fmt.Println("Got value ", val)
	val1 := <-notify1
	fmt.Println("Got value ", val1)
}
