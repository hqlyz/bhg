package main

import (
	"fmt"
	"log"
	"net/http"
	"plugin-scanner/scanner"
)


var Users = []string{"admin", "manager"}
var Passwords = []string{"admin", "manager"}

var (
	res *scanner.Result
	req *http.Request
	resp *http.Response
	client *http.Client
	err error
)

type TomcatChecker struct {}

func (checker *TomcatChecker)Check(host string, port uint64) *scanner.Result {
	log.Println("Checking for Tomcat Manager...")
	res = new(scanner.Result)
	url := fmt.Sprintf("http://%s:%d/manager/html", host, port)
	if resp, err = http.Head(url); err != nil {
		fmt.Printf("http head error: %v\n", err)
		return res
	}
	log.Println("Host responded to /manager/html request")
	if resp.StatusCode != http.StatusUnauthorized || resp.Header.Get("WWW-Authenticate") == "" {
		fmt.Println("No need to authorize")
		return res
	}
	log.Println("Host requires authentication. Proceeding with password guessing...")
	client = new(http.Client)
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Println("Unable to build GET request")
 		return res 
	}
	for _, u := range Users {
		for _, p := range Passwords {
			req.SetBasicAuth(u, p)
			if resp, err = client.Do(req); err != nil {
				log.Printf("Http request failed: %v\n", err)
				continue
			}
			if resp.StatusCode == http.StatusOK {
				res.Vulnerable = true
				res.Details = fmt.Sprintf("Valid credentials found - %s:%s", u, p)
				return res
			}
		}
	}
	return res
}

func New() scanner.Checker {
	return new(TomcatChecker)
}