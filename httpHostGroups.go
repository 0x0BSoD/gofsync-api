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
	Class    string
	Subclass string
	SCIDs    string
}
type PuppetClassesWeb struct {
	Subclass     string      `json:"subclass"`
	SmartClasses []string    `json:"smart_classes,omitempty"`
	Overrides    []SCOParams `json:"overrides,omitempty"`
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
