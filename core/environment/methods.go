package environment

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, cfg *models.Config) {
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments",
		Host:    host,
	}))

	// Socket Broadcast ---
	msg := models.Step{
		Host:    host,
		Actions: "Getting Environments",
		State:   "",
	}
	utils.BroadCastMsg(cfg, msg)
	// ---

	beforeUpdate := DbAll(host, cfg)
	var afterUpdate []string

	environmentsResult, err := ApiAll(host, cfg)
	if err != nil {
		logger.Warning.Printf("Error on getting Environments:\n%q", err)
	}

	sort.Slice(environmentsResult.Results, func(i, j int) bool {
		return environmentsResult.Results[i].ID < environmentsResult.Results[j].ID
	})

	for _, env := range environmentsResult.Results {
		// Socket Broadcast ---
		msg := models.Step{
			Host:    host,
			Actions: "Saving Environments",
			State:   fmt.Sprintf("Parameter: %s", env.Name),
		}
		utils.BroadCastMsg(cfg, msg)
		// ---
		DbInsert(host, env.Name, env.ID, cfg)
		afterUpdate = append(afterUpdate, env.Name)
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i, afterUpdate) {
			DbDelete(host, i, cfg)
		}
	}
}
