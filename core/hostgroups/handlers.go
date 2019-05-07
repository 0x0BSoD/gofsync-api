package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// ===============================
// GET
// ===============================

// Get HG info from Foreman
func GetHGFHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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
}

func GetHGCheckHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		data := HostGroupCheck(params["host"], params["hgName"], cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting HG check: %s", err)
		}
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

func GetHGListHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		data := GetHGList(params["host"], cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting HG list: %s", err)
		}
	}
}

func GetAllHGListHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := GetHGAllList(cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting all HG list: %s", err)
		}
	}
}

func GetHGHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["swe_id"])
		data := GetHG(id, cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting HG: %s", err)
		}
	}
}

func GetAllHostsHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(cfg.Hosts)
		if err != nil {
			logger.Error.Printf("Error on getting hosts: %s", err)
		}
	}
}

// ===============================
// POST
// ===============================
func PostHGCheckHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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
}

func PostHGHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		user := context.Get(r, 0)
		var t models.HGPost
		err := decoder.Decode(&t)
		if err != nil {
			logger.Error.Printf("Error on POST HG: %s", err)
			return
		}

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

		if data.ExistId == -1 {

			fmt.Println(string(jDataBase))
			response, err := logger.ForemanAPI("POST", t.TargetHost, "hostgroups", string(jDataBase), cfg)
			fmt.Println(string(response.Body))

			if err == nil {
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
				}
			}
			if user != nil {
				logger.Info.Printf("%s : %s on %s", user.(string), "uploaded HG data", t.TargetHost)
			} else {
				logger.Info.Printf("%s : NOPE on %s", "uploaded HG data", t.TargetHost)
			}
			err = json.NewEncoder(w).Encode(string(response.Body))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error.Printf("Error on POST HG: %s", err)
				return
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
								if user != nil {
									logger.Info.Printf("%s : created Override ForemanID: %d on %s", user.(string), ovr.ScForemanId, t.TargetHost)
								} else {
									logger.Info.Printf("NOPE : created Override ForemanID: %d on %s", ovr.ScForemanId, t.TargetHost)
								}
								logger.Trace.Println(string(resp.Body))
							}
							if user != nil {
								logger.Info.Printf("NOPE : updated Override ForemanID: %d on %s", ovr.ScForemanId, t.TargetHost)
							} else {

							}
							if user != nil {
								logger.Info.Printf("%s : updated Override ForemanID: %d, Value: %s on %s", user.(string), ovr.ScForemanId, ovr.Value, t.TargetHost)
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
							logger.Info.Printf("%s : created Override ForemanID: %d on %s", user.(string), ovr.ScForemanId, t.TargetHost)
							logger.Trace.Println(string(resp.Body))
						}
					}
				}
			}

			user := context.Get(r, 0)
			if user != nil {
				logger.Info.Printf("%s : %s on %s", user.(string), "updated HG", t.TargetHost)
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
}
