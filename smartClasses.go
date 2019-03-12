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

//// OverrideValues
//type OverrideValues struct {
//	ID               int         `json:"id"`
//	Match            string      `json:"match"`
//	Value            interface{} `json:"value"`
//	UsePuppetDefault bool        `json:"use_puppet_default"`
//}
//
//// PClass
//type PClass struct {
//	Name       string `json:"name"`
//	ModuleMame string `json:"module_name"`
//	ID         int    `json:"id"`
//}

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
				insertSC(host, j)
			}
		}
	} else {
		for _, i := range r.Results {
			//fmt.Printf("SC Param: %s || %s\n", i.Parameter, host)
			insertSC(host, i)
		}
	}
}
