package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func getAPI(host string, params string) []byte {

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
	auth := configParser("./config.json")
	req.SetBasicAuth(auth.Username, auth.Pass)
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
