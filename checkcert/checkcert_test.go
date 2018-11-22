package checkcert

import (
	"fmt"
	"testing"
)

func TestParseHostPortDomain(t *testing.T) {
	t.Parallel()
	tt := []struct {
		in     string
		host   string
		port   int
		domain string
	}{
		{in: "example.com", host: "example.com", port: 443, domain: "example.com"},
		{in: "www.example.com:443:h1.example.com", host: "h1.example.com", port: 443, domain: "www.example.com"},
		{in: "www.example.com:8443", host: "www.example.com", port: 8443, domain: "www.example.com"},
		{in: "mail.example.com::host1.example.com", host: "host1.example.com", port: 443, domain: "mail.example.com"},
		{in: "www.low.com:8080:192.168.55.1", host: "192.168.55.1", port: 8080, domain: "www.low.com"},
	}
	for i, x := range tt {
		t.Run(fmt.Sprintf("sub test (%d) -> %v", i, x), func(st *testing.T) {
			host, port, domain := ParseHostPortDomain(x.in)
			if host != x.host {
				st.Errorf("host expected %s, got %s", x.host, host)
			}
			if port != x.port {
				st.Errorf("port expected %d, got %d", x.port, port)
			}
			if domain != x.domain {
				st.Errorf("domain expected %s, got %s", x.domain, domain)
			}
		})
	}
}
