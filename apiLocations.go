package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ===============================
// TYPES & VARS
// ===============================
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

// ===============
// GET
// ===============
func locations(host string) (Locations, error) {

	var result Locations
	bodyText := ForemanAPI("GET", host, "locations", "")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		return Locations{}, err
	}

	return result, nil
}

func locationsByHG(host string, hgID int, lastID int64) {

	var result Locations

	uri := fmt.Sprintf("hostgroups/%d/locations", hgID)
	bodyText := ForemanAPI("GET", host, uri, "")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}
	var locList []int64
	for _, loc := range result.Results {
		lId := checkLoc(host, loc.Name)
		locList = append(locList, int64(lId))
	}
	updateLocInHG(lastID, locList)
}
