package main

import (
	"encoding/json"
	"fmt"
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
	ID   int    `json:"id"`
	Name string `json:"name"`
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
type errStruct struct {
	Message string
	State   string
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
		log.Fatalf("Error on getting HG list: %s", err)
	}
}

func getAllHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := getHGAllList()
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting All HG list: %s", err)
	}
}

func getHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["swe_id"])
	data := getHG(id)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting SWE list: %s", err)
	}
}

func getAllHostsHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(globConf.Hosts)
	if err != nil {
		log.Fatalf("Error on getting SWE list: %s", err)
	}
}

// ===============================
// POST
// ===============================
func postHGCheckHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t HGPost
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatalf("Error on POST HG: %s", err)
	}
	data := postCheckHG(t.TargetHost, t.HgId)
	if err != nil {
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

func postHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t HGPost
	type POSTStructBase struct {
		HostGroup HostGroupBase `json:"hostgroup"`
	}
	type POSTStructOvrVal struct {
		OverrideValue struct {
			Match string `json:"match"`
			Value string `json:"value"`
		} `json:"override_value"`
	}
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatalf("Error on POST HG: %s", err)
	}
	data, err := postHG(t.SourceHost, t.TargetHost, t.HgId)
	jDataBase, _ := json.Marshal(POSTStructBase{data.BaseInfo})
	response, err := ForemanAPI("POST", t.TargetHost, "hostgroups", string(jDataBase))
	if err == nil {
		if len(data.Overrides) > 0 {
			for _, ovr := range data.Overrides {

				p := struct {
					Match string `json:"match"`
					Value string `json:"value"`
				}{Match: ovr.Match, Value: ovr.Value}

				d := POSTStructOvrVal{p}
				jDataOvr, _ := json.Marshal(d)
				uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ForemanId)
				resp, err := ForemanAPI("POST", t.TargetHost, uri, string(jDataOvr))
				if err != nil {
					err = json.NewEncoder(w).Encode(string(resp))
					if err != nil {
						log.Fatalf("Error on getting SWE list: %s", err)
					}
				}
			}
		}
		err = json.NewEncoder(w).Encode(string(response))
		if err != nil {
			log.Fatalf("Error on getting SWE list: %s", err)
		}
	}
}
