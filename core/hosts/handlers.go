package hosts

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/global"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
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

// Get HG info from Foreman
func ByHostgroupHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	response := ByHostgroup(params["host"], params["hgForemanId"], ctx)
	if response.StatusCode == 404 {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode("not found"); err != nil {
			utils.Error.Printf("Error on getting HG: %s", err)
		}
	} else {
		var result Hosts
		if err := json.Unmarshal(response.Body, &result); err != nil {
			utils.Error.Printf("Error on getting HG: %s", err)
		}
		if err := json.NewEncoder(w).Encode(result.Results); err != nil {
			utils.Error.Printf("Error on getting HG: %s", err)
		}
	}
}

func Update(w http.ResponseWriter, r *http.Request) {

	// VARS
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	global.Sync(params["host"], ctx)

	// =====
	if err := json.NewEncoder(w).Encode("ok"); err != nil {
		utils.Error.Printf("error while updating host: %s", err)
	}
}

//func ByHostgroupNameHttp(w http.ResponseWriter, _ *http.Request) {
// swagger:operation GET /hosts/all/hg/{hgName} host Host
//
// Returns a hosts list with target SWE
// ---
// consumes:
// - text/json
// produces:
// - text/json
// parameters:
// - name: hgName
//   in: path
//   description: HostGroup name for search.
//   required: true
//   type: string
// responses:
//   '200':
//     description: Host list
//     type: string
//w.Header().Set("Content-Type", "application/json")
//ctx := middleware.GetContext(r)
//params := mux.Vars(r)
//if err := r.ParseForm(); err != nil {
//	utils.Warning.Printf("Error on parsing parameters: %s", err)
//}
//if _, ok := r.Form["hostnames"]; ok {
//	data := ByHostgroupNameHostNames(params["hgName"], r.Form, ctx)
//	if err := json.NewEncoder(w).Encode(data); err != nil {
//		utils.Error.Printf("Error on getting HG: %s", err)
//	}
//} else {
//	data := ByHostgroupName(params["hgName"], r.Form, ctx)
//	if err := json.NewEncoder(w).Encode(data); err != nil {
//		utils.Error.Printf("Error on getting HG: %s", err)
//	}
//}
//}
