package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
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

//func RTAPI(method string, host string, params string) []byte {
//
//	var res *http.Response
//	transport := makeTransport()
//	client := &http.Client{Transport: transport}
//	defer transport.CloseIdleConnections()
//
//	switch method {
//	case "GET":
//		req, _ := http.NewRequest(method, "http://"+host+"/"+params, nil)
//		req.SetBasicAuth(globConf.Username, globConf.Pass)
//		res, _ = client.Do(req)
//	}
//	bodyText, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		log.Fatalf("%q:\n %s\n", err, bodyText)
//	}
//	return []byte(bodyText)
//}

func ForemanAPI(method string, host string, params string, payload string) ([]byte, error) {

	var res *http.Response

	transport := makeTransport()
	client := &http.Client{Transport: transport}
	defer transport.CloseIdleConnections()

	uri := fmt.Sprintf("https://%s/api/v2/%s", host, params)

	switch method {
	case "GET":
		req, _ := http.NewRequest("GET", uri, nil)
		req.SetBasicAuth(globConf.Username, globConf.Pass)
		res, _ = client.Do(req)
	case "POST":
		req, _ := http.NewRequest("POST", uri, strings.NewReader(payload))
		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth(globConf.Username, globConf.Pass)
		res, _ = client.Do(req)
	case "DELETE":
		req, _ := http.NewRequest("DELETE", uri, nil)
		req.SetBasicAuth(globConf.Username, globConf.Pass)
		res, _ = client.Do(req)
	case "PUT":
		req, _ := http.NewRequest("PUT", uri, strings.NewReader(payload))
		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth(globConf.Username, globConf.Pass)
		res, _ = client.Do(req)
	}

	if res != nil {
		bodyText, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []byte{}, fmt.Errorf("host: %s, statusCode: %d, uri: %s", host, res.StatusCode, res.Request.RequestURI)
		}
		defer res.Body.Close()

		return []byte(bodyText), nil
	}

	return []byte{}, fmt.Errorf("error in apiWrap, %s", params)
}
