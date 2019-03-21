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
type envCheckP struct {
	Host string `json:"host"`
	Env  string `json:"env"`
}

// ===============================
// GET
// ===============================
func getAllEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	data := getEnvList(params["host"])
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting HG list: %s", err)
	}
}

// ===============================
// POST
// ===============================
func postEnvCheckHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t envCheckP
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatalf("Error on POST EnvCheck: %s", err)
	}
	data := checkEnv(t.Host, t.Env)
	if err != nil {
		err = json.NewEncoder(w).Encode(errStruct{Message: err.Error(), State: "fail"})
		if err != nil {
			log.Fatalf("Error on EnvCheck: %s", err)
		}
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on EnvCheck: %s", err)
	}
}
