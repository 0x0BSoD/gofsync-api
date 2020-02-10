package locations

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(hostname string, ctx *user.GlobalCTX) {

	// Step LOG to stdout ======================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Locations :: Start",
		Host:    hostname,
	}))
	// =========================================

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.Location,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting Locations from foreman"},
	})
	// ---

	// from the DB
	beforeUpdate, _ := DbAll(ctx.Config.Hosts[hostname], ctx)

	// from a foreman
	locationsResult, err := ApiAll(hostname, ctx)
	if err != nil {
		utils.Warning.Printf("Error on getting Locations:\n%q", err)
	}

	sort.Slice(locationsResult.Results, func(i, j int) bool {
		return locationsResult.Results[i].ID < locationsResult.Results[j].ID
	})

	// store
	aLen := len(locationsResult.Results)
	bLen := len(beforeUpdate)

	var afterUpdate = make([]string, 0, aLen)
	count := 1
	for _, loc := range locationsResult.Results {

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			HostName:  hostname,
			Resource:  models.Environment,
			Operation: "sync",
			UserName:  ctx.Session.UserName,
			AdditionalData: models.CommonOperation{
				Message: "Saving Location",
				Item:    loc.Name,
				Total:   aLen,
				Current: count,
			},
		})
		// ---

		DbInsert(ctx.Config.Hosts[hostname], loc.ID, loc.Name, ctx)
		afterUpdate = append(afterUpdate, loc.Name)
		count++
	}
	sort.Strings(afterUpdate)

	// delete if don't have any errors
	if err == nil {
		if aLen != bLen {
			for _, i := range beforeUpdate {
				if !utils.StringInSlice(i, afterUpdate) {
					DbDelete(ctx.Config.Hosts[hostname], i, ctx)
				}
			}
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.Location,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting Locations from foreman done", Done: true},
	})
	// ---

	// Step LOG to stdout ======================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Locations :: Done",
		Host:    hostname,
	}))
	// =========================================
}
