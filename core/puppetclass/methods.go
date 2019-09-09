package puppetclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Puppet classes",
		Host:    host,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    host,
			Actions: "puppetClasses",
			Status:  ctx.Session.UserName,
			State:   "started",
		},
	})
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "getPC",
		Data: models.Step{
			Host:  host,
			State: "running",
		},
	})
	// ---

	beforeUpdate := DbAll(host, ctx)

	getAllPCResult, err := ApiAll(host, ctx)
	if err != nil {
		logger.Warning.Printf("Error on getting Puppet classes:\n%q", err)
	}

	count := 1
	subclassesLen := len(getAllPCResult)
	afterUpdate := make([]string, 0, subclassesLen)
	for className, subClasses := range getAllPCResult {

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			Operation: "getPC",
			Data: models.Step{
				Host:  host,
				State: "saving",
				Item:  className,
				Counter: struct {
					Current int `json:"current"`
					Total   int `json:"total"`
				}{count, len(getAllPCResult)},
			},
		})
		// ---

		subclassesLen := len(subClasses)
		updated := make([]string, 0, subclassesLen)
		for _, subClass := range subClasses {
			DbInsert(host, className, subClass.Name, subClass.ID, ctx)
			updated = append(updated, subClass.Name)
		}
		count++
		afterUpdate = append(afterUpdate, updated...)
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i.Subclass, afterUpdate) {
			DeletePuppetClass(host, i.Subclass, ctx)
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "done",
	})
	// ---
}
