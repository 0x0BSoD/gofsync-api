package smartclass

import (
	"encoding/json"
	cl "git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// ===============================
// GET
// ===============================
func GetOverridesByHGHttp(cfg *cl.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		data := GetOverridesHG(params["hgName"], cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting overrides: %s", err)
		}
	}
}
func GetOverridesByLocHttp(cfg *cl.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		data := GetOverridesLoc(params["locName"], params["host"], cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting location overrides: %s", err)
		}
	}
}
