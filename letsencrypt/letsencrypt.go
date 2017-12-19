package letsencrypt

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/helber/letsencrypt-dns/dns"
	"github.com/helber/letsencrypt-dns/linode"
)

// certbot certonly --preferred-challenges dns --manual -d ah-notifications-ahgora.ahgoracloud.com.br
// https://stackoverflow.com/questions/27322722/interact-with-external-application-from-within-code-golang

// TXTRecord generated to register
type TXTRecord struct {
	domain string
	key    string
}

func parseTopic(topic string) (TXTRecord, error) {
	keyName := ""
	keyValue := ""
	for i, item := range strings.Split(topic, "\n") {
		if i == 3 {
			keyName = strings.Split(item, " ")[0]
		}
		if i == 5 {
			keyValue = strings.Split(item, " ")[0]
		}
	}
	return TXTRecord{keyName, keyValue}, nil
}

// Call for domains using certbot
func Call(domain string, domains []string, done chan bool) error {
	var generatedDomains []string
	var recordList []linode.Record
	domainObj, err := linode.GetDomainObject(domain)
	if err != nil {
		return err
	}
	cmd := exec.Command("certbot", "certonly", "--preferred-challenges", "dns", "--manual")
	for _, sub := range domains {
		cmd.Args = append(cmd.Args, "-d")
		cmd.Args = append(cmd.Args, sub)
	}
	log.Printf("calling %s", cmd.Args)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if nil != err {
		return err
		// log.Fatalf("Error obtaining stdin: %s", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		return err
		// log.Fatalf("Error obtaining stdout: %s", err.Error())
	}
	reader := bufio.NewReader(stdout)
	// Parse stdout
	go func(reader io.Reader, topics int) {
		defer stdin.Close()
		parced := 0
		var newBuf []byte
		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {
			txt := scanner.Bytes()
			newBuf = append(newBuf, txt[0])
			// Header
			if bytes.Contains(newBuf, []byte("(Y)es/(N)o:")) {
				newBuf = []byte{}
				io.WriteString(stdin, "Y\n")
			}
			// Topic
			if bytes.Contains(newBuf, []byte("Press Enter to Continue")) {
				parced++
				toParse := string(newBuf)
				dom, err := parseTopic(toParse)
				if err == nil {
					log.Fatalf("Error parsing buffer:%s\n%s", err.Error(), toParse)
				}
				generatedDomains = append(generatedDomains, dom.domain)
				// Register new TXT record
				rec := linode.Record{Type: "TXT", Name: dom.domain, Target: dom.key, TTLSec: 300}
				record, err := linode.AddRecord(rec, domainObj)
				if err != nil {
					log.Fatalf("can't create new record:%s", err.Error())
				}
				recordList = append(recordList, record)
				if parced >= topics {
					propagation := make(chan bool)
					dns.WaitForPropagation(generatedDomains, 10*time.Minute, propagation)
					// wait for propagation before press ENTER
					<-propagation
				}
				newBuf = []byte{}
				io.WriteString(stdin, "\n")
			}
		}
	}(reader, len(domains))
	// Wait
	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}
	cmd.Wait()
	// Clean registered domains
	for _, record := range recordList {
		err := linode.RemoveRecord(record, domainObj)
		if err != nil {
			log.Println("error removing", record)
		}
	}
	done <- true
	return nil
}

// CreateCommandForDomains create a certbot command call for a list of domains
func CreateCommandForDomains(domains []string) string {
	cmd := "certbot certonly --preferred-challenges dns --manual "
	for _, dom := range domains {
		cmd += "-d "
		cmd += dom + " "
	}
	return strings.TrimSpace(cmd)
}
