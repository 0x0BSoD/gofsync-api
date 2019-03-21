package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func makeTransport() *http.Transport {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return transport
}

func RTAPI(method string, host string, params string) []byte {

	var res *http.Response
	transport := makeTransport()
	client := &http.Client{Transport: transport}
	defer transport.CloseIdleConnections()

	switch method {
	case "GET":
		req, _ := http.NewRequest(method, "http://"+host+"/"+params, nil)
		req.SetBasicAuth(globConf.Username, globConf.Pass)
		res, _ = client.Do(req)
	}
	bodyText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}
	return []byte(bodyText)
}

func ForemanAPI(method string, host string, params string, payload string) []byte {

	var res *http.Response
	transport := makeTransport()
	client := &http.Client{Transport: transport}
	defer transport.CloseIdleConnections()

	switch method {
	case "GET":
		req, _ := http.NewRequest("GET", "https://"+host+"/api/v2/"+params, nil)
		req.SetBasicAuth(globConf.Username, globConf.Pass)
		res, _ = client.Do(req)
	case "POST":
		req, _ := http.NewRequest("POST", "https://"+host+"/api/v2/"+params, strings.NewReader(payload))
		fmt.Println(req.Method)
		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth(globConf.Username, globConf.Pass)
		res, _ = client.Do(req)
	}
	bodyText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(res.Request.RequestURI)
		log.Fatalf("%s || %q:\n %s\n", host, err, bodyText)
	}
	if res.StatusCode == 500 {
		log.Println(res.Request.RequestURI)
		log.Fatalf("%s || %s\n", host, bodyText)
	}
	return []byte(bodyText)
}
