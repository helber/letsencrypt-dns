package main

import (
	"flag"
	"log"

	"github.com/helber/letsencrypt-dns/checkcert"
	mylog "github.com/helber/letsencrypt-dns/log"
)

func main() {
	flag.Parse()
	mylog.InitLogs()
	// res, err := checkcert.CheckHost("pw2-socket.ahgoracloud.com.br:443")
	// if err != nil {
	// 	log.Printf("error %v", err)
	// }
	// log.Println(res)
	// res, err = checkcert.CheckHost("ahgoracloud.com.br:443")
	// if err != nil {
	// 	log.Printf("error %v", err)
	// }
	// log.Println(res)
	// res, err = checkcert.CheckHost("console.ahgoracloud.com.br:443")
	// if err != nil {
	// 	log.Printf("error %v", err)
	// }
	// log.Println(res)
	res, err := checkcert.CheckHost("www.ahgora.com.br:443")
	if err != nil {
		log.Printf("error %v", err)
	}
	log.Println(res)
}
