package smartclass

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/middleware"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// ===============================
// GET
// ===============================
func GetOverridesByHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	session := middleware.GetConfig(r)
	params := mux.Vars(r)
	data := GetOverridesHG(params["hgName"], &session)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting overrides: %s", err)
	}
}
func GetOverridesByLocHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	session := middleware.GetConfig(r)
	params := mux.Vars(r)
	data := GetOverridesLoc(params["locName"], params["host"], &session)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting location overrides: %s", err)
	}
}

func GetSCDataByIdHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	session := middleware.GetConfig(r)
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["sc_id"])
	data := GetSCData(id, &session)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting overrides: %s", err)
	}
}
