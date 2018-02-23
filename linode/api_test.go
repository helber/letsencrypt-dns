package linode

import (
	"encoding/json"
	"testing"
)

const result = `
{
    "data": [
        {
            "axfr_ips": [],
            "ttl_sec": 300,
            "id": 40831,
            "refresh_sec": 300,
            "type": "master",
            "group": "",
            "status": "active",
            "domain": "ahgorasistemas.com.br",
            "soa_email": "postmaster@ahgora.com.br",
            "description": "",
            "expire_sec": 0,
            "master_ips": [],
            "retry_sec": 300
        },
        {
            "axfr_ips": [],
            "ttl_sec": 300,
            "id": 907590,
            "refresh_sec": 300,
            "type": "master",
            "group": "",
            "status": "active",
            "domain": "ahgoracloud.com.br",
            "soa_email": "postmaster@ahgora.com.br",
            "description": "",
            "expire_sec": 604800,
            "master_ips": [],
            "retry_sec": 300
        }
    ],
    "pages": 1,
    "page": 1,
    "results": 2
}`

const records = `{
    "data": [
    {
        "priority": 10,
        "tag": null,
        "protocol": null,
        "id": 8521741,
        "weight": 5,
        "type": "CNAME",
        "target": "console.ahgoracloud.com.br",
        "port": 80,
        "service": null,
        "name": "ah-notifications-ahgora",
        "ttl_sec": 0
    },
    {
        "priority": 0,
        "tag": null,
        "protocol": null,
        "id": 8541091,
        "weight": 0,
        "type": "TXT",
        "target": "XXXXXXX",
        "port": 0,
        "service": null,
        "name": "_test_2",
        "ttl_sec": 0
    },
    {
        "priority": 10,
        "tag": null,
        "protocol": null,
        "id": 8541092,
        "weight": 5,
        "type": "TXT",
        "target": "lxOINBCHdAwZsOFwj4rPI5WBgIEGH9WNIbVoLdFUoRk",
        "port": 80,
        "service": null,
        "name": "_acme-challenge.ah-notifications-ahgora",
        "ttl_sec": 300
    }
    ],
    "page": 1,
    "pages": 1,
    "results": 3
}`

func TestUnmarshal(t *testing.T) {
	rec := Record{Type: "TXT", Name: "_acme-challenge.xpto.ahgoracloud.com.br", Target: "xpto"}
	out, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("can't Marshal %v(%s)", rec, err)
	}
	t.Logf("OUT=%s", out)

}
