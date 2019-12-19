package smartclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
)

func Sync(hostname string, ctx *user.GlobalCTX) {

	hostID := ctx.Config.Hosts[hostname]

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: Started",
		Host:    hostname,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    hostname,
			Actions: "smartClasses",
			Status:  ctx.Session.UserName,
			State:   "started",
		},
	})
	// ---

	//beforeUpdate := GetForemanIDs(hostID, ctx)

	smartClassesResult, err := GetAll(hostname, ctx)
	if err != nil {
		utils.Warning.Printf("error on getting Smart Classes and Overrides:\n%q", err)
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Storing Smart classes",
		Host:    hostname,
	}))

	//aLen := len(smartClassesResult)
	//bLen := len(beforeUpdate)

	//var afterUpdate = make([]int, 0, aLen)

	for _, i := range smartClassesResult {
		//afterUpdate = append(afterUpdate, i.ID)
		fmt.Printf("{INSERT SC} %s || %s \n", i.Parameter, i.PuppetClass.Name)
		InsertSC(hostID, i, ctx)
	}

	//fmt.Println(utils.PrintJsonStep(models.Step{
	//	Actions: "Checking Smart classes",
	//	Host:    hostname,
	//}))
	//
	//sort.Ints(beforeUpdate)
	//sort.Ints(afterUpdate)
	//
	//if aLen != bLen {
	//	for _, i := range beforeUpdate {
	//		if !utils.Search(afterUpdate, i) {
	//			DeleteSmartClass(hostID, i, ctx)
	//		}
	//	}
	//}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    hostname,
			Actions: "smartClasses",
			Status:  ctx.Session.UserName,
			State:   "done",
		},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: Done",
		Host:    hostname,
	}))
}
