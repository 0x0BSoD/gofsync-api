package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	. "github.com/0x0BSoD/splitter"
	"sort"
	"sync"
)

// ===============
// INSERT
// ===============
// Get Smart Classes from Foreman
func GetAll(host string, cfg *models.Config) ([]models.SCParameter, error) {
	var r models.SCParameters
	var resultId []int
	var result []models.SCParameter

	uri := fmt.Sprintf("smart_class_parameters?per_page=%d", cfg.Api.GetPerPage)
	response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}
	for _, i := range r.Results {
		resultId = append(resultId, i.ID)
	}

	sort.Ints(resultId)
	var wg sync.WaitGroup

	tasks := make(chan int)
	resChan := make(chan models.SCParameter)
	WORKERS := 6
	queue := SplitToQueue(resultId, WORKERS)
	wg.Add(WORKERS)

	// Spin up workers ===
	for i := 0; i < WORKERS; i++ {
		go worker(i, tasks, resChan, host, &wg, cfg)
	}

	// Send tasks to him ===
	var wgPool sync.WaitGroup
	var lock sync.Mutex
	for i := range queue {
		wgPool.Add(1)
		go func(ids []int, r *[]models.SCParameter, wg *sync.WaitGroup) {
			defer wg.Done()
			for i := 0; i < len(ids); i++ {
				tasks <- ids[i]
				lock.Lock() // w/o lock values may drop from result because 'condition race'
				*r = append(*r, <-resChan)
				lock.Unlock()
			}
		}(queue[i], &result, &wgPool)
	}
	wgPool.Wait()

	// Sort by ID ===
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

func worker(wrkID int,
	in <-chan int,
	out chan<- models.SCParameter,
	host string,
	wg *sync.WaitGroup,
	cfg *models.Config) {
	defer wg.Done()
	var d models.SCParameter
	for {
		i := <-in

		uri := fmt.Sprintf("smart_class_parameters/%d", i)
		response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
		if response.StatusCode != 200 {
			fmt.Println(i, response.StatusCode)
		}

		err := json.Unmarshal(response.Body, &d)
		if err != nil {
			logger.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
		} else {
			out <- d
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
