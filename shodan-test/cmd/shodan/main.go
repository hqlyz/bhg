package main

import (
	"fmt"
	"log"
	"os"
	"shodan-test/shodan"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: shodan searchterm")
	}
	c := shodan.New("noED5KJvpANoc2xmm1ni5ERFmLRREDZT")
	apiInfo, err := c.APIInfo()
	if err != nil {
		log.Fatalf("api-info error: %v\n", err)
	}
	fmt.Printf("scan credits: %d\tquery credits: %d\n", apiInfo.ScanCredits, apiInfo.QueryCredits)
	hostSearch, err := c.HostSearch(os.Args[1])
	if err != nil {
		log.Fatalf("host-search error: %v\n", err)
	}
	for _, host := range hostSearch.Matches {
		fmt.Printf("%18s%8d\n", host.IPStr, host.Port)
	}
}
