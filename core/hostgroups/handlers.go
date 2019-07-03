package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/middleware"
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
		w.Write([]byte("410 - Foreman server gone"))
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
	data := GetHG(ID, ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on updating HG: %s", err)
	}
}

func GetHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := GetHGList(params["host"], ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetAllHGListHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	data := GetHGAllList(ctx)
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
	data := GetHG(id, ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG: %s", err)
	}
}

func GetAllHostsHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	data := AllHosts(ctx)
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
	decoder := json.NewDecoder(r.Body)
	var t HGElem
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST HG: %s", err)
		return
	}

	envID := environment.DbForemanID(params["host"], t.Environment, ctx)
	locationsIDs := locations.DbAllForemanID(params["host"], ctx)
	pID, _ := strconv.Atoi(t.ParentId)

	if envID != -1 {
		existId := CheckHGID(t.Name, params["host"], ctx)
		data, _ := HGDataNewItem(params["host"], t, ctx)
		base := HWPostRes{
			BaseInfo: HostGroupBase{
				Name:           t.Name,
				EnvironmentId:  envID,
				LocationIds:    locationsIDs,
				ParentId:       pID,
				PuppetClassIds: data.BaseInfo.PuppetClassIds,
			},
			ExistId:    existId,
			Overrides:  data.Overrides,
			Parameters: data.Parameters,
		}

		fmt.Println(existId)
		fmt.Println(base.ExistId)

		if base.ExistId == -1 {
			resp, err := PushNewHG(base, params["host"], ctx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error.Printf("Error on POST HG: %s", err)
				_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error on POST HG: %s", err))
				return
			}
			// Send response to client
			_ = json.NewEncoder(w).Encode(resp)
		} else {
			resp, err := UpdateHG(base, params["host"], ctx)
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

						ctx.Session.WSMessage <- msg

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

						ctx.Session.WSMessage <- msg
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
