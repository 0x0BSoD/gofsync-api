package API

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// Return all puppet classes by host
func (Get) All(host string, ctx *user.GlobalCTX) (map[string][]PuppetClass, error) {

	// VARS
	var pcResult PuppetClasses
	var result = make(map[string][]PuppetClass)

	// =====
	response, err := utils.ForemanAPI("GET", host, "puppetclasses?format=json", "", ctx)
	if err == nil {
		err := json.Unmarshal(response.Body, &pcResult)
		if err != nil {
			utils.Error.Printf("%q:\n %q\n", err, response)
		}
		for className, class := range pcResult.Results {
			for _, subClass := range class {
				result[className] = append(result[className], subClass)
			}
		}
	}

	return result, nil
}

// Return all puppet classes by host group id
func (Get) ByHostGroupID(host string, hgID int, bdId int, ctx *user.GlobalCTX) (map[string][]PuppetClass, error) {

	// VARS
	var pcResult PuppetClasses
	var result = make(map[string][]PuppetClass)

	// ======
	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	response, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err == nil {
		err := json.Unmarshal(response.Body, &pcResult)
		if err != nil {
			utils.Error.Printf("%q:\n %q\n", err, response)
			return nil, err
		}
		for className, class := range pcResult.Results {
			for _, subClass := range class {
				result[className] = append(result[className], subClass)
			}
		}
		return result, nil
	} else {
		return nil, err
	}

}

//func UpdateSCID(host string, ctx *user.GlobalCTX) {
//
//	fmt.Println(utils.PrintJsonStep(models.Step{
//		Actions: "Match smart classes to puppet class ID's",
//		Host:    host,
//	}))
//
//	var ids []int
//	PuppetClasses := DbAll(host, ctx)
//	for _, pc := range PuppetClasses {
//		ids = append(ids, pc.ForemanId)
//	}
//
//	var r PCResult
//
//	// ver 2 ===
//	// Create a new WorkQueue.
//	wq := utils.New()
//	// This sync.WaitGroup is to make sure we wait until all of our work
//	// is done.
//	var wg sync.WaitGroup
//
//	//fmt.Println(len(ids))
//
//	for _, j := range ids {
//		wg.Add(1)
//		go func(ID int) {
//			wq <- func() {
//				defer wg.Done()
//				var tmp smartclass.PCSCParameters
//				uri := fmt.Sprintf("puppetclasses/%d", ID)
//				response, _ := logger.ForemanAPI("GET", host, uri, "", ctx)
//				if response.StatusCode != 200 {
//					fmt.Println("PuppetClasses updates, ID:", ID, response.StatusCode, host)
//				}
//				err := json.Unmarshal(response.Body, &tmp)
//				if err != nil {
//					logger.Error.Printf("%q:\n %q\n", err, response)
//				}
//
//				r.Add(tmp)
//
//			}
//		}(j)
//	}
//	// Wait for all of the work to finish, then close the WorkQueue.
//	wg.Wait()
//	close(wq)
//
//	//fmt.Println(len(r.resSlice))
//
//	for _, pc := range r.resSlice {
//		DbUpdate(host, pc, ctx)
//	}
//}
