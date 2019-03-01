package main

import (
	"encoding/json"
	"github.com/0x0bsod/foremanGetter/entitys"
	"log"
	"sort"
	"strings"
)

func getLocations(host string) {
	var result entitys.Locations
	//fmt.Printf("Getting from %s \n", host)
	bodyText := getAPI(host, "locations")

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
	//fmt.Println(listLocations)

	insertToLocations(host, strings.Join(listLocations, ","))
}
