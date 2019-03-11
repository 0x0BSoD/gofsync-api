package main

import (
	"encoding/json"
	"log"
	"sort"
	"strings"
)

// PuppetClasses container
type Locations struct {
	Results  []*Location            `json:"results"`
	Total    int                    `json:"total"`
	SubTotal int                    `json:"subtotal"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
	Search   string                 `json:"search"`
	Sort     map[string]interface{} `json:"sort"`
}

// PuppetClass structure
type Location struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func getLocations(host string) {

	var result Locations
	//fmt.Printf("Getting from %s \n", host)
	bodyText := ForemanAPI("GET", host, "locations", "")

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
