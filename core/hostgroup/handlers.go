package hostgroup

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/hostgroup/DB"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func GetHGHttp(w http.ResponseWriter, r *http.Request) {
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
