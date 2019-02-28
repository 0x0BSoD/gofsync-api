package main

import (
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
	"log"
	"sort"
	"strings"
)

func getLocations(host string) {
	var result entitys.Locations
	fmt.Printf("Getting from %s \n", host)
	bodyText := getAPI(host, "locations")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}
	var listLocations []string
	for _, loc := range result.Results {
		//fmt.Println("    ID          :  ", loc.ID)
		//fmt.Println("    Name        :  ", loc.Name)
		listLocations = append(listLocations, strings.ToUpper(loc.Name))
		//fmt.Println("    Title       :  ", loc.Title)
		//fmt.Println("    CreatedAt   :  ", loc.CreatedAt)
		//fmt.Println("    UpdatedAt   :  ", loc.UpdatedAt)
	}
	sort.Strings(listLocations)
	fmt.Println(listLocations)

	insertToLocations(host, strings.Join(listLocations, ","))
}
