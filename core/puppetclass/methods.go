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
		Actions: "Getting Puppet classes :: Start",
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

	allPuppetClasses := DbAll(host, ctx)
	beforeUpdate := make([]int, 0, len(allPuppetClasses))
	for _, i := range allPuppetClasses {
		beforeUpdate = append(beforeUpdate, i.ForemanId)
	}

	getAllPCResult, err := ApiAll(host, ctx)
	if err != nil {
		logger.Warning.Printf("Error on getting Puppet classes:\n%q", err)
	}

	count := 1

	subclassesLen := len(getAllPCResult)
	afterUpdate := make([]int, 0, subclassesLen)

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
		updated := make([]int, 0, subclassesLen)
		for _, subClass := range subClasses {
			DbInsert(host, className, subClass.Name, subClass.ForemanID, ctx)
			updated = append(updated, subClass.ForemanID)
		}
		count++
		afterUpdate = append(afterUpdate, updated...)
	}

	sort.Ints(afterUpdate)
	sort.Ints(beforeUpdate)

	for _, i := range beforeUpdate {
		if !utils.Search(afterUpdate, i) {
			DeletePuppetClass(host, i, ctx)
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "done",
	})
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    host,
			Actions: "puppetClasses",
			Status:  ctx.Session.UserName,
			State:   "done",
		},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Puppet classes :: Done",
		Host:    host,
	}))
}
