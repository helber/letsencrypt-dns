package checkcert

import (
	"crypto/tls"
	"log"
	"net"
	"strings"
	"time"
)

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
func CheckHost(hostPort string, res chan<- HostResult) {
	log.Printf("started > %s", hostPort)
	start := time.Now()
	result := HostResult{
		Host:       hostPort,
		ExpireDays: -1,
	}

	splt := strings.Split(hostPort, ":")
	domainName, port := splt[0], splt[1]
	ip, err := net.LookupHost(domainName)
	if err != nil {
		log.Printf("Could not resolve domain name, %v.\n\n", domainName)
		log.Printf("Either supply a valid domain name or use the -i switch to supply the ip address.\n")
		log.Printf("Domain name lookups are not performed when the user provides the ip address.\n")
		res <- result
	}
	ipAddress := ip[0] + ":" + port
	//Connect network
	ipConn, err := net.DialTimeout("tcp", ipAddress, 5*time.Second)
	if err != nil {
		log.Printf("Could not connect to %v - %v\n", ipAddress, domainName)
		res <- result
	}
	defer ipConn.Close()
	// Configure tls to look at domainName
	config := tls.Config{ServerName: domainName}
	// Connect to tls
	conn := tls.Client(ipConn, &config)
	defer conn.Close()
	// Handshake with TLS to get cert
	hsErr := conn.Handshake()
	if hsErr != nil {
		log.Printf("Client connected to: %v\n", conn.RemoteAddr())
		log.Printf("Cert Failed for %v - %v\n", ipAddress, domainName)
		res <- result
	}
	log.Printf("Client connected to: %v\n", conn.RemoteAddr())
	log.Printf("Cert Checks OK\n")
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
