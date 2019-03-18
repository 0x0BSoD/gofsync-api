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
type PuppetClasses struct {
	Results  map[string][]*PuppetClass `json:"results"`
	Total    int                       `json:"total"`
	SubTotal int                       `json:"subtotal"`
	Page     int                       `json:"page"`
	PerPage  int                       `json:"per_page"`
	Search   string                    `json:"search"`
}

// PuppetClass structure
type PuppetClass struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	SmartClassesId []int
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ===============
// GET
// ===============
// Get all Puppet Classes and insert to base
func getAllPC(host string) {

	var r PuppetClasses
	uri := fmt.Sprintf("puppetclasses?format=json&per_page=%d", globConf.PerPage)
	bodyText := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(bodyText, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}

	if r.Total > globConf.PerPage {
		pagesRange := Pager(r.Total)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("puppetclasses?format=json&page=%d&per_page=%d", i, globConf.PerPage)
			bodyText := ForemanAPI("GET", host, uri, "")
			err := json.Unmarshal(bodyText, &r)
			if err != nil {
				log.Fatalf("%q:\n %s\n", err, bodyText)
			}

			for className, cl := range r.Results {
				for _, sublcass := range cl {
					insertPC(host, className, sublcass.Name, sublcass.ID)
				}
			}
		}
	} else {
		for className, cl := range r.Results {
			for _, sublcass := range cl {
				insertPC(host, className, sublcass.Name, sublcass.ID)
			}
		}
	}
}

//
func getPCIdOnHost(host string, PCname string) int {
	var r PuppetClass
	uri := fmt.Sprintf("puppetclasses/%s", PCname)
	fmt.Println(uri)
	body := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, body)
	}
	return r.ID
}

// ===============
// INSERT
// ===============
// Get Puppet Classes by host group and insert to Host Group
func insertPCByHg(host string, hgID int, bdId int64) {
	var result PuppetClasses
	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	bodyText := ForemanAPI("GET", host, uri, "")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}
	var pcIDs []int64
	for className, cl := range result.Results {
		for _, sublcass := range cl {
			lastId := insertPC(host, className, sublcass.Name, sublcass.ID)
			fmt.Printf("PC: %s, %d || %s\n", sublcass.Name, lastId, host)
			if lastId != -1 {
				pcIDs = append(pcIDs, lastId)
			}
		}
	}
	updatePCinHG(bdId, pcIDs)
}
