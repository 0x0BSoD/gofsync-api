package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ===============================
// TYPES & VARS
// ===============================
// For Getting SWE from RackTables
type RTSWE struct {
	Name      string `json:"name"`
	BaseTpl   string `json:"basetpl"`
	OsVersion string `json:"osversion"`
	SWEStatus string `json:"swestatus"`
}
type RTSWES []RTSWE

// For Getting SWE from Foreman
type SWE struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Title               string `json:"title"`
	SubnetID            int    `json:"subnet_id"`
	SubnetName          string `json:"subnet_name"`
	OperatingSystemID   int    `json:"operatingsystem_id"`
	OperatingSystemName string `json:"operatingsystem_name"`
	DomainID            int    `json:"domain_id"`
	DomainName          string `json:"domain_name"`
	EnvironmentID       int    `json:"environment_id"`
	EnvironmentName     string `json:"environment_name"`
	ComputeProfileId    int    `json:"compute_profile_id"`
	ComputeProfileName  string `json:"compute_profile_name"`
	Ancestry            string `json:"ancestry,omitempty"`
	PuppetProxyId       int    `json:"puppet_proxy_id"`
	PuppetCaProxyId     int    `json:"puppet_ca_proxy_id"`
	PTableId            int    `json:"ptable_id"`
	PTableName          string `json:"ptable_name"`
	MediumId            int    `json:"medium_id"`
	MediumName          string `json:"medium_name"`
	ArchitectureId      int    `json:"architecture_id"`
	ArchitectureName    int    `json:"architecture_name"`
	RealmId             int    `json:"realm_id"`
	RealmName           string `json:"realm_name"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
type SWEContainer struct {
	Results  []SWE  `json:"results"`
	Total    int    `json:"total"`
	SubTotal int    `json:"subtotal"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	Search   string `json:"search"`
}

//  Host Group parameters
type HostGroupPContainer struct {
	Results  []HostGroupP `json:"results"`
	Total    int          `json:"total"`
	SubTotal int          `json:"subtotal"`
	Page     int          `json:"page"`
	PerPage  int          `json:"per_page"`
	Search   string       `json:"search"`
}
type HostGroupP struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Priority int    `json:"priority"`
}
type HostGroupS struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

// ===============================
// GET
// ===============================
// Get SWE from RackTables
func (swe RTSWE) Get(host string) RTSWES {
	var r RTSWES
	body := RTAPI("GET", host,
		"api/rchwswelookups/search?q=name~.*&fields=name,osversion,basetpl,swestatus&format=json")

	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("%q:\n %s\n", err, body)
		return []RTSWE{}
	}
	return r
}

// ===================================
// Get SWE from Foreman
func (swe SWE) Get(host string) {
	var r SWEContainer
	uri := fmt.Sprintf("hostgroups?format=json&per_page=%d&search=label+~+SWE", globConf.PerPage)
	body := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, body)
	}

	if r.Total > globConf.PerPage {
		var resultsContainer []SWE
		pagesRange := Pager(r.Total)
		for i := 1; i <= pagesRange; i++ {

			uri := fmt.Sprintf("hostgroups?format=json&page=%d&per_page=%d&search=label+~+SWE", i, globConf.PerPage)
			body := ForemanAPI("GET", host, uri, "")
			err := json.Unmarshal(body, &r)
			if err != nil {
				log.Fatalf("%q:\n %s\n", err, body)
			}
			resultsContainer = append(resultsContainer, r.Results...)
		}
		for _, i := range resultsContainer {
			sJson, _ := json.Marshal(i)
			lastId := insertHG(i.Name, host, string(sJson), i.ID)
			if lastId != -1 {
				getPCByHg(host, i.ID, lastId)
				hgParams(host, lastId, i.ID)
				locationsByHG(host, i.ID, lastId)
			}
		}
	} else {
		for _, i := range r.Results {
			sJson, _ := json.Marshal(i)
			lastId := insertHG(i.Name, host, string(sJson), i.ID)
			if lastId != -1 {
				getPCByHg(host, i.ID, lastId)
				hgParams(host, lastId, i.ID)
				locationsByHG(host, i.ID, lastId)
			}
		}
	}
}

// Get SWE Parameters from Foreman
func hgParams(host string, dbID int64, sweID int) {
	var r HostGroupPContainer
	uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&per_page=%d", sweID, globConf.PerPage)
	body := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, body)
	}

	if r.Total > globConf.PerPage {
		pagesRange := Pager(r.Total)
		for i := 1; i <= pagesRange; i++ {

			uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&page=%d&per_page=%d", sweID, i, globConf.PerPage)
			body := ForemanAPI("GET", host, uri, "")
			err := json.Unmarshal(body, &r)
			if err != nil {
				log.Fatalf("%q:\n %s\n", err, body)
			}
			for _, j := range r.Results {
				insertHGP(dbID, j.Name, j.Value, j.Priority)
			}
		}
	} else {
		for _, i := range r.Results {
			insertHGP(dbID, i.Name, i.Value, i.Priority)
		}
	}
}
