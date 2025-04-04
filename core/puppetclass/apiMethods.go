package puppetclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sync"
)

// ===============
// GET
// ===============
// Get all Puppet Classes and insert to base
func ApiAll(hostname string, ctx *user.GlobalCTX) (map[string][]PuppetClass, error) {

	var pcResult PuppetClasses
	var result = make(map[string][]PuppetClass)

	// check items count
	uri := fmt.Sprintf("puppetclasses?format=json&per_page=%d", ctx.Config.Api.GetPerPage)
	response, err := logger.ForemanAPI("GET", hostname, uri, "", ctx)
	if err == nil {
		err := json.Unmarshal(response.Body, &pcResult)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, response)
		}

		if pcResult.Total > ctx.Config.Api.GetPerPage {
			pagesRange := utils.Pager(pcResult.Total, ctx.Config.Api.GetPerPage)
			for i := 1; i <= pagesRange; i++ {
				uri := fmt.Sprintf("puppetclasses?format=json&page=%d&per_page=%d", i, ctx.Config.Api.GetPerPage)
				response, err := logger.ForemanAPI("GET", hostname, uri, "", ctx)
				if err == nil {
					err := json.Unmarshal(response.Body, &pcResult)
					if err != nil {
						return result, err
					}

					for className, class := range pcResult.Results {
						result[className] = append(result[className], class...)
					}
				}
			}
		} else {
			for className, class := range pcResult.Results {
				result[className] = append(result[className], class...)
			}
		}
	}

	return result, nil
}

// Get Puppet Classes by host group and insert to Host Group
func ApiByHG(hostname string, hgID int, bdId int, ctx *user.GlobalCTX) []int {

	var result PuppetClasses
	var foremanSCIds []int

	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	response, err := logger.ForemanAPI("GET", hostname, uri, "", ctx)
	if err == nil {
		err := json.Unmarshal(response.Body, &result)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, response)
		}
		var pcIDs []int
		for className, cl := range result.Results {
			for _, subclass := range cl {
				foremanSCIds = append(foremanSCIds, subclass.ForemanID)
				lastId := DbInsert(ctx.Config.Hosts[hostname], subclass.ForemanID, className, subclass.Name, ctx)
				if lastId != -1 {
					pcIDs = append(pcIDs, lastId)
				}
			}
		}
		DbUpdatePcID(bdId, pcIDs, ctx)
	}
	return foremanSCIds
}

// Just get Puppet Classes by host group
func ApiByHGJson(hostname string, hgID int, ctx *user.GlobalCTX) map[string][]PuppetClass {

	var result PuppetClasses

	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	response, err := logger.ForemanAPI("GET", hostname, uri, "", ctx)
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
	resSlice []smartclass.PCSCParameters
}

func (r *PCResult) Add(pc smartclass.PCSCParameters) {
	r.Lock()
	r.resSlice = append(r.resSlice, pc)
	r.Unlock()
}

func UpdateSCID(hostname string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Match smart classes to puppet class ID's :: Started",
		Host:    hostname,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.PuppetClass,
		Operation:      "updatingSmartClassesIDs",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Match smart classes to puppet class ID's"},
	})
	// ---

	PuppetClasses := DbAll(ctx.Config.Hosts[hostname], ctx)
	var ids = make([]int, 0, len(PuppetClasses))
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

	for _, j := range ids {
		wg.Add(1)
		go func(ID int) {
			wq <- func() {
				defer wg.Done()
				var tmp smartclass.PCSCParameters
				uri := fmt.Sprintf("puppetclasses/%d", ID)
				response, _ := logger.ForemanAPI("GET", hostname, uri, "", ctx)
				if response.StatusCode != 200 {
					fmt.Println("PuppetClasses updates, ID:", ID, response.StatusCode, hostname)
				}
				err := json.Unmarshal(response.Body, &tmp)
				if err != nil {
					logger.Error.Printf("%q:\n %q\n", err, response)
				}

				r.Add(tmp)
			}
		}(j)
	}
	// Wait for all the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	for _, pc := range r.resSlice {
		DbUpdate(ctx.Config.Hosts[hostname], pc, ctx)
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.PuppetClass,
		Operation:      "updatingSmartClassesIDs",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Match smart classes to puppet class ID's", Done: true},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Match smart classes to puppet class ID's :: Done",
		Host:    hostname,
	}))
}
