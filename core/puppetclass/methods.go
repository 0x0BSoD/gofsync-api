package puppetclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/puppetclass/API"
	"git.ringcentral.com/archops/goFsync/core/puppetclass/DB"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
	"sync"
)

func Sync(host string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Puppet classes",
		Host:    host,
	}))

	// VARS
	var (
		aGet API.Get
		dGet DB.Get
		dIns DB.Insert
		dDel DB.Delete
	)

	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		msg, _ := json.Marshal(models.Step{
			Host:    host,
			Actions: "Getting Puppet Classes",
			State:   "",
		})
		ctx.Session.SendMsg(msg)
	}
	// ---

	beforeUpdate := dGet.All(host, ctx)
	var afterUpdate []string

	getAllPCResult, err := aGet.All(host, ctx)
	if err != nil {
		utils.Warning.Printf("Error on getting Puppet classes:\n%q", err)
	}

	count := 1
	for className, subClasses := range getAllPCResult {

		// Socket Broadcast ---
		if ctx.Session.PumpStarted {
			msg, _ := json.Marshal(models.Step{
				Host:    host,
				Actions: "Saving Puppet Class",
				State:   fmt.Sprintf("Puppet Class: %s %d/%d", className, count, len(getAllPCResult)),
			})
			ctx.Session.SendMsg(msg)
		}
		// ---

		for _, subClass := range subClasses {
			dIns.Insert(host, className, subClass.Name, subClass.ForemanID, ctx)
			afterUpdate = append(afterUpdate, subClass.Name)
		}
		count++
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i.Subclass, afterUpdate) {
			err := dDel.BySubclass(host, i.Subclass, ctx)
			if err != nil {
				utils.Warning.Printf("error while deleteing puppet class:\n%q", err)
			}
		}
	}
}

func FillSmartClassIDs(host string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Match smart classes to puppet class ID's",
		Host:    host,
	}))

	// VARS
	var ids []int
	var gAPI DB.Get
	var uDB DB.Update

	// ==============
	PuppetClasses := gAPI.All(host, ctx)
	for _, pc := range PuppetClasses {
		ids = append(ids, pc.ForemanID)
	}

	var r PCResult

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
				var tmp smartclass.PCSCParameters
				uri := fmt.Sprintf("puppetclasses/%d", ID)
				response, _ := utils.ForemanAPI("GET", host, uri, "", ctx)
				if response.StatusCode != 200 {
					fmt.Println("PuppetClasses updates, ID:", ID, response.StatusCode, host)
				}
				err := json.Unmarshal(response.Body, &tmp)
				if err != nil {
					utils.Error.Printf("%q:\n %q\n", err, response)
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
		uDB.SmartClassIDs(host, pc, ctx)
	}
}
