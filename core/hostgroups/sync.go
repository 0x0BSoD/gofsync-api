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
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    hostname,
			Actions: "hostGroups",
			Status:  ctx.Session.UserName,
			State:   "started",
		},
	})
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "getHG",
		Data: models.Step{
			Host:  hostname,
			State: "running",
		},
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
			Operation: "submitHG",
			Data: models.Step{
				Host:  hostname,
				Item:  i.Name,
				State: "saving",
				Counter: struct {
					Current int `json:"current"`
					Total   int `json:"total"`
				}{idx + 1, len(results)},
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
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    hostname,
			Actions: "hostGroups",
			Status:  ctx.Session.UserName,
			State:   "done",
		},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Filling HostGroups :: Done",
		Host:    hostname,
	}))
}
