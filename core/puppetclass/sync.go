package puppetclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(hostname string, ctx *user.GlobalCTX) {

	hostID := ctx.Config.Hosts[hostname]

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Puppet classes :: Start",
		Host:    hostname,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.PuppetClass,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting Puppet Classes from foreman"},
	})
	// ---

	allPuppetClasses := DbAll(hostID, ctx)
	beforeUpdate := make([]int, 0, len(allPuppetClasses))
	for _, i := range allPuppetClasses {
		beforeUpdate = append(beforeUpdate, i.ForemanId)
	}

	getAllPCResult, err := ApiAll(hostname, ctx)
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
			HostName:  hostname,
			Resource:  models.PuppetClass,
			Operation: "sync",
			UserName:  ctx.Session.UserName,
			AdditionalData: models.CommonOperation{
				Message: "Saving PuppetClass",
				Item:    className,
				Total:   subclassesLen,
				Current: count,
			},
		})
		// ---

		subclassesLen := len(subClasses)
		updated := make([]int, 0, subclassesLen)
		count2 := 0
		for _, subClass := range subClasses {
			// Socket Broadcast ---
			ctx.Session.SendMsg(models.WSMessage{
				Broadcast: false,
				HostName:  hostname,
				Resource:  models.PuppetClass,
				Operation: "sync",
				UserName:  ctx.Session.UserName,
				AdditionalData: models.CommonOperation{
					Message: "Saving PuppetClass subclass",
					Item:    subClass.Name,
					Total:   len(subClasses),
					Current: count2,
				},
			})
			// ---
			fmt.Printf("{INSERT PC} %s || %s \n", className, subClass.Name)
			DbInsert(hostID, subClass.ForemanID, className, subClass.Name, ctx)
			updated = append(updated, subClass.ForemanID)
			count2++
		}
		count++
		afterUpdate = append(afterUpdate, updated...)
	}

	sort.Ints(afterUpdate)
	sort.Ints(beforeUpdate)

	fmt.Println("{Deleting PC}")
	for _, i := range beforeUpdate {
		fmt.Println(i)
		if !utils.Search(afterUpdate, i) {
			fmt.Println("GOT:", i)
			DeletePuppetClass(hostID, i, ctx)
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.PuppetClass,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting Puppet Classes from foreman", Done: true},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Puppet classes :: Done",
		Host:    hostname,
	}))
}
