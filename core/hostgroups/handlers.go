package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ===============================
// GET
// ===============================

// Get HG info from Foreman
func GetHGFHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data, err := HostGroupJson(params["host"], params["hgName"], ctx)
	if (HgError{}) != err {
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
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := HostGroupCheck(params["host"], params["hgName"], ctx)
	if data.Error == "error -1" {
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte("410 - Foreman server gone"))
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG check: %s", err)
	}
}

func GetHGUpdateInBaseHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	ID := HostGroup(params["host"], params["hgName"], ctx)
	data := Get(ID, ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on updating HG: %s", err)
	}
}

func GetHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := OnHost(params["host"], ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetAllHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	data := All(ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting all HG list: %s", err)
	}
}

func GetHGHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["swe_id"])
	data := Get(id, ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG: %s", err)
	}
}

func GetAllHostsHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	data := PuppetHosts(ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting hosts: %s", err)
	}
}

func CommitGitHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["swe_id"])
	err := CommitJsonByHgID(id, params["host"], ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("ok")
		logger.Error.Printf("Error on getting hosts: %s", err)
	}

	err = json.NewEncoder(w).Encode("ok")
	if err != nil {
		logger.Error.Printf("Error on getting hosts: %s", err)
	}
}

// ===============================
// POST
// ===============================
func PostHGCheckHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var t HGPost
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
	}
	data := PostCheckHG(t.TargetHost, t.SourceHgId, ctx)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting SWE list: %s", err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	// Decode HostGroup
	decoder := json.NewDecoder(r.Body)
	var hostGroupJSON HGElem
	err := decoder.Decode(&hostGroupJSON)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	envID := environment.DbForemanID(params["host"], hostGroupJSON.Environment, ctx)
	locationsIDs := locations.DbAllForemanID(params["host"], ctx)
	pID, _ := strconv.Atoi(hostGroupJSON.ParentId)

	if envID != -1 {
		existId := FID(hostGroupJSON.Name, params["host"], ctx)
		NewHostGroup, _ := HGDataNewItem(params["host"], hostGroupJSON, ctx)

		// Brand new crafted host group
		toSubmit := HWPostRes{
			BaseInfo: HostGroupBase{
				Name:           hostGroupJSON.Name,
				EnvironmentId:  envID,
				LocationIds:    locationsIDs,
				ParentId:       pID,
				PuppetClassIds: NewHostGroup.BaseInfo.PuppetClassIds,
			},
			ExistId:    existId,
			Overrides:  NewHostGroup.Overrides,
			Parameters: NewHostGroup.Parameters,
		}

		// Check Environment

		if toSubmit.ExistId == -1 {
			resp, err := PushNewHG(toSubmit, params["host"], ctx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error.Printf("Error on POST HG: %s", err)
				_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
				return
			}
			// Send response to client
			_ = json.NewEncoder(w).Encode(resp)
		} else {
			resp, err := UpdateHG(toSubmit, params["host"], ctx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error.Printf("Error on POST HG: %s", err)
				_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
				return
			}
			// Send response to client
			_ = json.NewEncoder(w).Encode(resp)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Printf("Error on Create HG: %s", err)
		_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
	}

}

func Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var t HGPost
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	// Get data from DB ====================================================
	data, err := HGDataItem(t.SourceHost, t.TargetHost, t.SourceHgId, ctx)
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
		resp, err := PushNewHG(data, t.TargetHost, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Printf("Error on POST HG: %s", err)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
		}
		// Send response to client
		_ = json.NewEncoder(w).Encode(resp)
	} else {
		resp, err := UpdateHG(data, t.TargetHost, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Printf("Error on PUT HG: %s", err)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on PUT HG: %s", err))
		}
		// Send response to client
		_ = json.NewEncoder(w).Encode(resp)
	}

}

func BatchPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var postBody map[string][]BatchPostStruct
	err := decoder.Decode(&postBody)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	var uniqHGS []string
	var sourceHost string
	for _, HGs := range postBody {
		for _, hg := range HGs {
			if !utils.StringInSlice(hg.HGName, uniqHGS) {
				uniqHGS = append(uniqHGS, hg.HGName)
				sourceHost = hg.SHost
			}
		}
	}

	idx := 1
	for _, HG := range uniqHGS {
		// Socket Broadcast ---
		data := models.Step{
			Actions: "Updating Source HostGroups",
			State:   fmt.Sprintf("HostGroup: %s %d/%d", HG, idx, len(uniqHGS)),
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
		ID := HostGroup(sourceHost, HG, ctx)
		_ = Get(ID, ctx)
		idx++
	}

	// =================================================================================================================
	// Create a new WorkQueue.
	wq := utils.New()
	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	num := 0
	var wg sync.WaitGroup
	for _, HGs := range postBody {
		wg.Add(1)
		startTime := time.Now()
		fmt.Printf("Worker %d started\tjobs: %d\t %q\n", num, len(HGs), startTime)
		var lock sync.Mutex

		go func(HGs []BatchPostStruct, wID int, st time.Time) {
			wq <- func() {
				defer func() {
					fmt.Printf("Worker %d done\t %q\n", wID, startTime)
					wg.Done()
				}()
				for _, hg := range HGs {
					lock.Lock()
					if hg.Environment.TargetID != -1 {
						hg.InProgress = true
						hg.Done = false
						msg, _ := json.Marshal(hg)

						ctx.Session.SendMsg(msg)

						// Get data from DB ====================================================
						data, err := HGDataItem(hg.SHost, hg.THost, hg.Foreman.SourceID, ctx)
						if err != nil {
							logger.Error.Println(err)
							return
						}
						var resp string
						// Submit host group ====================================================
						if data.ExistId == -1 {
							resp, err = PushNewHG(data, hg.THost, ctx)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								logger.Error.Printf("Error on POST HG: %s", err)
								_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
							}
						} else {
							resp, err = UpdateHG(data, hg.THost, ctx)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								logger.Error.Printf("Error on PUT HG: %s", err)
								_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on PUT HG: %s", err))
							}
						}

						hg.Done = true
						hg.InProgress = false
						hg.HTTPResp = resp
						msg, _ = json.Marshal(hg)

						ctx.Session.SendMsg(msg)
					}
					lock.Unlock()
				}
			}
		}(HGs, num, startTime)

		lock.Lock()
		num++
		lock.Unlock()
	}
	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	_ = json.NewEncoder(w).Encode(postBody)
}

func SubmitLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)

	type param struct {
		Match string `json:"match"`
		Value string `json:"value"`
	}
	//d := POSTStructOvrVal{p}
	//jDataOvr, _ := json.Marshal(d)
	//uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
	//resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
	//if err != nil {
	//	return err
	//}
	//logger.Info.Println(string(resp.Body), resp.RequestUri)

	var t struct {
		Name   string                          `json:"name"`
		Source string                          `json:"source"`
		Target string                          `json:"target"`
		Data   []smartclass.OverrideParameters `json:"data"`
	}
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	// Get data from DB ====================================================
	ExistId := locations.DbID(t.Target, t.Name, ctx)

	// Submit host group ====================================================
	if ExistId == -1 {
		//POST /api/locations
		//{
		//	"location": {
		//	"name": "Test Location"
		//}
		//}
		type newLoc struct {
			Location struct {
				Name string `json:"name"`
			} `json:"location"`
		}

		_json, _ := json.Marshal(newLoc{
			Location: struct {
				Name string `json:"name"`
			}{Name: t.Name},
		})

		resp, err := logger.ForemanAPI("POST", t.Target, "locations", string(_json), ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err)
		}
		logger.Info.Println(string(resp.Body), resp.RequestUri)

		for _, i := range t.Data {
			//fmt.Println(i.PuppetClass)
			for _, ovr := range i.Parameters {
				//fmt.Println(ovr.Name)
				p := param{
					Value: ovr.Value,
					Match: fmt.Sprintf("location=%s", t.Name),
				}
				fmt.Println("Source FID:", ovr.ParameterForemanId)
				var ScForemanId int
				if t.Source != t.Target {
					targetSC := smartclass.GetSC(t.Target, i.PuppetClass, ovr.Name, ctx)
					ScForemanId = targetSC.ForemanID
					fmt.Println("Target FID:", targetSC.ForemanID)
				} else {
					fmt.Println("Target FID:", ovr.ParameterForemanId)
					ScForemanId = ovr.ParameterForemanId
				}
				_json, _ := json.Marshal(p)
				//fmt.Println(string(_json))
				//fmt.Println("----")

				uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ScForemanId)
				resp, err := logger.ForemanAPI("POST", t.Target, uri, string(_json), ctx)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(err)
				}
				logger.Info.Println(string(resp.Body), resp.RequestUri)

			}
			//fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		}
		// Send response to client
		_ = json.NewEncoder(w).Encode(t)
	} else {
		type newLoc struct {
			Location struct {
				Name string `json:"name"`
			} `json:"location"`
		}

		_json, _ := json.Marshal(newLoc{
			Location: struct {
				Name string `json:"name"`
			}{Name: t.Name},
		})

		resp, err := logger.ForemanAPI("POST", t.Target, "locations", string(_json), ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err)
		}
		logger.Info.Println(string(resp.Body), resp.RequestUri)

		for _, i := range t.Data {
			//fmt.Println(i.PuppetClass)
			for _, ovr := range i.Parameters {
				//fmt.Println(ovr.Name)
				p := param{
					Value: ovr.Value,
					Match: fmt.Sprintf("location=%s", t.Name),
				}
				fmt.Println("Source FID:", ovr.ParameterForemanId)
				var ScForemanId int
				if t.Source != t.Target {
					targetSC := smartclass.GetSC(t.Target, i.PuppetClass, ovr.Name, ctx)
					ScForemanId = targetSC.ForemanID
					fmt.Println("Target FID:", targetSC.ForemanID)
				} else {
					fmt.Println("Target FID:", ovr.ParameterForemanId)
					ScForemanId = ovr.ParameterForemanId
				}
				_json, _ := json.Marshal(p)
				//fmt.Println(string(_json))
				//fmt.Println("----")

				uri := fmt.Sprintf("smart_class_parameters/%d/override_values/%d", ScForemanId, ovr.OverrideForemanId)
				resp, err := logger.ForemanAPI("PUT", t.Target, uri, string(_json), ctx)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(err)
				}
				logger.Info.Println(string(resp.Body), resp.RequestUri)

			}
			//fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		}
		// Send response to client
		_ = json.NewEncoder(w).Encode(t)
	}
}

// ===============================
// PUT
// ===============================
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	Sync(params["host"], ctx)
	err := json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}
