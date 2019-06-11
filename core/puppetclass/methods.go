package puppetclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, cfg *models.Config) {
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
	utils.BroadCastMsg(cfg, msg)
	// ---

	beforeUpdate := DbAll(host, cfg)
	var afterUpdate []string

	getAllPCResult, err := ApiAll(host, cfg)
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
		utils.BroadCastMsg(cfg, msg)
		// ---

		for _, subClass := range subClasses {
			DbInsert(host, className, subClass.Name, subClass.ID, cfg)
			afterUpdate = append(afterUpdate, subClass.Name)
		}
		count++
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i.Subclass, afterUpdate) {
			DeletePuppetClass(host, i.Subclass, cfg)
		}
	}
}
