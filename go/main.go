package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	name := strings.ToLower(os.Args[1])

	// Get list of all current TLDs.
	// https://data.iana.org/TLD/tlds-alpha-by-domain.txt
	resp, err := http.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		log.Fatalf("Failed to fetch https://data.iana.org/TLD/tlds-alpha-by-domain.txt: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read list of TLDs: %v", err)
	}

	r := strings.Split(string(body), "\n")[1:]
	var tlds []string
	// My kingdom for pipes, filter, etc.
	for _, tld := range r {
		// Clean the various XN--NNNNN entries.
		if di := strings.Index(tld, "-"); di > 0 {
			tld = tld[0:di]
		}
		tld = strings.ToLower(tld)
		if len(tld) == 0 {
			continue
		}
		if len(tlds) == 0 || tlds[len(tlds)-1] != tld {
			tlds = append(tlds, tld)
		}
	}

	for _, tld := range tlds {
		if i := strings.LastIndex(name, tld); i > 0 {
			if i == 0 {
				// Can't have an empty domain; skip.
				continue
			}
			if i+len(tld) == len(name) {
				fmt.Printf("%s.%s\n", name[0:i], tld)
			} else {
				fmt.Printf("%s.%s/%s\n", name[0:i], tld, name[i+len(tld):])
			}
		}
	}
}
