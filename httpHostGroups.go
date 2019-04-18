package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
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
func getHGFHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		data, err := hostGroupJson(params["host"], params["hgName"], cfg)
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

func getHGCheckHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		data := hostGroupCheck(params["host"], params["hgName"], cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting HG check: %s", err)
		}
	}
}

func getHGUpdateInBaseHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		hostGroup(params["host"], params["hgName"], cfg)
		err := json.NewEncoder(w).Encode("ok")
		if err != nil {
			logger.Error.Printf("Error on updating HG: %s", err)
		}
	}
}

func getHGListHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		data := getHGList(params["host"], cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting HG list: %s", err)
		}
	}
}

func getAllHGListHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := getHGAllList(cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting all HG list: %s", err)
		}
	}
}

func getHGHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["swe_id"])
		data := getHG(id, cfg)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting HG: %s", err)
		}
	}
}

func getAllHostsHttp(cfg *models.Config) http.HandlerFunc {
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
func postHGCheckHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		var t models.HGPost
		err := decoder.Decode(&t)
		if err != nil {
			logger.Error.Printf("Error on POST HG: %s", err)
		}
		data := postCheckHG(t.TargetHost, t.SourceHgId, cfg)
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error.Printf("Error on getting SWE list: %s", err)
		}
	}
}

func postHGHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		var t models.HGPost
		err := decoder.Decode(&t)
		if err != nil {
			logger.Error.Printf("Error on POST HG: %s", err)
			return
		}

		data, err := postHG(t.SourceHost, t.TargetHost, t.SourceHgId, cfg)
		if err != nil {
			logger.Error.Printf("Error on POST HG: %s", err)
			return
		}

		jDataBase, _ := json.Marshal(models.POSTStructBase{data.BaseInfo})

		if data.BaseInfo.ExistId == -1 {
			response, err := logger.ForemanAPI("POST", t.TargetHost, "hostgroups", string(jDataBase), cfg)
			if err == nil {
				if len(data.Overrides) > 0 {
					for _, ovr := range data.Overrides {

						fmt.Println(ovr)

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
							err = json.NewEncoder(w).Encode("Foreman Api Error - 500")
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								logger.Error.Printf("Error on POST HG: %s", err)
							}
						}
						logger.Info.Println(string(resp.Body), resp.RequestUri)
					}
				}
				// Commit new HG for target host
				hostGroup(t.TargetHost, data.BaseInfo.Name, cfg)
			}
			user := context.Get(r, 0)
			if user != nil {
				logger.Info.Printf("%s : %s on %s data: %s", user.(string), "uploaded HG data", t.TargetHost, string(response.Body))
			}
			err = json.NewEncoder(w).Encode(string(response.Body))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error.Printf("Error on POST HG: %s", err)
				return
			}
		} else {
			uri := fmt.Sprintf("hostgroups/%d", data.BaseInfo.ExistId)
			response, err := logger.ForemanAPI("PUT", t.TargetHost, uri, string(jDataBase), cfg)
			if err == nil {
				if len(data.Overrides) > 0 {
					for _, ovr := range data.Overrides {

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
								err = json.NewEncoder(w).Encode("Foreman Api Error - 500")
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
									err = json.NewEncoder(w).Encode("Foreman Api Error - 500")
									if err != nil {
										w.WriteHeader(http.StatusInternalServerError)
										logger.Error.Printf("Error on POST HG: %s", err)
									}
								}
								logger.Info.Println(string(resp.Body))
							}
							logger.Info.Println(string(resp.Body))
						} else {
							uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
							resp, err := logger.ForemanAPI("POST", t.TargetHost, uri, string(jDataOvr), cfg)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								err = json.NewEncoder(w).Encode("Foreman Api Error - 500")
								if err != nil {
									w.WriteHeader(http.StatusInternalServerError)
									logger.Error.Printf("Error on POST HG: %s", err)
								}
							}
							logger.Info.Println(string(resp.Body))
						}
					}
				}
				// Commit new HG for target host
				hostGroup(t.TargetHost, data.BaseInfo.Name, cfg)
			}

			user := context.Get(r, 0)
			if user != nil {
				logger.Info.Printf("%s : %s on %s data: %s", user.(string), "updated HG data", t.TargetHost, string(response.Body))
			}
			err = json.NewEncoder(w).Encode(string(response.Body))
			if err != nil {
				logger.Error.Printf("Error on PUT HG: %s", err)
				return
			}
		}
	}
}
