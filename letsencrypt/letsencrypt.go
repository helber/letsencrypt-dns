package letsencrypt

import (
	"fmt"
	"os/exec"
	"strings"
)

// certbot certonly --preferred-challenges dns --manual -d ah-notifications-ahgora.ahgoracloud.com.br
// https://stackoverflow.com/questions/27322722/interact-with-external-application-from-within-code-golang

// Call for domains using certbot
func Call(domains []string) {
	// cmd := exec.Command("test.py", "certonly", "--preferred-challenges", "dns", "--manual")
	cmd := exec.Command("rm", "-i", "f1", "f2", "f3")
	// for _, domain := range domains {
	// 	cmd.Args = append(cmd.Args, "-d")
	// 	cmd.Args = append(cmd.Args, domain)
	// }
	fmt.Println(cmd.Args)
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
