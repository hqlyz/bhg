package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("hello %s\n", r.TLS.PeerCertificates[0].Subject.CommonName)
	fmt.Fprint(w, "Athentication Successful")
}

func main() {
	http.HandleFunc("/hello", hello)

	clientCerts, err := ioutil.ReadFile("../../client-certs/clientCrt.pem")
	if err != nil {
		log.Fatalln(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(clientCerts)

	tlsConfig := &tls.Config{
		ClientCAs:  pool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	server := &http.Server{
		TLSConfig: tlsConfig,
		Addr:      ":9443",
	}
	log.Println("server is starting...")
	log.Fatalln(server.ListenAndServeTLS("../../server-certs/serverCrt.pem", "../../server-certs/serverKey.pem"))
}
