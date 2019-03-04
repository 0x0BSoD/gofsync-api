package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func postRTAPI(host string, params string) []byte {

	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer transport.CloseIdleConnections()

	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("POST", "http://"+host+"/"+params, nil)
	req.SetBasicAuth(Config.Username, Config.Pass)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}

	return []byte(bodyText)
}

func getForemanAPI(host string, params string) []byte {

	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer transport.CloseIdleConnections()

	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", "https://"+host+"/api/v2/"+params, nil)
	req.SetBasicAuth(Config.Username, Config.Pass)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}

	return []byte(bodyText)
}
