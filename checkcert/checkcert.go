package checkcert

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// HostResult Results
type HostResult struct {
	Host        string
	ExpireDays  int64
	Err         error
	Issuer      string
	TLSVersion  string
	ElapsedTime time.Duration
}

// CheckHost check cert
func CheckHost(host string, port int, domain string, res chan<- HostResult) {
	full := fmt.Sprintf("%s:%d:%s", host, port, domain)
	log.Printf("started > %s", full)
	start := time.Now()
	result := HostResult{
		Host:       full,
		ExpireDays: -1,
	}
	ip, err := net.LookupHost(host)
	if err != nil {
		log.Printf("Could not resolve domain name, %v.\n\n", host)
		result.Err = err
		res <- result
		return
	}
	ipAddress := fmt.Sprintf("%s:%d", ip[0], port)
	//Connect network
	ipConn, err := net.DialTimeout("tcp", ipAddress, time.Second*2)
	if err != nil {
		log.Printf("Could not connect to %v - %v\n", ipAddress, host)
		result.Err = err
		res <- result
		return
	}
	log.Printf("Connected to %v - %v\n", ipAddress, host)
	defer ipConn.Close()
	// Configure tls to look at domain
	config := tls.Config{ServerName: domain}
	// Connect to tls
	conn := tls.Client(ipConn, &config)
	defer conn.Close()
	// Handshake with TLS to get cert
	hsErr := conn.Handshake()
	if hsErr != nil {
		log.Printf("Client connected to: %v\n", conn.RemoteAddr())
		log.Printf("Cert Failed for %v - %v\n", ipAddress, domain)
		result.Err = err
		res <- result
		return
	}
	timeNow := time.Now()
	state := conn.ConnectionState()
	for i, v := range state.PeerCertificates {
		switch i {
		case 0:
			log.Println("Server key information:")
			switch v.Version {
			case 3:
				log.Printf("\tVersion: TLS v1.2\n")
				result.TLSVersion = "TLS v1.2"
			case 2:
				log.Printf("\tVersion: TLS v1.1\n")
				result.TLSVersion = "TLS v1.1"
			case 1:
				log.Printf("\tVersion: TLS v1.0\n")
				result.TLSVersion = "TLS v1.0"
			case 0:
				log.Printf("\tVersion: SSL v3\n")
				result.TLSVersion = "SSL v3"
			}
			log.Printf("\tCN:\t %v\n\tOU:\t %v\n\tOrg:\t %v\n", v.Subject.CommonName, v.Subject.OrganizationalUnit, v.Subject.Organization)
			log.Printf("\tCity:\t %v\n\tState:\t %v\n\tCountry: %v\n", v.Subject.Locality, v.Subject.Province, v.Subject.Country)
			log.Printf("SSL Certificate Valid:\n\tFrom:\t %v\n\tTo:\t %v\n", v.NotBefore, v.NotAfter)
			na := v.NotAfter.Sub(timeNow).Hours()
			expiresIn := int64(na)
			result.ExpireDays = expiresIn / 24
			log.Printf("Valid Certificate DNS:\n")
			if len(v.DNSNames) >= 1 {
				for dns := range v.DNSNames {
					log.Printf("\t%v\n", v.DNSNames[dns])
				}
			} else {
				log.Printf("\t%v\n", v.Subject.CommonName)
			}
		case 1:
			log.Printf("Issued by:\n\t%v\n\t%v\n\t%v\n", v.Subject.CommonName, v.Subject.OrganizationalUnit, v.Subject.Organization)
			result.Issuer = v.Subject.Organization[0]
		default:
			break
		}
	}
	t := time.Now()
	result.ElapsedTime = t.Sub(start)
	log.Printf("finished %v in %v", result.Host, result.ElapsedTime)
	res <- result
}

// ParseHostPortDomain parse host:port:domain
func ParseHostPortDomain(info string) (host string, port int, domain string) {
	splt := strings.Split(info, ":")
	port = 443
	if len(splt) > 1 {
		portn, err := strconv.Atoi(splt[1])
		if err == nil {
			port = portn
		}
	}
	host = splt[0]
	if len(splt) > 2 {
		domain = splt[2]
	} else {
		domain = host
	}
	return
}

// CheckHostsParallel Return a slice of host results given some hosts
//
func CheckHostsParallel(hosts ...string) (res []HostResult) {
	results := make(chan HostResult, len(hosts))
	for _, dom := range hosts {
		wg.Add(1)
		host, port, domain := ParseHostPortDomain(dom)
		go CheckHost(host, port, domain, results)
	}
	// fmt.Fprintf(os.Stderr, "Created [%d] checks\n", len(hosts))
	for range hosts {
		resT := <-results
		log.Println(resT)
		res = append(res, resT)
		// fmt.Fprintf(os.Stderr, "HOST [%s] done in %s\n", resT.Host, resT.ElapsedTime)
		wg.Done()
	}
	wg.Wait()
	return
}
