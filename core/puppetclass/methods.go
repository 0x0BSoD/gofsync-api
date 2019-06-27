package puppetclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, ss *models.Session) {
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
	utils.BroadCastMsg(ss, msg)
	// ---

	beforeUpdate := DbAll(host, ss)
	var afterUpdate []string

	getAllPCResult, err := ApiAll(host, ss)
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
		utils.BroadCastMsg(ss, msg)
		// ---

		for _, subClass := range subClasses {
			DbInsert(host, className, subClass.Name, subClass.ID, ss)
			afterUpdate = append(afterUpdate, subClass.Name)
		}
		count++
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i.Subclass, afterUpdate) {
			DeletePuppetClass(host, i.Subclass, ss)
		}
	}
}
