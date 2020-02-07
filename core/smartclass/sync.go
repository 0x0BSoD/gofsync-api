package smartclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(hostname string, ctx *user.GlobalCTX) {

	hostID := ctx.Config.Hosts[hostname]

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: Started",
		Host:    hostname,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.SmartClass,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting Smart classes from foreman"},
	})
	// ---

	beforeUpdate := GetForemanIDs(hostID, ctx)

	smartClassesResult, err := GetAll(hostname, ctx)
	if err != nil {
		utils.Warning.Printf("error on getting Smart Classes and Overrides:\n%q", err)
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Storing Smart classes",
		Host:    hostname,
	}))

	aLen := len(smartClassesResult)
	bLen := len(beforeUpdate)

	var afterUpdate []int

	count := 0
	for _, i := range smartClassesResult {
		_, err := InsertSC(hostID, i, ctx)
		if err == nil {
			// Socket Broadcast ---
			ctx.Session.SendMsg(models.WSMessage{
				Broadcast: false,
				HostName:  hostname,
				Resource:  models.SmartClass,
				Operation: "sync",
				UserName:  ctx.Session.UserName,
				AdditionalData: models.CommonOperation{
					Message: "Saving SmartClass parameter",
					Item:    i.Parameter,
					Total:   aLen,
					Current: count,
				},
			})
			// ---
			fmt.Printf("{INSERT SC} %s || %s \n", i.Parameter, i.PuppetClass.Name)
			afterUpdate = append(afterUpdate, i.ID)
			count++
		} else {
			utils.Warning.Printf("error on inserting Smart Classes and Overrides:\n%q", err)
		}
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Checking Smart classes",
		Host:    hostname,
	}))

	sort.Ints(beforeUpdate)
	sort.Ints(afterUpdate)

	if aLen != bLen {
		for _, i := range beforeUpdate {
			if !utils.Search(afterUpdate, i) {
				fmt.Printf("{DELETE SC} ForemanID: %d\n", i)
				err := DeleteSmartClass(hostID, i, ctx)
				if err != nil {
					utils.Warning.Printf("error on deleting Smart Class:\n%q", err)
				}
			}
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.SmartClass,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting Smart classes from foreman", Done: true},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes :: Done",
		Host:    hostname,
	}))
}
