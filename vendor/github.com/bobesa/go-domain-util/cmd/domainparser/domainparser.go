package main // import "github.com/bobesa/go-domain-util/cmd/domainparser"

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// tld contains single tld info
type tld map[string]*tld

func (t *tld) Source() string {
	// Report nil if nothing is present
	if len(*t) == 0 {
		return "nil"
	}

	// Create set of keys (for sorting)
	keys, i := make([]string, len(*t)), 0
	for key := range *t {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	// Create source code based on sorted keys
	str := "&tld{\n"
	for _, key := range keys {
		str += `"` + key + `": `
		str += (*t)[key].Source()
		str += ",\n"
	}
	return str + "}"
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

func main() {
	// Get path as argument
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Do the http request
	resp, err := http.Get("https://publicsuffix.org/list/public_suffix_list.dat")
	checkError(err)

	// Read the listing from request body
	b, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	// Generate basic tree
	tlds := &tld{}

	// Parse text as separate lines
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if line != "" && (len(line) < 2 || line[:2] != "//") {
			currentTLD := tlds
			parts := strings.Split(line, ".")
			for p := len(parts) - 1; p >= 0; p-- {
				part := parts[p]
				if nextTLD, exists := (*currentTLD)[part]; !exists {
					nextTLD = &tld{}
					(*currentTLD)[part] = nextTLD
					currentTLD = nextTLD
				} else {
					currentTLD = nextTLD
				}
			}
		}
	}

	// Create tlds file
	source := `package domainutil

	// tld contains single tld info
	type tld map[string]*tld

	// tlds holds all informations about correct tlds
	var tlds = ` + tlds.Source()

	// Run gofmt to format the code
	cmd := exec.Command("gofmt")
	cmd.Stdin = strings.NewReader(source)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// Write results
	dir := filepath.Dir(args[0])
	outputName := filepath.Join(dir, "tlds.go")
	err = ioutil.WriteFile(outputName, out, 0644)
	checkError(err)
}
