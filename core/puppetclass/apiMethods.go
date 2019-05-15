package puppetclass

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
// GET
// ===============
// Get all Puppet Classes and insert to base
func GetAllPC(host string, cfg *models.Config) (map[string][]models.PuppetClass, error) {

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
func GetPCByHg(host string, hgID int, bdId int, cfg *models.Config) []int {
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
				lastId := InsertPC(host, className, subclass.Name, subclass.ID, cfg)
				if lastId != -1 {
					pcIDs = append(pcIDs, lastId)
				}
			}
		}
		UpdatePCinHG(bdId, pcIDs, cfg)
	}
	return foremanSCIds
}

// Just get Puppet Classes by host group
func GetPCByHgJson(host string, hgID int, cfg *models.Config) map[string][]models.PuppetClass {

	var result models.PuppetClasses
	//var pcIDs []int64

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
	//var r models.PCSCParameters

	PCss := GetAllPCBase(host, cfg)
	pcIDs := make([]int, 0, len(PCss))
	for k := range PCss {
		pcIDs = append(pcIDs, PCss[k].ForemanID)
	}
	sort.Ints(pcIDs)
	var wg sync.WaitGroup

	tasks := make(chan int)
	resChan := make(chan models.PCSCParameters)
	var data []models.PCSCParameters
	WORKERS := 6
	//queue := SplitToQueue(pcIDs, WORKERS)
	wg.Add(WORKERS)

	//for i:= range queue {
	//	fmt.Println("===================")
	//	fmt.Println(queue[i])
	//}

	//Spin up workers ===
	for i := 0; i < WORKERS; i++ {
		go worker(i, tasks, resChan, host, &wg, cfg)
	}

	//// =====
	var wgPool sync.WaitGroup
	var lock sync.Mutex
	for i := range pcIDs {
		wgPool.Add(1)
		go func(_id int, r *[]models.PCSCParameters, wg *sync.WaitGroup) {
			defer wg.Done()
			//for i := 0; i < len(ids); i++ {
			tasks <- _id
			lock.Lock() // w/o lock values may drop from result because 'condition race'
			*r = append(*r, <-resChan)
			lock.Unlock()
			//}
		}(pcIDs[i], &data, &wgPool)
	}
	wgPool.Wait()
	// =====

	// Store that ===
	for i := range data {
		UpdatePC(host, data[i].Name, data[i], cfg)
	}
}

func worker(wrkID int,
	in <-chan int,
	out chan<- models.PCSCParameters,
	host string,
	wg *sync.WaitGroup,
	cfg *models.Config) {
	defer wg.Done()
	var r models.PCSCParameters
	for {
		i := <-in

		//fmt.Printf("W: %d got task, pcId: %d, HOST: %s\n", wrkID, i, host)

		uri := fmt.Sprintf("puppetclasses/%d", i)
		response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
		if response.StatusCode != 200 {
			fmt.Println("PuppetClasses updates, ID:", i, response.StatusCode, host)
		}

		err := json.Unmarshal(response.Body, &r)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, response)
		}
		out <- r
	}
}
