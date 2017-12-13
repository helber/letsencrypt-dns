package letsencrypt

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// certbot certonly --preferred-challenges dns --manual -d ah-notifications-ahgora.ahgoracloud.com.br
// https://stackoverflow.com/questions/27322722/interact-with-external-application-from-within-code-golang

// TXTRecord generated to register
type TXTRecord struct {
	domain string
	key    string
}

func parseTopic(topic string) (TXTRecord, error) {
	var pos int
	keyName := ""
	keyValue := ""
	for _, item := range strings.Split(topic, "\n") {
		if pos == 3 {
			keyName = strings.Split(item, " ")[0]
		}
		if pos == 5 {
			keyValue = strings.Split(item, " ")[0]
		}
		pos++
	}
	return TXTRecord{keyName, keyValue}, nil
}

// Call for domains using certbot
func Call(domain string, domains []string, recordChan chan TXTRecord, propagation chan bool, done chan bool) {
	cmd := exec.Command("certbot", "certonly", "--preferred-challenges", "dns", "--manual")
	for _, sub := range domains {
		cmd.Args = append(cmd.Args, "-d")
		cmd.Args = append(cmd.Args, sub)
	}
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if nil != err {
		log.Fatalf("Error obtaining stdin: %s", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		log.Fatalf("Error obtaining stdout: %s", err.Error())
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
				dom, err := parseTopic(string(newBuf))
				if err == nil {
					log.Println(dom)
				}
				recordChan <- TXTRecord{dom.domain, dom.key}
				if parced >= topics {
					// wait for propagation
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
	done <- true
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
