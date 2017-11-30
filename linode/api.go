package linode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIToken to access Linode API
var APIToken string

// Domain object
type Domain struct {
	ID          int      `json:"id"`
	Domain      string   `json:"domain"`
	Domaintype  string   `json:"type"`
	AxfrIps     []string `json:"axfr_ips"`
	Group       string   `json:"group"`
	Status      string   `json:"status"`
	SoaEmail    string   `json:"soa_email"`
	Description string   `json:"description"`
	MasterIps   []string `json:"master_ips"`
	ExpireSec   int      `json:"expire_sec"`
	EetrySec    int      `json:"retry_sec"`
	TTLSec      int      `json:"ttl_sec"`
	RefreshSec  int      `json:"refresh_sec"`
}

// DomainResult Parse linode domain list result
type DomainResult struct {
	Data    []Domain `json:"data"`
	Pages   int      `json:"pages"`
	Page    int      `json:"page"`
	Results int      `json:"results"`
}

// Record on domain
type Record struct {
	ID       int    `json:"id,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Target   string `json:"target,omitempty"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Tag      string `json:"tag,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Port     int    `json:"port,omitempty"`
	Service  string `json:"service,omitempty"`
	TTLSec   int    `json:"ttl_sec,omitempty"`
}

// RecordResult Parse linode record list result
type RecordResult struct {
	Data    []Record `json:"data"`
	Pages   int      `json:"pages"`
	Page    int      `json:"page"`
	Results int      `json:"results"`
}

// GetDomains Get all Domains
func GetDomains() ([]Domain, error) {
	var jsonObjs DomainResult
	cli := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.linode.com/v4/domains/", nil)
	req.Header.Add("Authorization", "Bearer "+APIToken)
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println("Error loading API", err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(body, &jsonObjs)
		if err != nil {
			fmt.Printf("Can't Unmarshal %s", err)
		}
	}
	return jsonObjs.Data, err
}

// AddRecord create new record on linode
func AddRecord(r Record, d Domain) error {
	cli := &http.Client{}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("Can't Marshal %s", err)
	}
	uri := fmt.Sprintf("https://api.linode.com/v4/domains/%d/records/", d.ID)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+APIToken)
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// CreateNewTXTRecord create new record
func CreateNewTXTRecord(domain string, name string, value string) error {
	domains, err := GetDomains()
	if err == nil {
		for _, dom := range domains {
			if dom.Domain == domain {
				rec := Record{Type: "TXT", Name: name, Target: value, TTLSec: 300}
				err := AddRecord(rec, dom)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return err
	}
	return nil
}
