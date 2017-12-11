package letsencrypt

import "testing"

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
