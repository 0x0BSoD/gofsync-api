package environment

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(hostname string, ctx *user.GlobalCTX) error {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments :: Start",
		Host:    hostname,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    hostname,
			Actions: "environment",
			Status:  ctx.Session.UserName,
			State:   "started",
		},
	})

	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "getEnv",
		Data: models.Step{
			Host:  hostname,
			State: "running",
		},
	})
	// ---

	// ==========
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments :: Getting data",
		Host:    hostname,
	}))

	beforeUpdate := DbGetByHost(ctx.Config.Hosts[hostname], ctx)

	environmentsResult, err := ApiGetAll(hostname, ctx)
	if err != nil {
		utils.Error.Printf("list all Environments:\n%q", err)
		return err
	}

	sort.Slice(environmentsResult.Results, func(i, j int) bool {
		return environmentsResult.Results[i].ID < environmentsResult.Results[j].ID
	})

	aLen := len(environmentsResult.Results)
	bLen := len(beforeUpdate)

	// ==========
	var afterUpdate = make([]string, 0, aLen)
	for _, env := range environmentsResult.Results {

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			Operation: "getEnv",
			Data: models.Step{
				Host:   hostname,
				Status: "saving",
				Item:   env.Name,
			},
		})
		// ---

		fmt.Println(utils.PrintJsonStep(models.Step{
			Actions: "Getting Environments :: Saving " + env.Name,
			Host:    hostname,
		}))

		ctx.GlobalLock.Lock()
		var codeInfoURL SvnUrlInfo
		codeInfoDIR, errD := RemoteDIRGetSVNInfoName(hostname, env.Name, ctx)
		var errU error
		if errD != nil {
			// Socket Broadcast ---
			ctx.Session.SendMsg(models.WSMessage{
				Broadcast: false,
				Operation: "getEnv",
				Data: models.Step{
					Host:   hostname,
					Status: "error::" + errD.Error(),
					Item:   env.Name,
				},
			})
			// ---
			utils.Warning.Println("no SWE code on host:", env.Name)
		} else {
			r := DbGetRepo(ctx.Config.Hosts[hostname], ctx)

			codeInfoURL, errU = RemoteURLGetSVNInfoName(hostname, env.Name, r, ctx)
			if errU != nil {
				utils.Trace.Println("no SWE code in repo:", env.Name)
			}
		}
		ctx.GlobalLock.Unlock()

		state := "absent"
		if errD == nil && errU == nil {
			state = compareInfo(codeInfoDIR, codeInfoURL)
		}
		repo := DbGetRepo(ctx.Config.Hosts[hostname], ctx)
		if repo == "" {
			repo = "svn://svn.dins.ru/Vportal/trunk/setup/automation/puppet/environments/"
		}
		DbInsert(ctx.Config.Hosts[hostname], env.Name, repo, state, env.ID, codeInfoDIR, ctx)
		afterUpdate = append(afterUpdate, env.Name)

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			Operation: "getEnv",
			Data: models.Step{
				Host:   hostname,
				Status: "done",
				Item:   env.Name,
			},
		})
		// ---
	}
	sort.Strings(afterUpdate)

	// ==========
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments :: Checking outdated",
		Host:    hostname,
	}))

	if aLen != bLen {
		for _, i := range beforeUpdate {
			if !utils.StringInSlice(i, afterUpdate) {
				DbDelete(ctx.Config.Hosts[hostname], i, ctx)
			}
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    hostname,
			Actions: "environment",
			Status:  ctx.Session.UserName,
			State:   "done",
		},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments :: Done",
		Host:    hostname,
	}))

	return nil
}
