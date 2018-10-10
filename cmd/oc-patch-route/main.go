package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	mylog "github.com/helber/letsencrypt-dns/log"
	flag "github.com/spf13/pflag"
)

// /opt/openshift/bin/oc patch route $ROUTE -p '{"spec":{"termination":"edge","tls":{"certificate":"'${CERT}'","key":"'${KEY}'"}}}'

// OpenshiftTLS fullcert and key
type OpenshiftTLS struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

// RouteSpec Openshift Router spec
type RouteSpec struct {
	Termination string       `json:"termination"`
	TLS         OpenshiftTLS `json:"tls"`
}

// OpenshiftPatch oc command route
type OpenshiftPatch struct {
	Spec RouteSpec `json:"spec"`
}

// GenerateOpenshiftPatch generate a struct to patch
func GenerateOpenshiftPatch(dirname string) (OpenshiftPatch, error) {
	full, err := ioutil.ReadFile(dirname + "fullchain1.pem")
	if err != nil {
		log.Fatalf("Error openning %v", err)
		return OpenshiftPatch{}, err
	}
	key, err := ioutil.ReadFile(dirname + "privkey1.pem")
	if err != nil {
		log.Fatalf("Error openning %v", err)
		return OpenshiftPatch{}, err
	}
	patch := OpenshiftPatch{
		Spec: RouteSpec{
			Termination: "edge",
			TLS: OpenshiftTLS{
				Certificate: string(full),
				Key:         string(key),
			},
		},
	}
	return patch, nil
}

// GetOcCommand return a oc command line
func GetOcCommand(namespace string, route string, patch OpenshiftPatch) (string, error) {
	j, err := json.Marshal(patch)
	if err != nil {
		log.Fatalf("Error %s", err)
		return "", err
	}
	occmd := fmt.Sprintf("oc patch route %s -n %s -p '%s'", route, namespace, j)
	return occmd, nil
}

func main() {
	domain := flag.StringP("domain", "d", "", "Domain name")
	certdir := flag.StringP("certdir", "c", ".", "Directory base of certificated")
	namespace := flag.StringP("namespace", "n", "", "Openshift namespace")
	route := flag.StringP("route", "r", "", "Route name")
	flag.Parse()
	mylog.InitLogs()

	// domain := "ahgoracloud.com.br"
	fullcertdir := "" + *certdir + *domain + "/"
	patch, err := GenerateOpenshiftPatch(fullcertdir)
	if err != nil {
		log.Fatalf("Error %s", err)
	}

	occmd, err := GetOcCommand(*namespace, *route, patch)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	fmt.Println(occmd)
}
