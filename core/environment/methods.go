package environment

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments",
		Host:    host,
	}))

	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    host,
			Actions: "Getting Environments",
			State:   "",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
	// ---

	beforeUpdate := DbAll(host, ctx)
	var afterUpdate []string

	environmentsResult, err := ApiAll(host, ctx)
	if err != nil {
		logger.Warning.Printf("Error on getting Environments:\n%q", err)
	}

	sort.Slice(environmentsResult.Results, func(i, j int) bool {
		return environmentsResult.Results[i].ID < environmentsResult.Results[j].ID
	})

	for _, env := range environmentsResult.Results {

		// Socket Broadcast ---
		if ctx.Session.PumpStarted {
			data := models.Step{
				Host:    host,
				Actions: "Saving Environments",
				State:   fmt.Sprintf("Parameter: %s", env.Name),
			}
			msg, _ := json.Marshal(data)
			ctx.Session.SendMsg(msg)
		}
		// ---

		DbInsert(host, env.Name, env.ID, ctx)
		afterUpdate = append(afterUpdate, env.Name)
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i, afterUpdate) {
			DbDelete(host, i, ctx)
		}
	}
}
