package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
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
		req, _ := http.NewRequest(method, "https://"+host+"/api/v2/"+params, nil)
		req.SetBasicAuth(globConf.Username, globConf.Pass)
		res, _ = client.Do(req)
	}
	bodyText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}
	return []byte(bodyText)
}
