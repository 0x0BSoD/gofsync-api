package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/logger"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// ===============================
// TYPES & VARS
// ===============================
type HGElem struct {
	ID            int                           `json:"id"`
	ForemanID     int                           `json:"foreman_id"`
	Name          string                        `json:"name"`
	Environment   string                        `json:"environment"`
	ParentId      string                        `json:"parent_id"`
	Params        []HGParam                     `json:"params,omitempty"`
	PuppetClasses map[string][]PuppetClassesWeb `json:"puppet_classes"`
}
type HGListElem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type HGParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type PC struct {
	ID        int
	ForemanId int
	Class     string
	Subclass  string
	SCIDs     string
}
type PuppetClassesWeb struct {
	Subclass     string      `json:"subclass"`
	SmartClasses []string    `json:"smart_classes,omitempty"`
	Overrides    []SCOParams `json:"overrides,omitempty"`
}
type HGPost struct {
	SourceHost string `json:"source_host"`
	TargetHost string `json:"target_host"`
	TargetHgId int    `json:"target_hg_id"`
	SourceHgId int    `json:"source_hg_id"`
}
type errStruct struct {
	Message string
	State   string
}
type POSTStructBase struct {
	HostGroup HostGroupBase `json:"hostgroup"`
}
type POSTStructOvrVal struct {
	OverrideValue struct {
		Match string `json:"match"`
		Value string `json:"value"`
	} `json:"override_value"`
}

// ===============================
// GET
// ===============================

// Get HG info from Foreman
func getHGFHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	data, err := hostGroupJson(params["host"], params["hgName"])
	if (errs{}) != err {
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			log.Fatalf("Error on getting HG list: %s", err)
		}
	} else {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Fatalf("Error on getting HG list: %s", err)
		}
	}

}
func getHGCheckHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	data := hostGroupCheck(params["host"], params["hgName"])
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting HG list: %s", err)
	}
}
func getHGUpdateInBaseHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	hostGroup(params["host"], params["hgName"])
	err := json.NewEncoder(w).Encode("ok")
	if err != nil {
		log.Fatalf("Error on getting HG list: %s", err)
	}
}
func getHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	data := getHGList(params["host"])
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting HG list: %s", err)
	}
}

func getAllHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := getHGAllList()
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting All HG list: %s", err)
	}
}

func getHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["swe_id"])
	data := getHG(id)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalf("Error on getting SWE list: %s", err)
	}
}

func getAllHostsHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clientIp := r.Header.Get("X-Forwarded-For")
	logger.Info.Printf("%s got /hosts", clientIp)

	err := json.NewEncoder(w).Encode(globConf.Hosts)
	if err != nil {
		log.Fatalf("Error on getting SWE list: %s", err)
	}
}

// ===============================
// POST
// ===============================
func postHGCheckHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t HGPost
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}
	data := postCheckHG(t.TargetHost, t.SourceHgId)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting SWE list: %s", err)
		return
	}
}

func postHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t HGPost
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	data, err := postHG(t.SourceHost, t.TargetHost, t.SourceHgId)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	jDataBase, _ := json.Marshal(POSTStructBase{data.BaseInfo})

	fmt.Println(string(jDataBase))

	response, err := ForemanAPI("POST", t.TargetHost, "hostgroups", string(jDataBase))
	if err == nil {
		if len(data.Overrides) > 0 {
			for _, ovr := range data.Overrides {

				p := struct {
					Match string `json:"match"`
					Value string `json:"value"`
				}{Match: ovr.Match, Value: ovr.Value}

				d := POSTStructOvrVal{p}
				jDataOvr, _ := json.Marshal(d)
				uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ForemanId)

				fmt.Println(string(jDataOvr))

				resp, err := ForemanAPI("POST", t.TargetHost, uri, string(jDataOvr))

				fmt.Println(string(resp))

				if err != nil {
					err = json.NewEncoder(w).Encode(string(resp))
					if err != nil {
						logger.Error.Printf("Error on POST HG: %s", err)
						return
					}
				}
			}
		}

		// Commit new HG for target host
		hostGroup(t.TargetHost, data.BaseInfo.Name)

		err = json.NewEncoder(w).Encode(string(response))
		if err != nil {
			logger.Error.Printf("Error on POST HG: %s", err)
			return
		}
	}
}

func postHGUpdateHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	fmt.Println(r)
	fmt.Println(r.Body)

	var t HGPost
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG JSON Decode: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deleteHG(t.TargetHost, t.TargetHgId)

	data, err := postHG(t.SourceHost, t.TargetHost, t.SourceHgId)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jDataBase, _ := json.Marshal(POSTStructBase{data.BaseInfo})

	response, err := ForemanAPI("POST", t.TargetHost, "hostgroups", string(jDataBase))
	if err == nil {
		if len(data.Overrides) > 0 {
			for _, ovr := range data.Overrides {

				p := struct {
					Match string `json:"match"`
					Value string `json:"value"`
				}{Match: ovr.Match, Value: ovr.Value}

				d := POSTStructOvrVal{p}
				jDataOvr, _ := json.Marshal(d)
				uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ForemanId)

				resp, err := ForemanAPI("POST", t.TargetHost, uri, string(jDataOvr))

				if err != nil {
					err = json.NewEncoder(w).Encode(string(resp))
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						logger.Error.Printf("Error on POST HG: %s", err)
						return
					}
				}
			}
		}

		// Commit new HG for target host
		hostGroup(t.TargetHost, data.BaseInfo.Name)

		err = json.NewEncoder(w).Encode(string(response))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Printf("Error on POST HG: %s", err)
			return
		}
	}
}
