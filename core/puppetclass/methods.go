package puppetclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/puppetclass/API"
	"git.ringcentral.com/archops/goFsync/core/puppetclass/DB"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
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
			dIns.Insert(host, className, subClass.Name, subClass.ID, ctx)
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
