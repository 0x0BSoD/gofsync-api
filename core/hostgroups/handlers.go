package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ===============================
// GET
// ===============================
func GetForemanID(ctx *user.GlobalCTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx.Set(&user.Claims{Username: "srv_foreman"}, "fake")

		params := mux.Vars(r)
		data := ForemanID(ctx.Config.Hosts[params["host"]], params["hgName"], ctx)

		utils.SendResponse(w, "error on getting foremanId for env: %s", data)
	}
}

// Get HG info from Foreman
func GetHGFHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data, err := HostGroupJson(params["host"], params["hgName"], ctx)

	if (HgError{}) != err {
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			utils.Error.Printf("Error on getting HG: %s", err)
			return
		}
	} else {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			utils.Error.Printf("Error on getting HG: %s", err)
		}
	}
}

func GetHGCheckNameHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := HostGroupCheckName(params["host"], params["hgName"], ctx)

	if data.Error == "error -1" {
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte("410 - Foreman server gone"))
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG check: %s", err)
	}
}

func GetHGCheckIDHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := HostGroupCheckID(params["host"], params["id"], ctx)

	if data.Error == "error -1" {
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte("410 - Foreman server gone"))
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG check: %s", err)
	}
}

func GetHGCheckUAHttp(ctx *user.GlobalCTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx.Set(&user.Claims{Username: "srv_foreman"}, "fake")

		params := mux.Vars(r)
		data := HostGroupCheckName(params["host"], params["hgName"], ctx)

		if data.Error == "error -1" {
			w.WriteHeader(http.StatusGone)
			_, _ = w.Write([]byte("410 - Foreman server gone"))
			return
		}
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			utils.Error.Printf("Error on getting HG check: %s", err)
		}
	}
}

func GetHGUpdateInBaseHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	ID, err := HostGroup(params["host"], params["hgName"], ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	data := Get(ID, ctx)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on updating HG: %s", err)
	}
}

func GetHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := OnHost(ctx.Config.Hosts[params["host"]], ctx)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetAllHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	data := All(ctx)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		utils.Error.Printf("Error on getting all HG list: %s", err)
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
		utils.Error.Printf("Error on getting HG: %s", err)
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
		_ = json.NewEncoder(w).Encode(err)
		utils.Error.Printf("Error on getting hosts: %s", err)
		return
	}

	err = json.NewEncoder(w).Encode("ok")
	if err != nil {
		utils.Error.Printf("Error on getting hosts: %s", err)
	}
}

func CompareHG(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	ts := time.Now()
	fmt.Println("Compare started")

	ID := ID(ctx.Config.Hosts[params["host"]], params["host"], ctx)
	dbHG := Get(ID, ctx)
	foremanHG, _ := HostGroupJson(params["host"], params["hgName"], ctx)

	cr := CompareHGWorker(dbHG, foremanHG)
	fmt.Println("Compare done, ", time.Since(ts))

	err := json.NewEncoder(w).Encode(cr)
	if err != nil {
		utils.Error.Printf("Error on getting hosts: %s", err)
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
		utils.Error.Printf("Error on POST HG: %s", err)
		return
	}
	data := PostCheckHG(t.TargetHost, t.SourceHgId, ctx)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting SWE list: %s", err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	hostID := ctx.Config.Hosts[params["host"]]

	// Decode HostGroup
	decoder := json.NewDecoder(r.Body)
	var hostGroupJSON HGElem
	err := decoder.Decode(&hostGroupJSON)
	if err != nil {
		utils.Error.Printf("Error on POST HG: %s", err)
		return
	}

	envID := environment.ForemanID(hostID, hostGroupJSON.Environment, ctx)
	locationsIDs := locations.DbAllForemanID(hostID, ctx)
	pID := ForemanID(hostID, "SWE", ctx)

	if envID != -1 {
		existId := ForemanID(hostID, hostGroupJSON.Name, ctx)
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
				utils.Error.Printf("Error on POST HG: %s", err)
				_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
				return
			}
			// Send response to client
			_ = json.NewEncoder(w).Encode(resp)
		} else {
			resp, err := UpdateHG(toSubmit, params["host"], ctx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				utils.Error.Printf("Error on POST HG: %s", err)
				_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
				return
			}
			// Send response to client
			_ = json.NewEncoder(w).Encode(resp)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Error.Printf("Error on Create HG: %s", err)
		_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
	}

}

func Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	hostID := ctx.Config.Hosts[params["host"]]

	// Decode HostGroup
	decoder := json.NewDecoder(r.Body)
	var hostGroupJSON HGElem
	err := decoder.Decode(&hostGroupJSON)
	if err != nil {
		utils.Error.Printf("Error on POST HG: %s", err)
		return
	}

	envID := environment.ForemanID(hostID, hostGroupJSON.Environment, ctx)
	locationsIDs := locations.DbAllForemanID(hostID, ctx)
	pID := ForemanID(hostID, "SWE", ctx)

	if envID != -1 {
		existId := hostGroupJSON.ForemanID
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

		resp, err := UpdateHG(toSubmit, params["host"], ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Error.Printf("Error on POST HG: %s", err)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
			return
		}
		// Send response to client
		_ = json.NewEncoder(w).Encode(resp)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Error.Printf("Error on Create HG: %s", err)
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
		utils.Error.Printf("Error on POST HG: %s", err)
		return
	}

	// Get data from DB ====================================================
	data, err := HGDataItem(t.SourceHost, t.TargetHost, t.SourceHgId, ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(fmt.Sprintf("Foreman Api Error: %q", err))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Error.Printf("Error on POST HG: %s", err)
		}
		return
	}

	scIDs := ResolveSmartClasses(ctx.Config.Hosts[t.TargetHost], data.BaseInfo.PuppetClassIds, ctx)
	for _, i := range scIDs {
		err := setOverridable(t.TargetHost, i, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Error.Printf("Error on PUT HG: %s", err)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on PUT HG: %s", err))
		}
	}

	// Submit host group ====================================================
	if data.ExistId == -1 {
		resp, err := PushNewHG(data, t.TargetHost, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Error.Printf("Error on POST HG: %s", err)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
			return
		}
		// Send response to client
		_ = json.NewEncoder(w).Encode(resp)
	} else {
		resp, err := UpdateHG(data, t.TargetHost, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Error.Printf("Error on PUT HG: %s", err)
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

	var postBody struct {
		Batch        map[string]map[string]BatchPostStruct `json:"batch"`
		UpdateSource bool                                  `json:"update_source"`
	}

	err := decoder.Decode(&postBody)
	if err != nil {
		utils.Error.Printf("Error on POST HG: %s", err)
		return
	}

	uniqHGS := make([]string, 0, len(postBody.Batch))
	var sourceHost string
	var targetHost string
	for _, HGs := range postBody.Batch {
		tmpHGS := make([]string, 0, len(HGs))
		for _, hg := range HGs {
			if !utils.StringInSlice(hg.HGName, uniqHGS) {
				tmpHGS = append(tmpHGS, hg.HGName)
				sourceHost = hg.SHost
				targetHost = hg.THost
			}
		}
		uniqHGS = append(uniqHGS, tmpHGS...)
	}

	if postBody.UpdateSource {

		// Create a new WorkQueue.
		wq := utils.New()
		// This sync.WaitGroup is to make sure we wait until all of our work
		// is done.
		var wg sync.WaitGroup
		idx := 0

		for _, HG := range uniqHGS {
			var lock sync.Mutex
			wg.Add(1)
			go func(HG string) {
				wq <- func() {
					defer func() {
						wg.Done()
					}()

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						HostName:  targetHost,
						Resource:  models.HostGroup,
						Operation: "submit",
						UserName:  ctx.Session.UserName,
						AdditionalData: models.CommonOperation{
							Message:   "Updating HG",
							HostGroup: HG,
						},
					})
					// ---

					ID := ID(ctx.Config.Hosts[sourceHost], HG, ctx)
					dbHG := Get(ID, ctx)
					foremanHG, _ := HostGroupJson(sourceHost, HG, ctx)

					cr := CompareHGWorker(dbHG, foremanHG)

					if !cr {
						_, err := HostGroup(sourceHost, HG, ctx)
						if err != nil {
							utils.Error.Printf("Error on POST HG: %s", err)
							return
						}
					} else {
						fmt.Println("Update not needed")
					}

					lock.Lock()
					idx++
					lock.Unlock()
					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						HostName:  targetHost,
						Resource:  models.HostGroup,
						Operation: "submit",
						UserName:  ctx.Session.UserName,
						AdditionalData: models.CommonOperation{
							Message:   "Updating HG",
							Done:      true,
							HostGroup: HG,
						},
					})
					// ---
				}
			}(HG)
		}
		// Wait for all of the work to finish, then close the WorkQueue.
		wg.Wait()
		close(wq)

	}

	// =================================================================================================================
	// Create a new WorkQueue.
	wq := utils.New()
	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	num := 0
	var wgSecond sync.WaitGroup
	for _, HGs := range postBody.Batch {
		wgSecond.Add(1)
		startTime := time.Now()
		fmt.Printf("Worker %d started\tjobs: %d\t %q\n", num, len(HGs), startTime)
		var lock sync.Mutex

		go func(HGs map[string]BatchPostStruct, wID int, st time.Time) {
			wq <- func() {
				defer func() {
					fmt.Printf("Worker %d done\t %q\n", wID, startTime)
					wgSecond.Done()
				}()
				for _, hg := range HGs {
					lock.Lock()
					if hg.Environment.TargetID != -1 {

						// Socket Broadcast ---
						ctx.Session.SendMsg(models.WSMessage{
							Broadcast: false,
							HostName:  hg.THost,
							Resource:  models.HostGroup,
							Operation: "submit",
							UserName:  ctx.Session.UserName,
							AdditionalData: models.CommonOperation{
								Message:   "Uploading HG",
								HostGroup: hg.HGName,
							},
						})
						// ---

						hg.InProgress = true
						hg.Done = false

						// Get data from DB ====================================================
						data, err := HGDataItem(hg.SHost, hg.THost, hg.Foreman.SourceID, ctx)
						if err != nil {
							utils.Error.Println(err)
							return
						}

						scIDs := ResolveSmartClasses(ctx.Config.Hosts[hg.THost], data.BaseInfo.PuppetClassIds, ctx)
						for _, i := range scIDs {
							err := setOverridable(hg.THost, i, ctx)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								utils.Error.Printf("Error on PUT HG: %s", err)
								_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on PUT HG: %s", err))
							}
						}

						var resp string
						// Submit host group ====================================================
						if data.ExistId == -1 {
							resp, err = PushNewHG(data, hg.THost, ctx)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								utils.Error.Printf("Error on POST HG: %s", err)
								_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
								return
							}
						} else {
							resp, err = UpdateHG(data, hg.THost, ctx)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								utils.Error.Printf("Error on PUT HG: %s", err)
								_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on PUT HG: %s", err))
								return
							}
						}

						hg.Done = true
						hg.InProgress = false
						hg.HTTPResp = resp

					}
					lock.Unlock()

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						HostName:  hg.THost,
						Resource:  models.HostGroup,
						Operation: "submit",
						UserName:  ctx.Session.UserName,
						AdditionalData: models.CommonOperation{
							Message:   "Uploading HG done",
							Done:      true,
							HostGroup: hg.HGName,
						},
					})
					// ---
				}
			}
			lock.Lock()
			num++
			lock.Unlock()
		}(HGs, num, startTime)
	}
	// Wait for all of the work to finish, then close the WorkQueue.
	wgSecond.Wait()
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

	var t struct {
		Name   string                          `json:"name"`
		Source string                          `json:"source"`
		Target string                          `json:"target"`
		Data   []smartclass.OverrideParameters `json:"data"`
	}
	err := decoder.Decode(&t)
	if err != nil {
		utils.Error.Printf("Error on POST HG: %s", err)
		return
	}

	// Get data from DB ====================================================
	ExistId := locations.ID(ctx.Config.Hosts[t.Target], t.Name, ctx)

	// Submit new location ====================================================
	if ExistId == -1 {

		hgIDs := ForemanIDs(ctx.Config.Hosts[t.Target], ctx)
		envIDs := environment.ForemanIDs(ctx.Config.Hosts[t.Target], ctx)
		spID := environment.ApiGetSmartProxy(t.Target, ctx)

		type newLoc struct {
			Location struct {
				Name           string `json:"name"`
				HostgroupIds   []int  `json:"hostgroup_ids"`
				EnvironmentIds []int  `json:"environment_ids"`
				SmartProxyIds  []int  `json:"smart_proxy_ids"`
			} `json:"location"`
		}

		_json, _ := json.Marshal(newLoc{
			Location: struct {
				Name           string `json:"name"`
				HostgroupIds   []int  `json:"hostgroup_ids"`
				EnvironmentIds []int  `json:"environment_ids"`
				SmartProxyIds  []int  `json:"smart_proxy_ids"`
			}{Name: t.Name, HostgroupIds: hgIDs, EnvironmentIds: envIDs, SmartProxyIds: []int{spID}},
		})

		resp, err := utils.ForemanAPI("POST", t.Target, "locations", string(_json), ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err)
		}
		utils.Info.Println(string(resp.Body), resp.RequestUri)

		for _, i := range t.Data {
			for _, ovr := range i.Parameters {
				p := param{
					Value: ovr.Value,
					Match: fmt.Sprintf("location=%s", t.Name),
				}
				var ScForemanId int
				if t.Source != t.Target {
					targetSC := smartclass.GetSC(ctx.Config.Hosts[t.Target], i.PuppetClass, ovr.Name, ctx)
					ScForemanId = targetSC.ForemanID
				} else {
					ScForemanId = ovr.ParameterForemanId
				}
				_json, _ := json.Marshal(p)

				uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ScForemanId)
				resp, err := utils.ForemanAPI("POST", t.Target, uri, string(_json), ctx)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(err)
				}
				utils.Info.Println(string(resp.Body), resp.RequestUri)

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

		resp, err := utils.ForemanAPI("POST", t.Target, "locations", string(_json), ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err)
		}
		utils.Info.Println(string(resp.Body), resp.RequestUri)

		for _, i := range t.Data {
			for _, ovr := range i.Parameters {
				p := param{
					Value: ovr.Value,
					Match: fmt.Sprintf("location=%s", t.Name),
				}
				var ScForemanId int
				if t.Source != t.Target {
					targetSC := smartclass.GetSC(ctx.Config.Hosts[t.Target], i.PuppetClass, ovr.Name, ctx)
					ScForemanId = targetSC.ForemanID
				} else {
					ScForemanId = ovr.ParameterForemanId
				}
				_json, _ := json.Marshal(p)

				uri := fmt.Sprintf("smart_class_parameters/%d/override_values/%d", ScForemanId, ovr.OverrideForemanId)
				resp, err := utils.ForemanAPI("PUT", t.Target, uri, string(_json), ctx)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(err)
				}
				utils.Info.Println(string(resp.Body), resp.RequestUri)

			}
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
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}
