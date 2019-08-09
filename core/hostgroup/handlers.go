package hostgroup

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/hostgroup/DB"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetHGByNameHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)
	data, err := gDB.ByName(params["host"], params["hgName"], ctx)
	if err != nil {
		utils.Error.Printf("Error on getting HG: %s", err)
	}

	// =========
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG: %s", err)
	}
}

func GetHGByIDHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)
	ID, _ := strconv.Atoi(params["hgID"])
	data, err := gDB.ByID(ID, ctx)
	if err != nil {
		utils.Error.Printf("Error on getting HG: %s", err)
	}

	// =========
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG: %s", err)
	}
}

func GetAllHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get

	data := gDB.List(ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting all HG list: %s", err)
	}
}

func GetHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)

	data := gDB.ListByHost(params["host"], ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting all HG list: %s", err)
	}
}
