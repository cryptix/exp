package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {

	io.WriteString(w, "hello, world!\n")
	for i, c := range req.TLS.PeerCertificates {
		fmt.Fprintf(w, "%2d: %s\n", i, c.SerialNumber.String())
	}
}

func main() {
	http.HandleFunc("/hello", HelloServer)

	caCert, err := ioutil.ReadFile("../certsNkeys/ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: tlsConfig,
	}

	log.Println("listening")
	err = server.ListenAndServeTLS("server.pem", "server-key.pem") //private cert
	if err != nil {
		log.Fatal(err)
	}
}
