package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getPuppetClasses(host string, classID int) {

	//spaces := 10
	var result entitys.PuppetClasses

	fmt.Printf("Getting %d class.\n", classID)
	//fmt.Println("https://" + host + "/api/v2/hostgroups/" + strconv.Itoa(classID) + "/puppetclasses")

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
	req, _ := http.NewRequest("GET", "https://"+host+"/api/v2/hostgroups/"+strconv.Itoa(classID)+"/puppetclasses", nil)
	auth := configParser("./config.json")
	req.SetBasicAuth(auth.Username, auth.Pass)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	bodyText, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal([]byte(bodyText), &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	for index, cl := range result.Results {
		fmt.Printf("%s ====\n", index)
		for _, v := range cl {
			fmt.Println("    ID          :  ", v.ID)
			fmt.Println("    Name        :  ", v.Name)
			fmt.Println("    CreatedAt   :  ", v.CreatedAt)
			fmt.Println("    UpdatedAt   :  ", v.UpdatedAt)
		}
	}
}
