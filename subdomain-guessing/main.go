package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

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
			ips = append(ips, a.A.String())
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
		results = append(results, result{IPAddress: ip, Hostname: fqdn})
	}
	return results
}

type empty struct{}

func worker(tracker chan empty, fqdns chan string, gather chan []result, serverAddr string) {
	for fqdn := range fqdns {
		result := loopup(fqdn, serverAddr)
		if len(result) > 0 {
			gather <- result
		}
	}
	var e empty
	tracker <- e
}

func main() {
	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("-domain and -wordlist are required")
		os.Exit(1)
	}
	fmt.Println(*flWorkerCount, *flServerAddr)

	var results []result
	fqdns := make(chan string, *flWorkerCount)
	gather := make(chan []result)
	tracker := make(chan empty)

	fh, err := os.Open(*flWordlist)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)

	for i := 0; i < *flWorkerCount; i++ {
		go worker(tracker, fqdns, gather, *flServerAddr)
	}

	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s.%s", scanner.Text(), *flDomain)
	}

	go func() {
		for g := range gather {
			results = append(results, g...)
		}
		var e empty
		tracker <- e
	}()

	close(fqdns)
	for i := 0; i < *flWorkerCount; i++ {
		<- tracker
	}
	close(gather)
	<- tracker

	w := tabwriter.NewWriter(os.Stdout, 8, 4, 0, ' ', 0)
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\n", r.Hostname, r.IPAddress)
	}
	w.Flush()
}
