package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
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
	Parameter           string          `json:"parameter"`
	PuppetClass         PuppetClassInSc `json:"puppetclass"`
	ID                  int             `json:"id"`
	Description         string          `json:"description"`
	Override            bool            `json:"override"`
	ParameterType       string          `json:"parameter_type"`
	DefaultValue        interface{}     `json:"default_value"`
	UsePuppetDefault    bool            `json:"use_puppet_default"`
	Required            bool            `json:"required"`
	ValidatorType       string          `json:"validator_type"`
	ValidatorRule       string          `json:"validator_rule"`
	MergeOverrides      bool            `json:"merge_overrides"`
	AvoidDuplicates     bool            `json:"avoid_duplicates"`
	OverrideValueOrder  string          `json:"override_value_order"`
	OverrideValuesCount int             `json:"override_values_count"`
}

// PC for old Foremans
type PuppetClassInSc struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ModuleName string `json:"module_name"`
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
type logStatus struct {
	Name          string `json:"name"`
	Host          string `json:"host"`
	Current       int    `json:"current"`
	CurrentThread int    `json:"current_thread,omitempty"`
	TotalInThread int    `json:"total_in_thread,omitempty"`
	Total         int    `json:"total"`
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
	var resultId []int
	var result []SCParameter

	uri := fmt.Sprintf("smart_class_parameters?per_page=%d", globConf.PerPage)
	body, _ := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, body)
	}

	if r.Total > globConf.PerPage {
		pagesRange := Pager(r.Total)
		for i := 1; i <= pagesRange; i++ {

			//jsonLog, _ := json.Marshal(logStatus{
			//	Name:    "smart_class_parameters_id",
			//	Host:    host,
			//	Current: i,
			//	Total:   r.Total,
			//})
			//fmt.Println(string(jsonLog))

			uri := fmt.Sprintf("smart_class_parameters?page=%d&per_page=%d", i, globConf.PerPage)
			body, _ := ForemanAPI("GET", host, uri, "")
			err := json.Unmarshal(body, &r)
			if err != nil {
				return []SCParameter{}, err
			}
			for _, j := range r.Results {
				resultId = append(resultId, j.ID)
			}
		}
	} else {
		for _, i := range r.Results {
			resultId = append(resultId, i.ID)
		}
	}
	queue := splitToQueue(resultId, 6)
	var d SCParameter
	var wg sync.WaitGroup

	for tIdx, q := range queue {
		wg.Add(1)
		go func(tIdx int, q []int) {
			defer wg.Done()
			for _, sId := range q {
				//jsonLog, _ := json.Marshal(logStatus{
				//	Name:          "smart_class_parameters",
				//	Host:          host,
				//	Current:       idx,
				//	CurrentThread: tIdx,
				//	Total:         len(resultId),
				//	TotalInThread: len(q),
				//})
				//fmt.Println(string(jsonLog))

				uri := fmt.Sprintf("smart_class_parameters/%d", sId)
				body, _ := ForemanAPI("GET", host, uri, "")
				err := json.Unmarshal(body, &d)
				if err != nil {
					log.Printf("Error on getting override: %q \n%s\n", err, uri)
				} else {
					result = append(result, d)
				}
			}
		}(tIdx, q)
	}
	wg.Wait()
	return result, nil
}

// Get Smart Classes Overrides from Foreman
func scOverridesById(host string, ForemanID int) []OverrideValue {

	var r OverrideValues
	var result []OverrideValue

	uri := fmt.Sprintf("smart_class_parameters/%d/override_values?per_page=%d", ForemanID, globConf.PerPage)
	body, _ := ForemanAPI("GET", host, uri, "")
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, body)
	}

	if r.Total > globConf.PerPage {
		pagesRange := Pager(r.Total)
		for i := 1; i <= pagesRange; i++ {

			uri := fmt.Sprintf("smart_class_parameters/%d/override_values?page=%d&per_page=%d", ForemanID, i, globConf.PerPage)
			body, _ := ForemanAPI("GET", host, uri, "")
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
func smartClassByPC(host string) {
	var r PCSCParameters
	PCss := getAllPCBase(host)
	for _, ss := range PCss {

		//jsonLog, _ := json.Marshal(logStatus{
		//	Name:    "smart_class_foreman_id_" + ss.SubClass,
		//	Host:    host,
		//	Current: idx,
		//	Total:   len(PCss),
		//})
		//fmt.Println(string(jsonLog))

		uri := fmt.Sprintf("puppetclasses/%d", ss.ForemanID)
		bodyText, _ := ForemanAPI("GET", host, uri, "")

		err := json.Unmarshal(bodyText, &r)
		if err != nil {
			log.Fatalf("%q:\n %s\n", err, bodyText)
		}
		updatePC(host, ss.SubClass, r)
	}
}

func smartClassByPCJson(host string, pcId int) []SCParameter {

	var r SCParameters

	uri := fmt.Sprintf("puppetclasses/%d/smart_class_parameters", pcId)
	bodyText, _ := ForemanAPI("GET", host, uri, "")

	err := json.Unmarshal(bodyText, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}
	return r.Results
}

// ===
func smartClassByPCJsonV2(host string, pcId int) PCSCParameters {

	var r PCSCParameters

	uri := fmt.Sprintf("puppetclasses/%d", pcId)
	bodyText, _ := ForemanAPI("GET", host, uri, "")

	err := json.Unmarshal(bodyText, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}
	return r
}
func smartClassByFId(host string, foremanId int) SCParameter {
	var r SCParameter

	uri := fmt.Sprintf("smart_class_parameters/%d", foremanId)
	bodyText, _ := ForemanAPI("GET", host, uri, "")

	err := json.Unmarshal(bodyText, &r)
	if err != nil {
		log.Fatalf("%q:\n %s\n", err, bodyText)
	}
	return r
}
