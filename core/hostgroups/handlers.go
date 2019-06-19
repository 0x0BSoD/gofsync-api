package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/models"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// ===============================
// GET
// ===============================

// Get HG info from Foreman
func GetHGFHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	params := mux.Vars(r)
	data, err := HostGroupJson(params["host"], params["hgName"], cfg)
	if (models.HgError{}) != err {
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			logger.Error.Printf("Error on getting HG: %s", err)
		}
	} else {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting HG: %s", err)
		}
	}
}

func GetHGCheckHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	params := mux.Vars(r)
	data := HostGroupCheck(params["host"], params["hgName"], cfg)
	if data.Error == "error -1" {
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("410 - Foreman server gone"))
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG check: %s", err)
	}
}

func GetHGUpdateInBaseHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		HostGroup(params["host"], params["hgName"], cfg)
		err := json.NewEncoder(w).Encode("ok")
		if err != nil {
			logger.Error.Printf("Error on updating HG: %s", err)
		}
	}
}

func GetHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	params := mux.Vars(r)
	data := GetHGList(params["host"], cfg)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetAllHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	data := GetHGAllList(cfg)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting all HG list: %s", err)
	}
}

func GetHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["swe_id"])
	data := GetHG(id, cfg)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG: %s", err)
	}
}

func GetAllHostsHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	data := AllHosts(cfg)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting hosts: %s", err)
	}
}

// ===============================
// POST
// ===============================
func PostHGCheckHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	decoder := json.NewDecoder(r.Body)
	var t models.HGPost
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
	}
	data := PostCheckHG(t.TargetHost, t.SourceHgId, cfg)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting SWE list: %s", err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	cfg.Web.RunSocket = false
	params := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	var t models.HGElem
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	envID := environment.CheckPostEnv(params["host"], t.Environment, cfg)
	locationsIDs := locations.DbAllForemanID(params["host"], cfg)
	pID, _ := strconv.Atoi(t.ParentId)

	if envID != -1 {

		data, _ := HGDataNewItem(params["host"], t, cfg)
		//
		base := models.HWPostRes{
			BaseInfo: models.HostGroupBase{
				Name:           t.Name,
				EnvironmentId:  envID,
				LocationIds:    locationsIDs,
				ParentId:       pID,
				PuppetClassIds: data.BaseInfo.PuppetClassIds,
			},
			Overrides:  data.Overrides,
			Parameters: data.Parameters,
		}
		resp, err := PushNewHG(base, params["host"], cfg)
		fmt.Println(resp)
		fmt.Println(base.Parameters)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Printf("Error on POST HG: %s", err)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
			return
		}
		// Send response to client
		cfg.Web.RunSocket = true
		_ = json.NewEncoder(w).Encode(resp)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Printf("Error on Create HG: %s", err)
		_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
	}

}

func Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cfg := middleware.GetConfig(r)
	decoder := json.NewDecoder(r.Body)
	var t models.HGPost
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	// Get data from DB ====================================================
	data, err := HGDataItem(t.SourceHost, t.TargetHost, t.SourceHgId, cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(fmt.Sprintf("Foreman Api Error: %q", err))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Printf("Error on POST HG: %s", err)
		}
		return
	}

	// Submit host group ====================================================
	if data.ExistId == -1 {
		resp, err := PushNewHG(data, t.TargetHost, cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Printf("Error on POST HG: %s", err)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
		}
		// Send response to client
		_ = json.NewEncoder(w).Encode(resp)
	} else {
		resp, err := UpdateHG(data, t.TargetHost, cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Printf("Error on PUT HG: %s", err)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on PUT HG: %s", err))
		}
		// Send response to client
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// ===============================
// PUT
// ===============================
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	params := mux.Vars(r)
	Sync(params["host"], cfg)
	err := json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}
