package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
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
	user := cfg.Api.Username
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

		base := models.HostGroupBase{
			Name:           t.Name,
			EnvironmentId:  envID,
			LocationIds:    locationsIDs,
			ParentId:       pID,
			PuppetClassIds: data.BaseInfo.PuppetClassIds,
		}
		jDataBase, _ := json.Marshal(models.POSTStructBase{HostGroup: base})
		response, err := logger.ForemanAPI("POST", params["host"], "hostgroups", string(jDataBase), cfg)
		fmt.Println(string(response.Body))
		fmt.Println(response.StatusCode)
		if response.StatusCode == 201 {
			if len(data.Overrides) > 0 {
				for _, ovr := range data.Overrides {

					p := struct {
						Match string `json:"match"`
						Value string `json:"value"`
					}{Match: ovr.Match, Value: ovr.Value}
					d := models.POSTStructOvrVal{p}
					jDataOvr, _ := json.Marshal(d)
					uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
					resp, err := logger.ForemanAPI("POST", params["host"], uri, string(jDataOvr), cfg)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						err = json.NewEncoder(w).Encode(fmt.Sprintf("Foreman Api Error: %q", err))
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							logger.Error.Printf("Error on POST HG: %s", err)
						}
					}
					logger.Info.Printf("%s : created Override ForemanID: %d on %s", user, ovr.ScForemanId, params["host"])
					logger.Trace.Println(string(resp.Body))

				}
			}
			err = json.NewEncoder(w).Encode(base)
			if err != nil {
				logger.Error.Printf("Error on Create HG: %s", err)
			}
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(fmt.Sprintf("Error on Create HG: %s", response.Body)))
			logger.Error.Printf("Error on Create HG: %s", string(response.Body))
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error on Create HG: %s", err)))
		logger.Error.Printf("Error on Create HG: %s", err)
	}

}

func Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cfg := middleware.GetConfig(r)
	user := cfg.Api.Username

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
	jDataBase, _ := json.Marshal(models.POSTStructBase{data.BaseInfo})

	// Hoist group not exist on target ====================================================
	if data.ExistId == -1 {
		response, err := logger.ForemanAPI("POST", t.TargetHost, "hostgroups", string(jDataBase), cfg)
		if err == nil && response.StatusCode == 200 {
			if len(data.Overrides) > 0 {
				for _, ovr := range data.Overrides {

					// Socket Broadcast ---
					msg := models.Step{
						Host:    t.TargetHost,
						Actions: "Submitting overrides",
						State:   fmt.Sprintf("Parameter: %s", ovr.Value),
					}
					utils.BroadCastMsg(cfg, msg)
					// ---

					p := struct {
						Match string `json:"match"`
						Value string `json:"value"`
					}{Match: ovr.Match, Value: ovr.Value}
					d := models.POSTStructOvrVal{p}
					jDataOvr, _ := json.Marshal(d)
					uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
					resp, err := logger.ForemanAPI("POST", t.TargetHost, uri, string(jDataOvr), cfg)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						err = json.NewEncoder(w).Encode(fmt.Sprintf("Foreman Api Error: %q", err))
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							logger.Error.Printf("Error on POST HG: %s", err)
						}
					}
					logger.Info.Println(string(resp.Body), resp.RequestUri)
				}
				if user != "" {
					logger.Info.Printf("crated HG || %s : %s on %s", user, data.BaseInfo.Name, t.TargetHost)
				}
				err = json.NewEncoder(w).Encode(string(response.Body))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Error.Printf("Error on POST HG: %s", err)
				}
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error on POST HG: %s", err)))
			logger.Error.Printf("Error on POST HG: %s", err)
			logger.Error.Printf("Error on POST HG: %s", string(response.Body))
		}
	} else {
		uri := fmt.Sprintf("hostgroups/%d", data.ExistId)
		response, err := logger.ForemanAPI("PUT", t.TargetHost, uri, string(jDataBase), cfg)
		if err == nil {
			if len(data.Overrides) > 0 {
				for _, ovr := range data.Overrides {

					// Socket Broadcast ---
					msg := models.Step{
						Host:    t.TargetHost,
						Actions: "Updating overrides",
						State:   fmt.Sprintf("Parameter: %s", ovr.Value),
					}
					utils.BroadCastMsg(cfg, msg)
					// ---

					p := struct {
						Match string `json:"match"`
						Value string `json:"value"`
					}{Match: ovr.Match, Value: ovr.Value}
					d := models.POSTStructOvrVal{p}
					jDataOvr, _ := json.Marshal(d)

					if ovr.OvrForemanId != -1 {
						uri := fmt.Sprintf("smart_class_parameters/%d/override_values/%d", ovr.ScForemanId, ovr.OvrForemanId)

						resp, err := logger.ForemanAPI("PUT", t.TargetHost, uri, string(jDataOvr), cfg)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							err = json.NewEncoder(w).Encode(fmt.Sprintf("Foreman Api Error: %q", err))
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								logger.Error.Printf("Error on POST HG: %s", err)
							}
						}
						if resp.StatusCode == 404 {
							uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
							resp, err := logger.ForemanAPI("POST", t.TargetHost, uri, string(jDataOvr), cfg)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								err = json.NewEncoder(w).Encode(fmt.Sprintf("Foreman Api Error: %q", err))
								if err != nil {
									w.WriteHeader(http.StatusInternalServerError)
									logger.Error.Printf("Error on POST HG: %s", err)
								}
							}
							if user != "" {
								logger.Info.Printf("%s : created Override ForemanID: %d on %s", user, ovr.ScForemanId, t.TargetHost)
							} else {
								logger.Info.Printf("NOPE : created Override ForemanID: %d on %s", ovr.ScForemanId, t.TargetHost)
							}
							logger.Trace.Println(string(resp.Body))
						}
						if user != "" {
							logger.Info.Printf("NOPE : updated Override ForemanID: %d on %s", ovr.ScForemanId, t.TargetHost)
						} else {

						}
						if user != "" {
							logger.Info.Printf("%s : updated Override ForemanID: %d, Value: %s on %s", user, ovr.ScForemanId, ovr.Value, t.TargetHost)
						} else {
							logger.Info.Printf("NOPE : updated Override ForemanID: %d, Value: %s on %s", ovr.ScForemanId, ovr.Value, t.TargetHost)
						}
					} else {
						uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
						resp, err := logger.ForemanAPI("POST", t.TargetHost, uri, string(jDataOvr), cfg)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							err = json.NewEncoder(w).Encode(fmt.Sprintf("Foreman Api Error: %q", err))
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								logger.Error.Printf("Error on POST HG: %s", err)
							}
						}
						logger.Info.Printf("%s : created Override ForemanID: %d on %s", user, ovr.ScForemanId, t.TargetHost)
						logger.Trace.Println(string(resp.Body))
					}
				}
			}
		}

		if user != "" {
			logger.Info.Printf("updated HG || %s : %s on %s", user, data.BaseInfo.Name, t.TargetHost)
		}

		// Socket Broadcast ---
		msg := models.Step{
			Host:    t.TargetHost,
			Actions: "Uploading Done!",
		}
		utils.BroadCastMsg(cfg, msg)
		// ---

		err = json.NewEncoder(w).Encode(string(response.Body))
		if err != nil {
			logger.Error.Printf("Error on PUT HG: %s", err)
		}
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
