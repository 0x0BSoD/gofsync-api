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
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ===============
// GET
// ===============
// Get Puppet Classes by host group
func getPCByHg(host string, hgID int) {

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
			lastId := insertPC(host, className, sublcass.Name)
			fmt.Printf("PC: %s, %d || %s\n", sublcass.Name, lastId, host)
			if lastId != -1 {
				pcIDs = append(pcIDs, lastId)
			}
		}
	}
	updatePCinHG(hgID, pcIDs)
}
