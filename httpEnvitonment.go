package main

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// ===============================
// GET
// ===============================
func getAllEnv(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		data := getEnvList(params["host"], cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting HG list: %s", err)
		}
	}
}

// ===============================
// POST
// ===============================
func postEnvCheckHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		var t models.EnvCheckP
		err := decoder.Decode(&t)
		if err != nil {
			logger.Error.Printf("Error on POST EnvCheck: %s", err)
		}
		data := checkEnv(t.Host, t.Env, cfg)
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on EnvCheck: %s", err)
		}
	}
}
