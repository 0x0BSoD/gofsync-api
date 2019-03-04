package main

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/foremanGetter/entitys"
	"log"
	"sort"
	"strings"
)

func getLocations(host string) {

	var result entitys.Locations
	//fmt.Printf("Getting from %s \n", host)
	bodyText := getForemanAPI(host, "locations")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}
	var listLocations []string
	for _, loc := range result.Results {
		listLocations = append(listLocations, strings.ToUpper(loc.Name))
	}
	sort.Strings(listLocations)
	insertToLocations(host, strings.Join(listLocations, ","))
}
