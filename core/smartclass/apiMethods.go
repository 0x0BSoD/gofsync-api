package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"sort"
	"sync"
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
	//WORKERS := 6
	//var wgWorkers sync.WaitGroup
	//var wgPool sync.WaitGroup
	//wgWorkers.Add(WORKERS)

	//var addResultMutex sync.Mutex
	//
	//taskChan := make(chan int)
	//resChan := make(chan models.SCParameter)

	//Spin up workers ===
	//for i := 0; i < WORKERS; i++ {
	//	go asyncWorker(i, taskChan, host, &addResultMutex, &wgWorkers, &result, cfg)
	//}

	// Send tasks to him ===
	for _, i := range resultId {
		r := worker(i, host, cfg)
		result.SmartClasses = append(result.SmartClasses, r)
		//wgPool.Add(1)
		//go func(_id int, wg *sync.WaitGroup) {
		//	defer wg.Done()
		//	taskChan <- _id
		//	addResult(<-resChan, r, &addResultMutex)
		//}(i, &wgPool)
	}
	//wgPool.Wait()

	// Sort by ID ===
	sort.Slice(result.SmartClasses, func(i, j int) bool {
		return result.SmartClasses[i].ID < result.SmartClasses[j].ID
	})

	return result.SmartClasses, nil
}

func addResult(i models.SCParameter, r *SmartClasses, mtx *sync.Mutex) {
	mtx.Lock()
	if i.OverrideValuesCount > 0 {
		fmt.Println(i.ID, i.OverrideValues)
		fmt.Println("====")
	}
	r.SmartClasses = append(r.SmartClasses, i)
	mtx.Unlock()
}

func worker(i int,
	host string,
	cfg *models.Config) models.SCParameter {
	var d models.SCParameter
	fmt.Printf("W: got task, scId: %d, HOST: %s\n", i, host)

	uri := fmt.Sprintf("smart_class_parameters/%d", i)
	response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
	if response.StatusCode != 200 {
		fmt.Println("SC Parameters, ID:", i, response.StatusCode, host)
	}

	err := json.Unmarshal(response.Body, &d)
	if err != nil {
		logger.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
	}
	return d
}

func asyncWorker(wrkID int,
	in <-chan int,
	host string,
	mtx *sync.Mutex,
	wg *sync.WaitGroup,
	r *SmartClasses,
	cfg *models.Config) {
	defer wg.Done()
	var d models.SCParameter
	for {
		i := <-in

		//fmt.Printf("W: %d got task, scId: %d, HOST: %s\n", wrkID, i, host)

		uri := fmt.Sprintf("smart_class_parameters/%d", i)
		response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
		if response.StatusCode != 200 {
			fmt.Println("SC Parameters, ID:", i, response.StatusCode, host)
		}

		err := json.Unmarshal(response.Body, &d)
		if err != nil {
			logger.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
		} else {
			if d.OverrideValuesCount > 0 {
				fmt.Println("WRK:", wrkID, i, d.OverrideValues)
				fmt.Println("====")
			}
			addResult(d, r, mtx)
		}
	}
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
