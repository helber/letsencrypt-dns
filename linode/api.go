package linode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
		return []Domain{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return []Domain{}, errors.New("linode api server error")
	}
	if resp.StatusCode >= 400 {
		return []Domain{}, errors.New("linode api unathorized")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &jsonObjs)
	if err != nil {
		fmt.Printf("Can't Unmarshal %s\n", err)
	}
	return jsonObjs.Data, err
}

// GetRecordResults get result object from domain and page number
func GetRecordResults(domain Domain, page int) (RecordResult, error) {
	var jsonObjs RecordResult
	cli := &http.Client{}
	uri := fmt.Sprintf("https://api.linode.com/v4/domains/%d/records/?page=%d", domain.ID, page)
	log.Printf("REQ=%v\n", uri)
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "Bearer "+APIToken)
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println("Error loading API", err)
		return RecordResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return RecordResult{}, errors.New("linode api server error")
	}
	if resp.StatusCode >= 400 {
		return RecordResult{}, errors.New("linode api unathorized")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &jsonObjs)
	if err != nil {
		fmt.Printf("Can't Unmarshal %s\n", err)
	}
	log.Printf("pages=%d page=%d results=%d\n", jsonObjs.Pages, jsonObjs.Page, jsonObjs.Results)
	// log.Println(jsonObjs)
	return jsonObjs, err
}

// AddRecord create new record on linode
func AddRecord(r Record, d Domain) (Record, error) {
	var jsonObjs Record
	cli := &http.Client{}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("Can't Marshal %s\n", err)
	}
	uri := fmt.Sprintf("https://api.linode.com/v4/domains/%d/records/", d.ID)
	log.Printf("REQ=%v\n", uri)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+APIToken)
	resp, err := cli.Do(req)
	if err != nil {
		return r, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &jsonObjs)
	if err != nil {
		fmt.Printf("Can't Unmarshal %s\n", err)
	}
	defer resp.Body.Close()
	return jsonObjs, nil
}

// RemoveRecord remove record from linode
func RemoveRecord(r Record, d Domain) error {
	cli := &http.Client{}
	data, err := json.Marshal(r)
	if err != nil {
		log.Printf("Can't Marshal %s\n", err)
	}
	uri := fmt.Sprintf("https://api.linode.com/v4/domains/%d/records/%d/", d.ID, r.ID)
	log.Printf("REQ=%v\n", uri)
	req, err := http.NewRequest("DELETE", uri, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+APIToken)
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// GetDomainObject get domain object from domain name
func GetDomainObject(domain string) (Domain, error) {
	domains, err := GetDomains()
	if err != nil {
		return Domain{}, err
	}
	for _, dom := range domains {
		if dom.Domain == domain {
			return dom, nil
		}
	}
	return Domain{}, errors.New("Domain not found")
}

// RemoveRecords Remove a list of Records from domain
func RemoveRecords(records []Record, domainObj Domain) {
	for _, rec := range records {
		err := RemoveRecord(rec, domainObj)
		if err != nil {
			log.Println(err)
		}
	}
}

// RemoveRecordByName Remove a txt record from linode
func RemoveRecordByName(record string, domain string) error {
	domainObj, err := GetDomainObject(domain)
	if err != nil {
		return err
	}
	result, err := GetRecordResults(domainObj, 1)
	if err != nil {
		return err
	}
	records := result.Data
	log.Printf("records page=%d len=%v cap=%v\n", 1, len(records), cap(records))
	for page := 1; page <= result.Page; {
		page = result.Page + 1
		result, err := GetRecordResults(domainObj, page)
		if err != nil {
			return err
		}
		for _, rec := range result.Data {
			records = append(records, rec)
		}
		log.Printf("records page=%d len=%v cap=%v\n", page, len(records), cap(records))
	}
	for _, r := range records {
		if r.Name == record {
			log.Println("Record Found", r)
			return RemoveRecord(r, domainObj)
		}
	}
	return errors.New("record not found")
}

// CreateNewTXTRecord create new record
func CreateNewTXTRecord(domain string, name string, value string) (Record, error) {
	domainObj, err := GetDomainObject(domain)
	if err != nil {
		return Record{}, err
	}
	rec := Record{Type: "TXT", Name: name, Target: value, TTLSec: 300}
	rec, err = AddRecord(rec, domainObj)
	if err != nil {
		return rec, err
	}

	return rec, nil
}
