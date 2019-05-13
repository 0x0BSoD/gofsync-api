package environment

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

func Sync(host string, cfg *models.Config) {
	// Locations ===
	//==========================================================================================================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments",
		Host:    host,
	}))

	beforeUpdate := GetEnvList(host, cfg)
	var afterUpdate []string

	environmentsResult, err := Environments(host, cfg)
	if err != nil {
		logger.Warning.Printf("Error on getting Environments:\n%q", err)
	}

	for _, env := range environmentsResult.Results {
		InsertToEnvironments(host, env.Name, env.ID, cfg)
		afterUpdate = append(afterUpdate, env.Name)
	}

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i, afterUpdate) {
			DeleteEnvironment(host, i, cfg)
		}
	}
}
