package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var url = "https://www.haitao.mei:60443/"
	req(url)

	// using localhost will raise an error, because the certificate is valid for *.haitao.mei, not localhost
	url = "https://www.haitao.mei:60443/helloAPI/Leopold"
	reqWithCert(url)
}

func req(url string) {
	fmt.Println("Request without cert....")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(contents))
}

func reqWithCert(url string) {
	fmt.Println("Request with cert....")
	// Load client cert
	// cert, err := tls.LoadX509KeyPair("certFile_location", "keyFile_location")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Load CA cert
	caCert, err := ioutil.ReadFile("tls/ca-crt.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		//Certificates: []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false, /* check if the cert is valid */
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(contents))
}
