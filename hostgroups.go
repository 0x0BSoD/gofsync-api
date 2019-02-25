package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func getHostgroups(host string) {

	spaces := 10
	var result []entitys.SWEs

	fmt.Printf("Getting from %s \n", host)

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
	req, _ := http.NewRequest("GET", "https://"+host+"/api/hostgroups?format=json&per_page="+count, nil)
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

	for _, item := range result {
		if item.Hostgroup.Name != "SWE" {
			fmt.Println(host + "  ==================================================")
			fmt.Println("Name             :  ", item.Hostgroup.Name)
			fmt.Println("ID               :  ", item.Hostgroup.ID)
			fmt.Println("SubnetID         :  ", item.Hostgroup.SubnetID)
			fmt.Println("OperatingsystemID:  ", item.Hostgroup.OperatingsystemID)
			fmt.Println("DomainID         :  ", item.Hostgroup.DomainID)
			fmt.Println("EnvironmentID    :  ", item.Hostgroup.EnvironmentID)
			fmt.Println("Ancestry         :  ", item.Hostgroup.Ancestry)

			if len(item.Hostgroup.Parameters) > 1 {
				fmt.Println("Parameters       :=>  ")
				for name, item := range item.Hostgroup.Parameters {
					length := len(name)
					strSpaces := giveMeSpaces(spaces - length)
					fmt.Printf("    %s%s:  %s\n", name, strSpaces, item)
				}
			} else {
				fmt.Println("Parameters       :   ", nil)
			}
			fmt.Println("PuppetclassIds   :  ", item.Hostgroup.PuppetclassIds)

			sJson, _ := json.Marshal(item.Hostgroup)

			if insertToSWE(item.Hostgroup.Name, host, string(sJson)) {
				fmt.Println("  ==================================================")
				fmt.Println(item.Hostgroup.Name + "  INSERTED")
				fmt.Println("  ==================================================")
			}
			fmt.Println()
		}
	}
}
