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
	var ids []int
	var result []models.SCParameter
	var writeLock sync.Mutex

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
					return result, err
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
	WORKERS := runtime.NumCPU()
	collector := StartDispatcher(WORKERS)
	var wg sync.WaitGroup

	for _, job := range CreateJobs(ids, host, &result, cfg) {
		wg.Add(1)
		collector.Work <- Work{
			ID:        job.ID,
			ForemanID: job.ForemanID,
			Host:      job.Host,
			Results:   job.Results,
			Cfg:       job.Cfg,
			Lock:      &writeLock,
			Wg:        &wg,
		}
	}

	return result, nil
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
