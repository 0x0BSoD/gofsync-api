package puppetclass

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

func Sync(host string, cfg *models.Config) {
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Puppet classes",
		Host:    host,
	}))

	beforeUpdate := GetAllPCDB(host, cfg)
	var afterUpdate []string

	getAllPCResult, err := GetAllPC(host, cfg)
	if err != nil {
		logger.Warning.Printf("Error on getting Puppet classes:\n%q", err)
	}

	for className, subClasses := range getAllPCResult {
		for _, subClass := range subClasses {
			InsertPC(host, className, subClass.Name, subClass.ID, cfg)
			afterUpdate = append(afterUpdate, subClass.Name)
		}
	}

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i.Subclass, afterUpdate) {
			DeletePuppetClass(host, i.Subclass, cfg)
		}
	}
}
