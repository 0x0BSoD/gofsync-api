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
	Results  map[string][]PuppetClass `json:"results"`
	Total    int                      `json:"total"`
	SubTotal int                      `json:"subtotal"`
	Page     int                      `json:"page"`
	PerPage  int                      `json:"per_page"`
	Search   string                   `json:"search"`
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
func getAllPC(host string) (map[string][]PuppetClass, error) {

	var pcResult PuppetClasses
	var result = make(map[string][]PuppetClass)

	// check items count
	uri := fmt.Sprintf("puppetclasses?format=json&per_page=%d", globConf.PerPage)
	bodyText, err := ForemanAPI("GET", host, uri, "")
	if err == nil {
		err := json.Unmarshal(bodyText, &pcResult)
		if err != nil {
			log.Fatalf("%q:\n %s\n", err, bodyText)
		}

		if pcResult.Total > globConf.PerPage {
			pagesRange := Pager(pcResult.Total)
			for i := 1; i <= pagesRange; i++ {
				uri := fmt.Sprintf("puppetclasses?format=json&page=%d&per_page=%d", i, globConf.PerPage)
				bodyText, err := ForemanAPI("GET", host, uri, "")
				if err == nil {
					err := json.Unmarshal(bodyText, &pcResult)
					if err != nil {
						return result, err
					}

					for className, cl := range pcResult.Results {
						for _, subClass := range cl {
							result[className] = append(result[className], subClass)
						}
					}
				}
			}
		} else {
			for className, cl := range pcResult.Results {
				for _, subClass := range cl {
					result[className] = append(result[className], subClass)
				}
			}
		}
	}

	return result, nil
}

// Get Puppet Classes by host group and insert to Host Group
func getPCByHg(host string, hgID int, bdId int64) {
	var result PuppetClasses
	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	bodyText, err := ForemanAPI("GET", host, uri, "")
	if err == nil {
		err := json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Fatalf("%q:\n %s\n", err, bodyText)
		}
		var pcIDs []int64
		for className, cl := range result.Results {
			for _, sublcass := range cl {
				lastId := insertPC(host, className, sublcass.Name, sublcass.ID)
				if lastId != -1 {
					pcIDs = append(pcIDs, lastId)
				}
			}
		}
		updatePCinHG(bdId, pcIDs)
	}

}

// ===============
// INSERT
// ===============
