package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/smartclass/DB"
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

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)

	match := fmt.Sprintf("hostgroup=SWE/%s", params["match"])

	data := gDB.OverridesByMatch(params["host"], match, ctx)

	// ===
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("error while getting overrides: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
	}
}

func GetOverridesByLocHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)

	match := fmt.Sprintf("location=%s", params["match"])

	data := gDB.OverridesByMatch(params["host"], match, ctx)

	// ===
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("error while getting overrides: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
	}
}

func GetSCDataByIdHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["sc_id"])
	data, err := gDB.ByID(id, ctx)
	if err != nil {
		logger.Error.Printf("error while getting overrides: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
	}

	// ===
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Printf("error while getting overrides: %s", err)
	}
}
