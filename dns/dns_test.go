package dns

import (
	"reflect"
	"testing"
)

func TestRealCheckLookupTxt(t *testing.T) {
	t.Log("resolving _acme-challenge.ah-notifications-ahgora.ahgoracloud.com.br")
	resp, err := CheckTxt("_acme-challenge.ah-notifications-ahgora.ahgoracloud.com.br")
	if err != nil {
		t.Fatal("lookup error ", err)
	}
	t.Log(resp)
}

func TestCheckTxt(t *testing.T) {
	type args struct {
		urn string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"ah-notifications",
			args{"_acme-challenge.ah-notifications-ahgora.ahgoracloud.com.br"},
			[]string{"lxOINBCHdAwZsOFwj4rPI5WBgIEGH9WNIbVoLdFUoRk"},
			false,
		},
		{
			"ah-notifications-err",
			args{"_acme-challengea.ah-notifications-ahgora.ahgoracloud.com.br"},
			[]string{"lxOINBCHdAwZsOFwj4rPI5WBgIEGH9WNIbVoLdFUoRk"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckTxt(tt.args.urn)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckTxt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckTxt() = %v, want %v", got, tt.want)
			}
		})
	}
}
