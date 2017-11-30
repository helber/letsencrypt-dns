package dns

import (
	"testing"
)

func realCheckLookupTxt(t *testing.T) {
	t.Log("resolving _acme-challenge.ah-notifications-ahgora.ahgoracloud.com.br")
	resp, err := CheckTxt("_acme-challenge.ah-notifications-ahgora.ahgoracloud.com.br")
	if err != nil {
		t.Fatal("lookup error ", err)
	}
	t.Log(resp)
}
