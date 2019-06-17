package puppetclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
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
// Result struct
type PCResult struct {
	sync.Mutex
	resSlice []models.PCSCParameters
}

func (r *PCResult) Add(ID models.PCSCParameters) {
	r.Lock()
	r.resSlice = append(r.resSlice, ID)
	r.Unlock()
}
func UpdateSCID(host string, cfg *models.Config) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Match smart classes to puppet class ID's",
		Host:    host,
	}))

	var ids []int
	PuppetClasses := DbAll(host, cfg)
	for _, pc := range PuppetClasses {
		ids = append(ids, pc.ForemanId)
	}

	var r PCResult

	// ver 2 ===
	// Create a new WorkQueue.
	wq := utils.New()
	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	var wg sync.WaitGroup

	fmt.Println(len(ids))

	for _, j := range ids {
		wg.Add(1)
		go func(ID int) {
			wq <- func() {
				defer wg.Done()
				var tmp models.PCSCParameters
				uri := fmt.Sprintf("puppetclasses/%d", ID)
				response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
				if response.StatusCode != 200 {
					fmt.Println("PuppetClasses updates, ID:", ID, response.StatusCode, host)
				}
				err := json.Unmarshal(response.Body, &r)
				if err != nil {
					logger.Error.Printf("%q:\n %q\n", err, response)
				}

				r.Add(tmp)

			}
		}(j)
	}
	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	fmt.Println(len(r.resSlice))

	for _, pc := range r.resSlice {
		DbUpdate(host, pc, cfg)
	}
}
