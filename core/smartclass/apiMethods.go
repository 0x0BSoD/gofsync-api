package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"runtime"
	"sort"
	"strconv"
	"sync"
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
func GetAll(hostname string, ctx *user.GlobalCTX) ([]SCParameter, error) {
	var r SCParameters
	var result SCResult

	// Get From Foreman ============================================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: From foreman",
		Host:    hostname,
	}))
	uri := fmt.Sprintf("smart_class_parameters?format=json&per_page=9999999")
	response, _ := logger.ForemanAPI("GET", hostname, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}

	resCount := r.Total
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: Total classes: " + strconv.Itoa(resCount),
		Host:    hostname,
	}))
	ids := make([]int, 0, resCount)
	idsTmp := make([]int, 0, resCount)
	for _, i := range r.Results {
		idsTmp = append(idsTmp, i.ID)
	}
	sort.Ints(idsTmp)
	lv := -1
	for _, i := range idsTmp {
		if lv != i {
			ids = append(ids, i)
		}
		lv = i
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: Uniq classes: " + strconv.Itoa(len(ids)),
		Host:    hostname,
	}))

	// ver 2 ===
	// Create a new WorkQueue.
	wq := utils.New()
	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	var wg sync.WaitGroup

	splitIDs := utils.SplitToQueue(ids, runtime.NumCPU())
	for num, ids := range splitIDs {
		wg.Add(1)

		go func(IDs []int, w int) {
			wq <- func() {
				defer wg.Done()
				for _, ID := range IDs {
					var r SCParameter
					uri := fmt.Sprintf("smart_class_parameters/%d", ID)
					response, _ := utils.ForemanAPI("GET", hostname, uri, "", ctx)
					if response.StatusCode != 200 {
						fmt.Println("SC Parameters, ID:", ID, response.StatusCode, hostname)
					}
					err := json.Unmarshal(response.Body, &r)
					if err != nil {
						logger.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
					}
					result.Add(r)
				}
			}
		}(ids, num)
	}

	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	return result.resSlice, nil
}

// Get Smart Classes Overrides from Foreman
func SCOverridesById(hostname string, ForemanID int, ctx *user.GlobalCTX) []OverrideValue {
	var r OverrideValues
	var result []OverrideValue

	uri := fmt.Sprintf("smart_class_parameters/%d/override_values?per_page=%d", ForemanID, ctx.Config.Api.GetPerPage)
	response, _ := logger.ForemanAPI("GET", hostname, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	if r.Total > ctx.Config.Api.GetPerPage {
		pagesRange := utils.Pager(r.Total, ctx.Config.Api.GetPerPage)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values?page=%d&per_page=%d", ForemanID, i, ctx.Config.Api.GetPerPage)
			response, _ := logger.ForemanAPI("GET", hostname, uri, "", ctx)
			err := json.Unmarshal(response.Body, &r)
			if err != nil {
				logger.Error.Printf("%q:\n %q\n", err, response)
			}
			result = append(result, r.Results...)
		}
	} else {
		result = append(result, r.Results...)
	}
	return result
}

func SCByPCJson(hostname string, pcId int, ctx *user.GlobalCTX) []SCParameter {
	var r SCParameters

	uri := fmt.Sprintf("puppetclasses/%d/smart_class_parameters", pcId)
	response, _ := logger.ForemanAPI("GET", hostname, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r.Results
}

// ===
func SCByPCJsonV2(hostname string, pcId int, ctx *user.GlobalCTX) PCSCParameters {
	var r PCSCParameters

	uri := fmt.Sprintf("puppetclasses/%d", pcId)
	response, _ := logger.ForemanAPI("GET", hostname, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r
}

func SCByFId(hostname string, foremanId int, ctx *user.GlobalCTX) SCParameter {
	var r SCParameter

	uri := fmt.Sprintf("smart_class_parameters/%d", foremanId)
	response, _ := logger.ForemanAPI("GET", hostname, uri, "", ctx)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r
}
