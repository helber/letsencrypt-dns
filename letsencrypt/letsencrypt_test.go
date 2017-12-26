package letsencrypt

import (
	"testing"
)

const certbotExample string = `
-------------------------------------------------------------------------------
Please read the Terms of Service at
https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf. You must
agree in order to register with the ACME server at
https://acme-v01.api.letsencrypt.org/directory
-------------------------------------------------------------------------------
(A)gree/(C)ancel: A

-------------------------------------------------------------------------------
Would you be willing to share your email address with the Electronic Frontier
Foundation, a founding partner of the Let's Encrypt project and the non-profit
organization that develops Certbot? We'd like to send you email about EFF and
our work to encrypt the web, protect its users and defend digital rights.
-------------------------------------------------------------------------------
(Y)es/(N)o: N
Obtaining a new certificate
Performing the following challenges:
dns-01 challenge for t1.ahgoracloud.com.br
dns-01 challenge for t2.ahgoracloud.com.br

-------------------------------------------------------------------------------
NOTE: The IP of this machine will be publicly logged as having requested this
certificate. If you're running certbot in manual mode on a machine that is not
your server, please ensure you're okay with that.

Are you OK with your IP being logged?
-------------------------------------------------------------------------------
(Y)es/(N)o: Y

-------------------------------------------------------------------------------
Please deploy a DNS TXT record under the name
_acme-challenge.t1.ahgoracloud.com.br with the following value:

7IIJd6NDIz39NrkRPc_QUC35BfuHUXUyaQYZ7ALAH-g

Before continuing, verify the record is deployed.
-------------------------------------------------------------------------------
Press Enter to Continue

-------------------------------------------------------------------------------
Please deploy a DNS TXT record under the name
_acme-challenge.t2.ahgoracloud.com.br with the following value:

oWLBE0JIfo4o83qRrB_4gMLvxYbLwNnAAMFMshJWd7c

Before continuing, verify the record is deployed.
-------------------------------------------------------------------------------
Press Enter to Continue
`

func TestCreateCommandForDomains(t *testing.T) {
	type args struct {
		domains []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"unique domain", args{[]string{"d1.example.com"}}, "certbot certonly --preferred-challenges dns --manual -d d1.example.com"},
		{"3 domain", args{[]string{"d1.example.com", "d2.example.com", "d3.example.com"}}, "certbot certonly --preferred-challenges dns --manual -d d1.example.com -d d2.example.com -d d3.example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateCommandForDomains(tt.args.domains); got != tt.want {
				t.Errorf("CreateCommandForDomains() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
