package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/miekg/dns"
)

var (
	flDomain      = flag.String("domain", "", "domain name to search")
	flWordlist    = flag.String("wordlist", "", "The wordlist to use for guessing.")
	flWorkerCount = flag.Int("c", 100, "The amount of workers to use.")
	flServerAddr  = flag.String("server", "8.8.8.8:53", "The DNS server to use.")
)

func init() {
	flag.Parse()
}

type result struct {
	IPAddress string
	Hostname  string
}

func lookupA(fqdn, dnsServer string) ([]string, error) {
	var msg dns.Msg
	var ips []string
	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	r, err := dns.Exchange(&msg, dnsServer)
	if err != nil {
		return ips, err
	}
	if len(r.Answer) < 1 {
		return ips, errors.New("no answers")
	}
	for _, answer := range r.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, string(a.A))
		}
	}
	return ips, nil
}

func lookupCNAME(fqdn, dnsServer string) ([]string, error) {
	var msg dns.Msg
	var fqdns []string
	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&msg, dnsServer)
	if err != nil {
		return fqdns, err
	}
	if len(in.Answer) < 1 {
		return fqdns, errors.New("no answers")
	}
	for _, answer := range in.Answer {
		if cname, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, cname.Target)
		}
	}
	return fqdns, nil
}

func loopup(fqdn, dnsServer string) []result {
	var results []result
	var cfqdn = fqdn
	for {
		cname, err := lookupCNAME(cfqdn, dnsServer)
		if err != nil {
			break
		}
		cfqdn = cname[0]
	}
	ips, err := lookupA(cfqdn, dnsServer)
	if err != nil {
		return results
	}
	for _, ip := range ips {
		results = append(results, result{IPAddress: ip, Hostname: cfqdn})
	}
	return results
}

func main() {
	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("-domain and -wordlist are required")
		os.Exit(1)
	}
	fmt.Println(*flWorkerCount, *flServerAddr)
}
