package letsencrypt

import (
	"bufio"
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

// CallAuto Call letsencrypt using automation
func CallAuto(domains []string, done chan bool) error {
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
