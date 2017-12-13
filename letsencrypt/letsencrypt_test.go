package letsencrypt

import (
	"io/ioutil"
	"testing"
)

const certbotExample string = `
certbot certonly --preferred-challenges dns --manual -d ahgoracloud.com.br -d console.ahgoracloud.com.br -d metrics.ahgoracloud.com.br -d logs.ahgoracloud.com.br -d node.ahgoracloud.com.br
Saving debug log to /var/log/letsencrypt/letsencrypt.log
Plugins selected: Authenticator manual, Installer None
Obtaining a new certificate
Performing the following challenges:
dns-01 challenge for ahgoracloud.com.br
dns-01 challenge for console.ahgoracloud.com.br
dns-01 challenge for metrics.ahgoracloud.com.br
dns-01 challenge for logs.ahgoracloud.com.br
dns-01 challenge for node.ahgoracloud.com.br

-------------------------------------------------------------------------------
NOTE: The IP of this machine will be publicly logged as having requested this
certificate. If you're running certbot in manual mode on a machine that is not
your server, please ensure you're okay with that.

Are you OK with your IP being logged?
-------------------------------------------------------------------------------
(Y)es/(N)o: Y

-------------------------------------------------------------------------------
Please deploy a DNS TXT record under the name
_acme-challenge.ahgoracloud.com.br with the following value:

OP3xHdUaRtLB6ou3UQ878UfuZ0zo_i0wAW7sWdSNqSw

Before continuing, verify the record is deployed.
-------------------------------------------------------------------------------
Press Enter to Continue

-------------------------------------------------------------------------------
Please deploy a DNS TXT record under the name
_acme-challenge.console.ahgoracloud.com.br with the following value:

k-wVRRiv5aBt7RhSFsmq2N6Qa-qU5Ykj5LoIsrk2rM4

Before continuing, verify the record is deployed.
-------------------------------------------------------------------------------
Press Enter to Continue

-------------------------------------------------------------------------------
Please deploy a DNS TXT record under the name
_acme-challenge.metrics.ahgoracloud.com.br with the following value:

IHsh4SiyYrsQF0mh__Lg2hGA0rlot02VguROLn07plM

Before continuing, verify the record is deployed.
-------------------------------------------------------------------------------
Press Enter to Continue

-------------------------------------------------------------------------------
Please deploy a DNS TXT record under the name
_acme-challenge.logs.ahgoracloud.com.br with the following value:

Z-Vn6VfLAHIlZBqmP92SOjY6afB1Lu5yM5110qwEeAE

Before continuing, verify the record is deployed.
-------------------------------------------------------------------------------
Press Enter to Continue

-------------------------------------------------------------------------------
Please deploy a DNS TXT record under the name
_acme-challenge.node.ahgoracloud.com.br with the following value:

VurKbQclijCZwGJl9sXDlrx4LABr2gU2BL_OiPmaUzw

Before continuing, verify the record is deployed.
-------------------------------------------------------------------------------
Press Enter to Continue
Waiting for verification...
Cleaning up challenges

IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live/ahgoracloud.com.br/fullchain.pem
   Your key file has been saved at:
   /etc/letsencrypt/live/ahgoracloud.com.br/privkey.pem
   Your cert will expire on 2018-03-11. To obtain a new or tweaked
   version of this certificate in the future, simply run certbot
   again. To non-interactively renew *all* of your certificates, run
   "certbot renew"
 - If you like Certbot, please consider supporting our work by:

   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
   Donating to EFF:                    https://eff.org/donate-le

`

func TestRmFiles(t *testing.T) {
	type args struct {
		domains []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"f1 f2", args{[]string{"f1", "f2"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, f := range tt.args.domains {
				t.Log(f)
				// golden := filepath.Join("testdata", f)
				ioutil.WriteFile(f, []byte("---"), 0644)
			}
			if got := Call(tt.args.domains); got != tt.want {
				t.Errorf("Call() = %v, want %v", got, tt.want)
			}
		})
	}

}

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
