package smartclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: Started",
		Host:    host,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    host,
			Actions: "smartClasses",
			Status:  ctx.Session.UserName,
			State:   "started",
		},
	})
	// ---

	beforeUpdate := GetForemanIDs(host, ctx)

	smartClassesResult, err := GetAll(host, ctx)
	if err != nil {
		utils.Warning.Printf("error on getting Smart Classes and Overrides:\n%q", err)
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Storing Smart classes",
		Host:    host,
	}))

	aLen := len(smartClassesResult)
	bLen := len(beforeUpdate)

	var afterUpdate = make([]int, 0, aLen)

	for _, i := range smartClassesResult {
		afterUpdate = append(afterUpdate, i.ID)
		InsertSC(host, i, ctx)
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Checking Smart classes",
		Host:    host,
	}))

	sort.Ints(beforeUpdate)
	sort.Ints(afterUpdate)

	if aLen != bLen {
		for _, i := range beforeUpdate {
			if !utils.Search(afterUpdate, i) {
				DeleteSmartClass(host, i, ctx)
			}
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    host,
			Actions: "smartClasses",
			Status:  ctx.Session.UserName,
			State:   "done",
		},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: Done",
		Host:    host,
	}))
}
