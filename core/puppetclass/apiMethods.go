package puppetclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"runtime"
	"sort"
	"sync"
)

// ===============
// GET
// ===============
// Get all Puppet Classes and insert to base
func ApiAll(host string, cfg *models.Config) (map[string][]models.PuppetClass, error) {

	var pcResult models.PuppetClasses
	var result = make(map[string][]models.PuppetClass)

	// check items count
	uri := fmt.Sprintf("puppetclasses?format=json&per_page=%d", cfg.Api.GetPerPage)
	response, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(response.Body, &pcResult)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, response)
		}

		if pcResult.Total > cfg.Api.GetPerPage {
			pagesRange := utils.Pager(pcResult.Total, cfg.Api.GetPerPage)
			for i := 1; i <= pagesRange; i++ {
				uri := fmt.Sprintf("puppetclasses?format=json&page=%d&per_page=%d", i, cfg.Api.GetPerPage)
				response, err := logger.ForemanAPI("GET", host, uri, "", cfg)
				if err == nil {
					err := json.Unmarshal(response.Body, &pcResult)
					if err != nil {
						return result, err
					}

					for className, class := range pcResult.Results {
						for _, subClass := range class {
							result[className] = append(result[className], subClass)
						}
					}
				}
			}
		} else {
			for className, class := range pcResult.Results {
				for _, subClass := range class {
					result[className] = append(result[className], subClass)
				}
			}
		}
	}

	return result, nil
}

// Get Puppet Classes by host group and insert to Host Group
func ApiByHG(host string, hgID int, bdId int, cfg *models.Config) []int {
	var result models.PuppetClasses
	var foremanSCIds []int

	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	response, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(response.Body, &result)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, response)
		}
		var pcIDs []int
		for className, cl := range result.Results {
			for _, subclass := range cl {
				foremanSCIds = append(foremanSCIds, subclass.ID)
				lastId := DbInsert(host, className, subclass.Name, subclass.ID, cfg)
				if lastId != -1 {
					pcIDs = append(pcIDs, lastId)
				}
			}
		}
		DbUpdatePcID(bdId, pcIDs, cfg)
	}
	return foremanSCIds
}

type PuppetClassesRes struct {
	PC []models.PCSCParameters
}

// Just get Puppet Classes by host group
func ApiByHGJson(host string, hgID int, cfg *models.Config) map[string][]models.PuppetClass {

	var result models.PuppetClasses

	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	response, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(response.Body, &result)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, response)
		}
	} else {
		logger.Warning.Printf("%q: getPCByHgJson", err)
	}
	return result.Results
}

//Update Smart Class ids in Puppet Classes
func UpdateSCID(host string, cfg *models.Config) {
	var ids []int

	WORKERS := runtime.NumCPU()
	PuppetClasses := DbAll(host, cfg)

	for _, pc := range PuppetClasses {
		ids = append(ids, pc.ForemanId)
	}

	var result []models.PCSCParameters
	collector := StartDispatcher(WORKERS)
	for _, job := range CreateJobs(ids, host, &result, cfg) {
		collector.Work <- Work{
			ID:        job.ID,
			ForemanID: job.ForemanID,
			Host:      job.Host,
			Results:   job.Results,
			Cfg:       job.Cfg,
		}
	}

	// Store that ===
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	for _, pc := range result {
		DbUpdate(host, pc, cfg)
	}
}

func addResult(i models.PCSCParameters, r *PuppetClassesRes, mtx *sync.Mutex) {
	mtx.Lock()
	r.PC = append(r.PC, i)
	mtx.Unlock()
}
