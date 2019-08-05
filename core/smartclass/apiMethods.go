package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"runtime"
	"sort"
	"sync"
	"time"
)

// ===============
// INSERT
// ===============

type SmartClasses struct {
	SmartClasses []SCParameter
}

// Get Smart Classes from Foreman
// Result struct
type SCResult struct {
	sync.Mutex
	resSlice []SCParameter
}

func (r *SCResult) Add(ID SCParameter) {
	r.Lock()
	r.resSlice = append(r.resSlice, ID)
	r.Unlock()
}
func GetAll(host string, ctx *user.GlobalCTX) ([]SCParameter, error) {
	var r SCParameters
	var ids []int
	var result SCResult

	// Get From Foreman ============================================
	uri := fmt.Sprintf("smart_class_parameters?per_page=%d", ctx.Config.Api.GetPerPage)
	response, _ := logger.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}

	// SC PAGER ============================================
	if r.Total > ctx.Config.Api.GetPerPage {
		pagesRange := utils.Pager(r.Total, ctx.Config.Api.GetPerPage)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("smart_class_parameters?format=json&page=%d&per_page=%d", i, ctx.Config.Api.GetPerPage)
			response, err := logger.ForemanAPI("GET", host, uri, "", ctx)
			if err == nil {
				err := json.Unmarshal(response.Body, &r)
				if err != nil {
					return result.resSlice, err
				}

				for _, i := range r.Results {
					ids = append(ids, i.ID)
				}
			}
		}
	} else {
		for _, i := range r.Results {
			ids = append(ids, i.ID)
		}
	}

	// Getting Additional data from foreman ===============================
	sort.Ints(ids)

	// ver 2 ===
	// Create a new WorkQueue.
	wq := utils.New()
	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	var wg sync.WaitGroup

	//fmt.Println(len(ids))

	splitIDs := utils.SplitToQueue(ids, runtime.NumCPU())
	for num, ids := range splitIDs {
		wg.Add(1)
		t := time.Now()
		fmt.Printf("Worker %d started\tjobs: %d\t %q\n", num, len(ids), t)
		go func(IDs []int, w int, s time.Time) {
			wq <- func() {
				defer wg.Done()
				for _, ID := range IDs {
					var r SCParameter
					uri := fmt.Sprintf("smart_class_parameters/%d", ID)
					response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
					if response.StatusCode != 200 {
						fmt.Println("SC Parameters, ID:", ID, response.StatusCode, host)
					}
					err := json.Unmarshal(response.Body, &r)
					if err != nil {
						logger.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
					}

					fmt.Println("Worker: ", w, "\t Parameter: ", r.Parameter)

					result.Add(r)
				}
				t := time.Since(s)
				fmt.Printf("Worker %d done\t%q\n", w, t)
			}
		}(ids, num, t)
	}

	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	//fmt.Println(len(result.resSlice))

	return result.resSlice, nil
}

// Get Smart Classes Overrides from Foreman
func SCOverridesById(host string, ForemanID int, ctx *user.GlobalCTX) []OverrideValue {
	var r OverrideValues
	var result []OverrideValue

	uri := fmt.Sprintf("smart_class_parameters/%d/override_values?per_page=%d", ForemanID, ctx.Config.Api.GetPerPage)
	response, _ := logger.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	if r.Total > ctx.Config.Api.GetPerPage {
		pagesRange := utils.Pager(r.Total, ctx.Config.Api.GetPerPage)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values?page=%d&per_page=%d", ForemanID, i, ctx.Config.Api.GetPerPage)
			response, _ := logger.ForemanAPI("GET", host, uri, "", ctx)
			err := json.Unmarshal(response.Body, &r)
			if err != nil {
				logger.Error.Printf("%q:\n %q\n", err, response)
			}
			for _, j := range r.Results {
				result = append(result, j)
			}
		}
	} else {
		for _, k := range r.Results {
			result = append(result, k)
		}
	}
	return result
}

func SCByPCJson(host string, pcId int, ctx *user.GlobalCTX) []SCParameter {
	var r SCParameters

	uri := fmt.Sprintf("puppetclasses/%d/smart_class_parameters", pcId)
	response, _ := logger.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r.Results
}

// ===
func SCByPCJsonV2(host string, pcId int, ctx *user.GlobalCTX) PCSCParameters {
	var r PCSCParameters
	uri := fmt.Sprintf("puppetclasses/%d", pcId)
	response, _ := logger.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r
}

func SCByFId(host string, foremanId int, ctx *user.GlobalCTX) SCParameter {
	var r SCParameter

	uri := fmt.Sprintf("smart_class_parameters/%d", foremanId)
	response, _ := logger.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r
}
