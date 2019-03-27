package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ===============================
// TYPES & VARS
// ===============================
// Smart Class Container
type SCParameters struct {
	Total    int           `json:"total"`
	SubTotal int           `json:"subtotal"`
	Page     int           `json:"page"`
	PerPage  int           `json:"per_page"`
	Search   string        `json:"search"`
	Results  []SCParameter `json:"results"`
}

// Smart Class
type SCParameter struct {
	Parameter           string      `json:"parameter"`
	PuppetClassName     string      `json:"puppetclass_name"`
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
type PCSCParameters struct {
	ID                   int             `json:"id"`
	Name                 string          `json:"name"`
	ModuleName           string          `json:"module_name"`
	SmartClassParameters []PCSCParameter `json:"smart_class_parameters"`
	Environments         []Environment   `json:"environments"`
	HostGroups           []HostGroupS    `json:"hostgroups"`
}
type PCSCParameter struct {
	ID          int    `json:"id"`
	Name        string `json:"parameter"`
	PuppetClass string `json:"puppetclass"`
}

// Return From Base
type SCGetRes struct {
	ForemanID int
	ID        int
	Type      string
}

// Return From Base
type SCGetResAdv struct {
	ID                  int
	ForemanId           int
	Name                string
	OverrideValuesCount int
	ValueType           string
	DefaultVal          interface{}
	Override            []SCOParams
}
type SCOParams struct {
	SmartClassId int    `json:"smart_class_id"`
	Parameter    string `json:"parameter"`
	Match        string `json:"match"`
	Value        string `json:"value"`
}

// ===============
// GET
// ===============

// ===============
// INSERT
// ===============
// Get Smart Classes from Foreman
func smartClasses(host string) ([]SCParameter, error) {

	var r SCParameters
	var result []SCParameter

	uri := fmt.Sprintf("smart_class_parameters?per_page=%d", globConf.PerPage)
	body := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, body)
	}

	if r.Total > globConf.PerPage {
		pagesRange := Pager(r.Total)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("smart_class_parameters?page=%d&per_page=%d", i, globConf.PerPage)
			body := ForemanAPI("GET", host, uri, "")
			err := json.Unmarshal(body, &r)
			if err != nil {
				return []SCParameter{}, err
			}
			for _, j := range r.Results {
				result = append(result, j)
			}
		}
	} else {
		for _, i := range r.Results {
			result = append(result, i)
		}
	}
	return result, nil
}

// Get Smart Classes Overrides from Foreman
func scOverridesById(host string, ForemanID int) []OverrideValue {

	var r OverrideValues
	var result []OverrideValue

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
				result = append(result, j)
			}
		}
	} else {
		for _, k := range r.Results {
			result = append(result, k)
		}
	}
	return result
}

//Update Smart Class ids in Puppet Classes
func insertSCByPC(host string) {
	var r PCSCParameters
	PCss := getAllPCBase(host)
	for _, ss := range PCss {
		uri := fmt.Sprintf("puppetclasses/%d", ss.ForemanID)
		bodyText := ForemanAPI("GET", host, uri, "")

		err := json.Unmarshal(bodyText, &r)
		if err != nil {
			log.Fatalf("%q:\n %s\n", err, bodyText)
		}
		updatePC(host, ss.SubClass, r)
	}
}
