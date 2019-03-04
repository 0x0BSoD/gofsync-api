package main

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/foremanGetter/entitys"
	"log"
)

func getRTHostGroups(host string) {
	var result []entitys.RTSWE

	bodyText := postRTAPI(host, "api/rchwswelookups/search?q=name~.*&fields=name,osversion,basetpl,swestatus&format=json")
	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	for _, i := range result {
		insertToSWE(i.Name, host, "[]")
	}

}

func getHostGroups(host string, count string) {
	var data entitys.SWEs

	bodyText := getForemanAPI(host, "hostgroups?format=json&per_page="+count+"&search=label+~+SWE%2F")

	err := json.Unmarshal(bodyText, &data)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}

	for _, i := range data.Results {
		sJson, _ := json.Marshal(i)
		insertToSWE(i.Name, host, string(sJson))
	}

	//return  pPrintCommit(result, host)

}
