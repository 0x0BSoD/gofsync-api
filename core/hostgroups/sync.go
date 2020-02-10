package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/foremans"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(hostname string, ctx *user.GlobalCTX) {
	// Host groups ===
	//==========================================================================================================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Filling HostGroups :: Started",
		Host:    hostname,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.HostGroup,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting HostGroups from foreman"},
	})
	// ---

	results := GetHostGroups(hostname, ctx)
	beforeUpdate := ForemanIDs(ctx.Config.Hosts[hostname], ctx)
	aLen := len(results)
	bLen := len(beforeUpdate)
	var afterUpdate = make([]int, 0, aLen)

	// RT SWEs =================================================================================================
	swes := RTBuildObj(foremans.PuppetHostEnv(ctx.Config.Hosts[hostname], ctx), ctx)

	for idx, i := range results {
		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			HostName:  hostname,
			Resource:  models.HostGroup,
			Operation: "sync",
			UserName:  ctx.Session.UserName,
			AdditionalData: models.CommonOperation{
				Message: "Saving HostGroup",
				State:   "saving",
				Item:    i.Name,
				Current: idx + 1,
				Total:   aLen,
			},
		})
		// ---

		sJson, _ := json.Marshal(i)
		sweStatus := GetFromRT(i.Name, swes)
		fmt.Printf("{INSERT HG} %s || %s \n", i.Name, sweStatus)
		lastId := Insert(ctx.Config.Hosts[hostname], i.ID, i.Name, string(sJson), sweStatus, ctx)
		afterUpdate = append(afterUpdate, i.ID)
		if lastId != -1 {
			puppetclass.ApiByHG(hostname, i.ID, lastId, ctx)
			HgParams(hostname, lastId, i.ID, ctx)
		}
	}

	if aLen != bLen {

		sort.Ints(afterUpdate)
		sort.Ints(beforeUpdate)

		for _, i := range beforeUpdate {
			if !utils.Search(afterUpdate, i) {
				fmt.Println("Deleting ... ", i, hostname)
				name := Name(ctx.Config.Hosts[hostname], i, ctx)
				Delete(ctx.Config.Hosts[hostname], i, ctx)
				rmJSON(name, hostname, ctx)
			}
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.HostGroup,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting HostGroups from foreman", Done: true},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Filling HostGroups :: Done",
		Host:    hostname,
	}))
}
