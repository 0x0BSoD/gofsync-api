package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"sort"
)

// ===============
// INSERT
// ===============

type SmartClasses struct {
	SmartClasses []models.SCParameter
}

// Get Smart Classes from Foreman
func GetAll(host string, cfg *models.Config) ([]models.SCParameter, error) {
	var r models.SCParameters
	var resultId []int
	var result SmartClasses

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
					return result.SmartClasses, err
				}

				for _, i := range r.Results {
					resultId = append(resultId, i.ID)
				}
			}
		}
	} else {
		for _, i := range r.Results {
			resultId = append(resultId, i.ID)
		}
	}
	// SC PAGER ============================================

	// Getting Data from foreman ===============================
	sort.Ints(resultId)
	//WORKERS := runtime.NumCPU()
	//var wgWorkers sync.WaitGroup
	//var wgPool sync.WaitGroup
	//wgPool.Add(len(resultId))
	//wgWorkers.Add(WORKERS)
	//var addResultMutex sync.Mutex
	//taskChan := make(chan int)
	//resChan := make(chan models.SCParameter, len(resultId))

	//Spin up workers ===
	//for i := 0; i < WORKERS; i++ {
	//	go asyncWorker(i, taskChan, resChan, host, &wgWorkers, cfg)
	//}

	// Send tasks to him ===
	for _, i := range resultId {

		// sync ====
		r := worker(i, host, cfg)
		result.SmartClasses = append(result.SmartClasses, r)

		// async ====
		//go func(_id int, r *SmartClasses, wg *sync.WaitGroup) {
		//	defer wg.Done()
		//	taskChan <- _id
		//addResult(<-resChan, r, &addResultMutex)
		//}(i, &result, &wgPool)
	}
	//wgPool.Wait()
	//close(taskChan)

	//for i := range resChan {
	//	result.SmartClasses = append(result.SmartClasses, i)
	//}

	// Sort by ID ===
	//sort.Slice(result.SmartClasses, func(i, j int) bool {
	//	return result.SmartClasses[i].ID < result.SmartClasses[j].ID
	//})

	return result.SmartClasses, nil
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
