package smartclass

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
		Actions: "Getting Smart classes",
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
	sort.Ints(beforeUpdate)

	smartClassesResult, err := GetAll(host, ctx)
	if err != nil {
		logger.Warning.Printf("Error on getting Smart Classes and Overrides:\n%q", err)
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Storing Smart classes",
		Host:    host,
	}))

	aLen := len(smartClassesResult)
	bLen := len(beforeUpdate)

	var afterUpdate = make([]int, aLen)

	for _, i := range smartClassesResult {
		afterUpdate = append(afterUpdate, i.ID)
		sort.Ints(afterUpdate)
		InsertSC(host, i, ctx)
	}

	if aLen != bLen {
		for _, i := range beforeUpdate {
			fmt.Println(utils.PrintJsonStep(models.Step{
				Actions: fmt.Sprintf("Checking ...%d", i),
				Host:    host,
			}))
			if !utils.Search(afterUpdate, i) {
				fmt.Println(utils.PrintJsonStep(models.Step{
					Actions: fmt.Sprintf("Deleting ...%d", i),
					Host:    host,
				}))
				DeleteSmartClass(host, i, ctx)
			}
		}
	}

}
