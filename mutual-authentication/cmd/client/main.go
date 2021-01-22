package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// GO 1.15版本废弃了Common Name
// 因此运行时需要加GODEBUG=x509ignoreCN=0运行环境
func main() {
	cert, err := tls.LoadX509KeyPair("../../client-certs/clientCrt.pem", "../../client-certs/clientKey.pem")
	if err != nil {
		log.Fatalln(err)
	}
	serverCert, err := ioutil.ReadFile("../../server-certs/serverCrt.pem")
	if err != nil {
		log.Fatalln(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(serverCert)

	tlsConfig := &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	resp, err := client.Get("https://localhost:9443/hello")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(body))
}
