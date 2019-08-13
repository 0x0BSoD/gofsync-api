package API

import (
	"encoding/json"
	"fmt"
	envDB "git.ringcentral.com/archops/goFsync/core/environment/DB"
	"git.ringcentral.com/archops/goFsync/core/puppetclass/DB"
	scDB "git.ringcentral.com/archops/goFsync/core/smartclass/DB"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"strconv"
	"strings"
	"sync"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

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
			return map[string][]PuppetClass{}, err
		}
		for className, class := range pcResult.Results {
			for _, subClass := range class {
				result[className] = append(result[className], subClass)
			}
		}
	}

	return result, nil
}

// Return Puppet Class by Foreman ID
func (Get) ByID(host string, pcId int, ctx *user.GlobalCTX) (map[string][]PuppetClass, error) {

	// VARS
	var pcResult PuppetClasses

	// =======
	uri := fmt.Sprintf("puppetclasses/%d", pcId)
	response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(response.Body, &pcResult)
	if err != nil {
		utils.Error.Printf("%q:\n %q\n", err, response)
		return map[string][]PuppetClass{}, err
	}
	return pcResult.Results, nil
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

// =====================================================================================================================
// UPDATE
// =====================================================================================================================

func (Update) SmartClassIDs(host string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Match smart classes to puppet class ID's",
		Host:    host,
	}))

	// VARS
	var ids []int
	var gAPI Get
	var iDB Insert

	PuppetClasses, err := gAPI.All(host, ctx)
	if err != nil {
		utils.Error.Printf("%q: error while getting puppet class data", err)
	}
	for _, PuppetClass := range PuppetClasses {
		for _, SubClass := range PuppetClass {
			ids = append(ids, SubClass.ForemanID)
		}
	}
	var r []PuppetClassDetailed
	var lock sync.Mutex

	// ============
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
				var tmp PuppetClassDetailed
				uri := fmt.Sprintf("puppetclasses/%d", ID)
				response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
				if response.StatusCode != 200 {
					fmt.Println("PuppetClasses updates, ID:", ID, response.StatusCode, host)
				}
				err := json.Unmarshal(response.Body, &tmp)
				if err != nil {
					utils.Error.Printf("%q:\n %q\n", err, response)
				}

				lock.Lock()
				r = append(r, tmp)
				lock.Unlock()

			}
		}(j)
	}
	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	for _, pc := range r {
		iDB.byID(host, pc, ctx)
	}
}

// =====================================================================================================================
// INSERT
// =====================================================================================================================

// Get and Insert to base by host group ID
func (Insert) Add(host string, hgID int, bdId int, ctx *user.GlobalCTX) {

	// VARS
	var result PuppetClasses
	var iDB DB.Insert
	var uDB DB.Update

	// ======
	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	response, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err == nil {
		err := json.Unmarshal(response.Body, &result)
		if err != nil {
			utils.Error.Printf("%q:\n %q\n", err, response)
		}
		var pcIDs []int
		for className, cl := range result.Results {
			for _, subclass := range cl {
				lastId := iDB.Insert(host, className, subclass.Name, subclass.ForemanID, ctx)
				if lastId != -1 {
					pcIDs = append(pcIDs, lastId)
				}
			}
		}
		uDB.HostGroupIDs(bdId, pcIDs, ctx)
	}
}

// Update puppet class in database and return id
func (Insert) byID(host string, parameters PuppetClassDetailed, ctx *user.GlobalCTX) {

	// VARS
	var (
		strScList  []string
		strEnvList []string
		scGDB      scDB.Get
		envGDB     envDB.Get
	)

	// =======
	for _, i := range parameters.SmartClassParameters {
		scID := scGDB.IDByForemanID(host, i.ForemanID, ctx)
		if scID != -1 {
			strScList = append(strScList, strconv.Itoa(scID))
		}
	}

	for _, i := range parameters.Environments {
		ID := envGDB.ID(host, i.Name, ctx)
		if ID != -1 {
			strEnvList = append(strEnvList, strconv.Itoa(ID))
		}
	}

	stmt, err := ctx.Config.Database.DB.Prepare("update puppet_classes set sc_ids=?, env_ids=? where host=? and foreman_id=?")
	if err != nil {
		utils.Error.Printf("%q, error while updating puppet class", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(
		strings.Join(strScList, ","),
		strings.Join(strEnvList, ","),
		host,
		parameters.ForemanID)
	if err != nil {
		utils.Warning.Printf("%q, error while updating puppet class", err)
	}
}
