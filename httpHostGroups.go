package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// ===============================
// TYPES & VARS
// ===============================
type HGElem struct {
	ID            int                           `json:"id"`
	Name          string                        `json:"name"`
	Environment   string                        `json:"environment"`
	ParentId      string                        `json:"parent_id"`
	Params        []HGParam                     `json:"params,omitempty"`
	PuppetClasses map[string][]PuppetClassesWeb `json:"puppet_classes"`
}
type HGListElem struct {
	ID   int
	Name string
}
type HGParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type PC struct {
	ID        int
	ForemanId int
	Class     string
	Subclass  string
	SCIDs     string
}
type PuppetClassesWeb struct {
	Subclass     string      `json:"subclass"`
	SmartClasses []string    `json:"smart_classes,omitempty"`
	Overrides    []SCOParams `json:"overrides,omitempty"`
}
type HGPost struct {
	SourceHost string `json:"source_host"`
	TargetHost string `json:"target_host"`
	HgId       int    `json:"hg_id"`
}

// For POST HG
// POST /api/hostgroups
//{
//  "hostgroup": {
//    "name": "TestHostgroup",
//    "puppet_proxy_id": 182953976
//  }
//}
type hgPOSTParams struct {
	Name           string   `json:"name"`
	Locations      []string `json:"locations"`
	ParentId       int      `json:"parent_id"`
	EnvironmentId  int      `json:"environment_id"`
	PuppetClassIds []int    `json:"puppetclass_ids"`
}

// For POST Override Params
// POST /api/smart_class_parameters/:smart_class_parameter_id/override_values
//{
//  "override_value": {
//    "match": "domain=example.com",
//    "value": "gdkWwYbkrO"
//  }
//}
type SCOerridePOSTParams struct {
	Match string `json:"match"`
	Value string `json:"value"`
}

// ===============================
// GET
// ===============================
func getHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	data := getHGList(params["host"])
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting SWE list: %s", err)
	}
}

func getHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["swe_id"])
	data := getHG(params["host"], id)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting SWE list: %s", err)
	}
}

func postHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t HGPost
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatalf("Error on POST HG: %s", err)
	}
	data, err := postHG(t.SourceHost, t.TargetHost, t.HgId)
	if err != nil {
		type errStruct struct {
			Message string
			State   string
		}
		err = json.NewEncoder(w).Encode(errStruct{Message: err.Error(), State: "fail"})
		if err != nil {
			log.Fatalf("Error on getting SWE list: %s", err)
		}
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting SWE list: %s", err)
	}
}
