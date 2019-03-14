package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// ===============================
// TYPES & VARS
// ===============================
type HGElem struct {
	ID            int
	Name          string
	Params        []HGParam
	PuppetClasses map[string][]PuppetClassesWeb
}
type HGListElem struct {
	ID   int
	Name string
}
type HGParam struct {
	Name  string
	Value string
}
type PC struct {
	Class    string
	Subclass string
	SCIDs    string
	//EnvIDs string
	//HGIDs string
}
type PuppetClassesWeb struct {
	Subclass     string      `json:"subclass"`
	SmartClasses []string    `json:"smart_classes"`
	Overrides    []SCOParams `json:"overrides"`
	//EnvIds []int `json:"env_ids"`
	//HostGroupsIds []int `json:"host_groups_ids"`
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
	data := getHG(params["host"], params["swe_id"])
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting SWE list: %s", err)
	}
}
