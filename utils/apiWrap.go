package utils

import (
	"crypto/tls"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
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

func ForemanAPI(method, host, params, payload string, ctx *user.GlobalCTX) (models.Response, error) {

	var res *http.Response

	transport := makeTransport()
	defer transport.CloseIdleConnections()

	uri := fmt.Sprintf("https://%s/api/v2/%s", host, params)
	var req *http.Request
	var err error
	switch method {
	case "GET":
		req, err = http.NewRequest("GET", uri, nil)
	case "POST":
		req, err = http.NewRequest("POST", uri, strings.NewReader(payload))
	case "DELETE":
		req, err = http.NewRequest("DELETE", uri, nil)
	case "PUT":
		req, err = http.NewRequest("PUT", uri, strings.NewReader(payload))
	}

	// Request
	if req != nil {
		client := &http.Client{Transport: transport}
		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth(ctx.Config.Api.Username, ctx.Config.Api.Password)
		res, err = client.Do(req)
		if err != nil {
			return models.Response{
				StatusCode: -1,
				Body:       []byte{},
				RequestUri: host,
			}, nil
		}
	} else {
		Error.Printf("error in apiWrap, %s", params)
		return models.Response{}, fmt.Errorf("error in apiWrap, %s", params)
	}

	// Response
	if res != nil {
		bodyText, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return models.Response{}, fmt.Errorf("host: %s, statusCode: %d, uri: %s", host, res.StatusCode, res.Request.RequestURI)
		}
		defer res.Body.Close()

		return models.Response{
			StatusCode: res.StatusCode,
			Body:       bodyText,
			RequestUri: res.Request.RequestURI,
		}, nil
	}

	Error.Printf("error in apiWrap, %s", params)
	return models.Response{
		StatusCode: res.StatusCode,
		Body:       []byte{},
		RequestUri: res.Request.RequestURI,
	}, nil
}

func RackTablesAPI(method, host, params, payload string, ctx *user.GlobalCTX) (models.Response, error) {
	var res *http.Response

	transport := makeTransport()
	defer transport.CloseIdleConnections()

	uri := fmt.Sprintf("https://%s/api/%s", host, params)
	var req *http.Request
	var err error
	switch method {
	case "GET":
		req, err = http.NewRequest("GET", uri, nil)
	case "POST":
		req, err = http.NewRequest("POST", uri, strings.NewReader(payload))
	case "DELETE":
		req, err = http.NewRequest("DELETE", uri, nil)
	case "PUT":
		req, err = http.NewRequest("PUT", uri, strings.NewReader(payload))
	}

	// Request
	if req != nil {
		client := &http.Client{Transport: transport}
		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth(ctx.Config.Api.Username, ctx.Config.Api.Password)
		res, err = client.Do(req)
		if err != nil {
			return models.Response{
				StatusCode: -1,
				Body:       []byte{},
				RequestUri: host,
			}, nil
		}
	} else {
		Error.Printf("error in apiWrap, %s", params)
		return models.Response{}, fmt.Errorf("error in apiWrap, %s", params)
	}

	// Response
	if res != nil {
		bodyText, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return models.Response{}, fmt.Errorf("host: %s, statusCode: %d, uri: %s", host, res.StatusCode, res.Request.RequestURI)
		}
		defer res.Body.Close()

		return models.Response{
			StatusCode: res.StatusCode,
			Body:       bodyText,
			RequestUri: res.Request.RequestURI,
		}, nil
	}

	Error.Printf("error in apiWrap, %s", params)
	return models.Response{
		StatusCode: res.StatusCode,
		Body:       []byte{},
		RequestUri: res.Request.RequestURI,
	}, nil
}
