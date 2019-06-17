package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
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
	SmartClasses []models.SCParameter
}

// Get Smart Classes from Foreman
// Result struct
type SCResult struct {
	sync.Mutex
	resSlice []models.SCParameter
}

func (r *SCResult) Add(ID models.SCParameter) {
	r.Lock()
	r.resSlice = append(r.resSlice, ID)
	r.Unlock()
}
func GetAll(host string, cfg *models.Config) ([]models.SCParameter, error) {
	var r models.SCParameters
	var ids []int
	var result SCResult

	// Get From Foreman ============================================
	uri := fmt.Sprintf("smart_class_parameters?per_page=%d", cfg.Api.GetPerPage)
	response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}

	// SC PAGER ============================================
	if r.Total > cfg.Api.GetPerPage {
		pagesRange := utils.Pager(r.Total, cfg.Api.GetPerPage)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("smart_class_parameters?format=json&page=%d&per_page=%d", i, cfg.Api.GetPerPage)
			response, err := logger.ForemanAPI("GET", host, uri, "", cfg)
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

	fmt.Println(len(ids))

	splitIDs := utils.SplitToQueue(ids, runtime.NumCPU())
	for num, ids := range splitIDs {
		wg.Add(1)
		t := time.Now()
		fmt.Printf("Worker %d started\tjobs: %d\t %q\n", num, len(ids), t)
		go func(IDs []int, w int, s time.Time) {
			wq <- func() {
				defer wg.Done()
				for _, ID := range IDs {
					var r models.SCParameter
					uri := fmt.Sprintf("smart_class_parameters/%d", ID)
					response, _ := utils.ForemanAPI("GET", host, uri, "", cfg)
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

	fmt.Println(len(result.resSlice))

	return result.resSlice, nil
}

// Get Smart Classes Overrides from Foreman
func SCOverridesById(host string, ForemanID int, cfg *models.Config) []models.OverrideValue {
	var r models.OverrideValues
	var result []models.OverrideValue

	uri := fmt.Sprintf("smart_class_parameters/%d/override_values?per_page=%d", ForemanID, cfg.Api.GetPerPage)
	response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	if r.Total > cfg.Api.GetPerPage {
		pagesRange := utils.Pager(r.Total, cfg.Api.GetPerPage)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values?page=%d&per_page=%d", ForemanID, i, cfg.Api.GetPerPage)
			response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
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

func SCByPCJson(host string, pcId int, cfg *models.Config) []models.SCParameter {
	var r models.SCParameters

	uri := fmt.Sprintf("puppetclasses/%d/smart_class_parameters", pcId)
	response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r.Results
}

// ===
func SCByPCJsonV2(host string, pcId int, cfg *models.Config) models.PCSCParameters {
	var r models.PCSCParameters
	uri := fmt.Sprintf("puppetclasses/%d", pcId)
	response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r
}

func SCByFId(host string, foremanId int, cfg *models.Config) models.SCParameter {
	var r models.SCParameter

	uri := fmt.Sprintf("smart_class_parameters/%d", foremanId)
	response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	return r
}
