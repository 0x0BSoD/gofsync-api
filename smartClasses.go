package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ===============================
// TYPES & VARS
// ===============================
// Smart Class
type SCParameter struct {
	Parameter           string      `json:"parameter"`
	ID                  int         `json:"id"`
	Description         string      `json:"description"`
	Override            bool        `json:"override"`
	ParameterType       string      `json:"parameter_type"`
	DefaultValue        interface{} `json:"default_value"`
	UsePuppetDefault    bool        `json:"use_puppet_default"`
	Required            bool        `json:"required"`
	ValidatorType       string      `json:"validator_type"`
	ValidatorRule       string      `json:"validator_rule"`
	MergeOverrides      bool        `json:"merge_overrides"`
	AvoidDuplicates     bool        `json:"avoid_duplicates"`
	OverrideValueOrder  string      `json:"override_value_order"`
	OverrideValuesCount int         `json:"override_values_count"`
	//CreatedAt           string            `json:"created_at"`
	//UpdatedAt           string            `json:"updated_at"`
	//PuppetClass         *PClass           `json:"puppetclass"`
	//OverrideValues      []*OverrideValues `json:"override_values"`
}

// Smart Class Container
type SCParameters struct {
	Total    int           `json:"total"`
	SubTotal int           `json:"subtotal"`
	Page     int           `json:"page"`
	PerPage  int           `json:"per_page"`
	Search   string        `json:"search"`
	Results  []SCParameter `json:"results"`
}

// OverrideValues Container
type OverrideValues struct {
	Total    int             `json:"total"`
	SubTotal int             `json:"subtotal"`
	Page     int             `json:"page"`
	PerPage  int             `json:"per_page"`
	Search   string          `json:"search"`
	Results  []OverrideValue `json:"results"`
}
type OverrideValue struct {
	ID               int         `json:"id"`
	Match            string      `json:"match"`
	Value            interface{} `json:"value"`
	UsePuppetDefault bool        `json:"use_puppet_default"`
}

// Return From Base
type SCGetRes struct {
	ForemanID int
	ID        int
	Type      string
}

// ===============
// GET
// ===============
// Get Smart Classes from Foreman
func getSmartClasses(host string) {
	var r SCParameters
	uri := fmt.Sprintf("smart_class_parameters?per_page=%d", globConf.PerPage)
	body := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, body)
	}

	if r.Total > globConf.PerPage {
		pagesRange := Pager(r.Total)
		for i := 1; i <= pagesRange; i++ {

			fmt.Printf("SC Param Page: %d of %d || %s\n", i, pagesRange, host)

			uri := fmt.Sprintf("smart_class_parameters?page=%d&per_page=%d", i, globConf.PerPage)
			body := ForemanAPI("GET", host, uri, "")
			err := json.Unmarshal(body, &r)
			if err != nil {
				log.Fatalf("%q:\n %s\n", err, body)
			}
			for _, j := range r.Results {
				//fmt.Printf("SC Param: %s || %s\n", j.Parameter, host)
				lastID := insertSC(host, j)
				if lastID != -1 {
					getSCOverridesById(host, j.ID, lastID, j.ParameterType)
				}
			}
		}
	} else {
		for _, i := range r.Results {
			//fmt.Printf("SC Param: %s || %s\n", i.Parameter, host)
			lastID := insertSC(host, i)
			if lastID != -1 {
				getSCOverridesById(host, i.ID, lastID, i.ParameterType)
			}
		}
	}
}

// Get Smart Classes Overrides from Foreman
func getSCOverridesById(host string, ForemanID int, ID int64, pType string) {
	var r OverrideValues
	uri := fmt.Sprintf("smart_class_parameters/%d/override_values?per_page=%d", ForemanID, globConf.PerPage)
	body := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, body)
	}

	if r.Total > globConf.PerPage {
		pagesRange := Pager(r.Total)
		for i := 1; i <= pagesRange; i++ {

			fmt.Printf("SC Param Page: %d of %d || %s\n", i, pagesRange, host)

			uri := fmt.Sprintf("smart_class_parameters/%d/override_values?page=%d&per_page=%d", ForemanID, i, globConf.PerPage)
			body := ForemanAPI("GET", host, uri, "")
			err := json.Unmarshal(body, &r)
			if err != nil {
				log.Fatalf("%q:\n %s\n", err, body)
			}

			for _, j := range r.Results {
				insertSCOverride(ID, j, pType)
			}
		}
	} else {
		for _, k := range r.Results {
			insertSCOverride(ID, k, pType)
		}
	}
}

// Get Smart Classes Overrides from Foreman
func getSCOverrides(host string) {
	data := getSCWithOverrides(host)
	var r OverrideValues
	items := len(data)
	for i := 0; i < items; i++ {
		//https://spb01-puppet.lab.nordigy.ru/api/v2/smart_class_parameters/173/override_values
		uri := fmt.Sprintf("smart_class_parameters/%d/override_values?per_page=%d", data[i].ForemanID, globConf.PerPage)
		body := ForemanAPI("GET", host, uri, "")
		err := json.Unmarshal(body, &r)
		if err != nil {
			log.Fatalf("%q:\n %s\n", err, body)
		}

		if r.Total > globConf.PerPage {
			pagesRange := Pager(r.Total)
			for i := 1; i <= pagesRange; i++ {

				fmt.Printf("SC Param Page: %d of %d || %s\n", i, pagesRange, host)

				uri := fmt.Sprintf("smart_class_parameters/%d/override_values?page=%d&per_page=%d", data[i].ForemanID, i, globConf.PerPage)
				body := ForemanAPI("GET", host, uri, "")
				err := json.Unmarshal(body, &r)
				if err != nil {
					log.Fatalf("%q:\n %s\n", err, body)
				}

				for _, j := range r.Results {
					//fmt.Printf("SC Override: %s || %s\n", j.Match, host)
					insertSCOverride(int64(data[i].ID), j, data[i].Type)
				}
			}
		} else {
			for _, k := range r.Results {
				//fmt.Printf("SC Override: %s || %s\n", k.Match, host)
				insertSCOverride(int64(data[i].ID), k, data[i].Type)
			}
		}
	}
}
