package API

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"runtime"
	"sort"
	"sync"
	"time"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Get Smart Classes from Foreman
func (Get) All(host string, ctx *user.GlobalCTX) ([]Parameter, error) {

	// VARS
	var (
		r      Parameters
		ids    []int
		result Result
	)

	// =======
	// Get From Foreman ============================================
	uri := fmt.Sprintf("smart_class_parameters?per_page=%d", ctx.Config.Api.GetPerPage)
	response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		utils.Error.Printf("%q:\n %q\n", err, response)
	}

	// SC PAGER ============================================
	if r.Total > ctx.Config.Api.GetPerPage {
		pagesRange := utils.Pager(r.Total, ctx.Config.Api.GetPerPage)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("smart_class_parameters?format=json&page=%d&per_page=%d", i, ctx.Config.Api.GetPerPage)
			response, err := utils.ForemanAPI("GET", host, uri, "", ctx)
			if err == nil {
				err := json.Unmarshal(response.Body, &r)
				if err != nil {
					return result.parameters, err
				}

				for _, i := range r.Results {
					ids = append(ids, i.ForemanID)
				}
			}
		}
	} else {
		for _, i := range r.Results {
			ids = append(ids, i.ForemanID)
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
	splitIDs := utils.SplitToQueue(ids, runtime.NumCPU())
	for num, ids := range splitIDs {
		wg.Add(1)
		t := time.Now()
		fmt.Printf("Worker %d started\tjobs: %d\t %q\n", num, len(ids), t)
		go func(IDs []int, w int, s time.Time) {
			wq <- func() {
				defer wg.Done()
				for _, ID := range IDs {
					var r Parameter
					uri := fmt.Sprintf("smart_class_parameters/%d", ID)
					response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
					if response.StatusCode != 200 {
						fmt.Println("SC Parameters, ID:", ID, response.StatusCode, host)
					}
					err := json.Unmarshal(response.Body, &r)
					if err != nil {
						utils.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
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

	return result.parameters, nil
}

// Get Smart Class Overrides by Foreman ID
func (Get) OverridesByID(host string, ForemanID int, ctx *user.GlobalCTX) ([]OverrideValue, error) {

	var ovr OverrideValues
	var result []OverrideValue

	uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ForemanID)
	response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &ovr)
	if err != nil {
		utils.Error.Printf("%q:\n %q\n", err, response)
		return nil, err
	}

	for _, k := range ovr.Results {
		result = append(result, k)
	}

	return result, nil
}

// Get Smart Class by Foreman ID
func (Get) ByID(host string, foremanId int, ctx *user.GlobalCTX) (Parameter, error) {

	var result Parameter

	uri := fmt.Sprintf("smart_class_parameters/%d", foremanId)
	response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &result)
	if err != nil {
		utils.Error.Printf("%q:\n %q\n", err, response)
		return Parameter{}, err
	}

	return result, nil
}

// Get Smart Classes by Puppet Class ID
func (Get) ByPuppetClassID(host string, pcId int, ctx *user.GlobalCTX) []Parameter {
	var result Parameters

	uri := fmt.Sprintf("puppetclasses/%d/smart_class_parameters", pcId)
	response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &result)
	if err != nil {
		utils.Error.Printf("%q:\n %q\n", err, response)
	}
	return result.Results
}
