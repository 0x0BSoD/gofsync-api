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
	msg := models.Step{
		Host:    host,
		Actions: "Getting Puppet Classes",
		State:   "",
	}
	utils.CastMsgToUser(ctx, msg)
	// ---

	beforeUpdate := DbAll(host, ctx)
	var afterUpdate []string

	getAllPCResult, err := ApiAll(host, ctx)
	if err != nil {
		logger.Warning.Printf("Error on getting Puppet classes:\n%q", err)
	}

	count := 1
	for className, subClasses := range getAllPCResult {

		// Socket Broadcast ---
		msg := models.Step{
			Host:    host,
			Actions: "Saving Puppet Class",
			State:   fmt.Sprintf("Puppet Class: %s %d/%d", className, count, len(getAllPCResult)),
		}
		utils.CastMsgToUser(ctx, msg)
		// ---

		for _, subClass := range subClasses {
			DbInsert(host, className, subClass.Name, subClass.ID, ctx)
			afterUpdate = append(afterUpdate, subClass.Name)
		}
		count++
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i.Subclass, afterUpdate) {
			DeletePuppetClass(host, i.Subclass, ctx)
		}
	}
}
