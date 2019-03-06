package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/foremanGetter/entitys"
	"log"
	"strconv"
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
		lastId := insertToSWE(i.Name, host, string(sJson))
		if lastId != -1 {
			getSWEParams(host, count, lastId, i.ID)
		}
	}

	//return  pPrintCommit(result, host)
}

func getSWEParams(host string, count string, dbID int64, sweID int) {
	var data entitys.SWEParameterContainer

	bodyText := getForemanAPI(host, "hostgroups/"+strconv.Itoa(sweID)+"/parameters?format=json&per_page="+count)

	err := json.Unmarshal(bodyText, &data)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}

	for _, i := range data.Results {
		insertToSWEParams(dbID, i.Name, i.Value, i.Priority)
		fmt.Println(i.Name)
		fmt.Println(i.Value)
		fmt.Println()
	}
}
