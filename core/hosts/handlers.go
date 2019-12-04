package hosts

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// ===============================
// GET
// ===============================

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

func NewHostHttp(ctx *user.GlobalCTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// VARS
		ctx.Set(&user.Claims{Username: "srv_foreman"}, "fake")
		//ctx := middleware.GetContext(r)

		var b NewHostParams
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&b)
		if err != nil {
			utils.Error.Printf("Error on POST NewHostHttp: %s", err)
		}

		fmt.Println(b)

		//name, foremanHost, envName, locName, hgName string, ctx *user.GlobalCTX
		response, err := CreateNewHost(b, ctx)
		if err != nil {
			utils.Error.Printf("Error on POST NewHostHttp: %s", err)
		}

		// ==========
		utils.SendResponse(w, "error while creating new host: %s", string(response))
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
