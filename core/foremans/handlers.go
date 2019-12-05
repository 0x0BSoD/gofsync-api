package foremans

import (
	"encoding/json"
	"github.com/gorilla/mux"

	//"git.ringcentral.com/archops/goFsync/core/global"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/utils"
	"net/http"
)

// ===============================
// GET
// ===============================

func GetAllHostsHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	data := PuppetHosts(ctx)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting hosts: %s", err)
	}
}

func GetTrendDataHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ctx := middleware.GetContext(r)

	data := getTrends(ctx.Config.Hosts[params["host"]], ctx)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting hosts: %s", err)
	}
}

//func Update(w http.ResponseWriter, r *http.Request) {
//	// VARS
//	ctx := middleware.GetContext(r)
//	params := mux.Vars(r)
//
//	global.Sync(params["host"], ctx)
//
//	// =====
//	if err := json.NewEncoder(w).Encode("ok"); err != nil {
//		utils.Error.Printf("error while updating host: %s", err)
//	}
//}
