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

const header0 string = `(A)gree/(C)ancel:`
const header1 string = `digital rights.
-------------------------------------------------------------------------------
(Y)es/(N)o:`

const header2 string = `Are you OK with your IP being logged?
-------------------------------------------------------------------------------
(Y)es/(N)o:`

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

// CallAuto
func CallAuto(domains []string, done chan bool) error {
	// certbot certonly --manual-public-ip-logging-ok --agree-tos -n -m sre@ahgora.com.br --preferred-challenges=dns --test-cert  --manual --manual-auth-hook /opt/certbot/validation.sh --manual-cleanup-hook /opt/certbot/clean.sh -d t1.ahgoracloud.com.br -d t2.ahgoracloud.com.br
	cmd := exec.Command(
		"certbot",
		"certonly",
		"--agree-tos",
		"--manual-public-ip-logging-ok",
		"-m",
		"sre@ahgora.com.br",
		"-n",
		"--preferred-challenges=dns",
		"--manual",
		"--manual-auth-hook",
		"letsencrypt-validate",
		"--manual-cleanup-hook",
		"letsencrypt-cleanup",
	)
	for _, sub := range domains {
		cmd.Args = append(cmd.Args, "-d")
		cmd.Args = append(cmd.Args, sub)
	}
	log.Println(cmd.Args)
	// Read output
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if nil != err {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		return err
	}
	reader := bufio.NewReader(stdout)
	// Parse stdout
	go func(reader io.Reader) {
		defer stdin.Close()
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			txt := scanner.Text()
			log.Println("-->", txt)
		}
	}(reader)

	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}
	log.Println("wait for command done")
	// Wait
	err = cmd.Wait()
	log.Println("command done")
	if err != nil {
		log.Printf("error=%v", err)
		done <- false
		return err
	}
	log.Println("send true to channel")
	done <- true
	log.Println("channel done")
	return nil
}

// Call for domains using certbot
func Call(domain string, domains []string, done chan bool) error {
	var generatedDomains []string
	var recordList []linode.Record
	domainObj, err := linode.GetDomainObject(domain)
	if err != nil {
		return err
	}
	cmd := exec.Command("certbot", "certonly", "-m", "sre@ahgora.com.br", "--preferred-challenges", "dns", "--manual")
	for _, sub := range domains {
		cmd.Args = append(cmd.Args, "-d")
		cmd.Args = append(cmd.Args, sub)
	}
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if nil != err {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		return err
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
			// Header0
			if bytes.Contains(newBuf, []byte(header0)) {
				newBuf = []byte{}
				io.WriteString(stdin, "A\n")
			}
			// Header1
			if bytes.Contains(newBuf, []byte(header1)) {
				newBuf = []byte{}
				io.WriteString(stdin, "N\n")
			}
			// Header2
			if bytes.Contains(newBuf, []byte(header2)) {
				newBuf = []byte{}
				io.WriteString(stdin, "Y\n")
			}
			// Topic
			if bytes.Contains(newBuf, []byte("Press Enter to Continue")) {
				parced++
				toParse := string(newBuf)
				dom, err := parseTopic(toParse)
				if err != nil {
					log.Fatalf("Error parsing buffer:%s\n%s", err, toParse)
				}
				generatedDomains = append(generatedDomains, dom.domain)
				name := strings.TrimSuffix(dom.domain, domain)
				// Register new TXT record
				rec := linode.Record{Type: "TXT", Name: name, Target: dom.key, TTLSec: 300}
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
					log.Println("all dns propagation done")
					io.WriteString(stdin, "\n")
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
	log.Println("wait for command done")
	cmd.Wait()
	log.Println("command done")
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
	log.Printf("CMD=%s", cmd)
	for _, dom := range domains {
		cmd += "-d "
		cmd += dom + " "
	}
	return strings.TrimSpace(cmd)
}
