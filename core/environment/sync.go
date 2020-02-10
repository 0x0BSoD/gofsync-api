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
		Broadcast:      false,
		HostName:       hostname,
		Resource:       models.Environment,
		Operation:      "sync",
		UserName:       ctx.Session.UserName,
		AdditionalData: models.CommonOperation{Message: "Getting Environments from foreman"},
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
	count := 1
	for _, env := range environmentsResult.Results {
		failed := false

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			HostName:  hostname,
			Resource:  models.Environment,
			Operation: "sync",
			UserName:  ctx.Session.UserName,
			AdditionalData: models.CommonOperation{
				Message: "Saving Environment",
				Item:    env.Name,
				Total:   aLen,
				Current: count,
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
			failed = true
			// Socket Broadcast ---
			ctx.Session.SendMsg(models.WSMessage{
				Broadcast: false,
				HostName:  hostname,
				Resource:  models.Environment,
				Operation: "sync",
				UserName:  ctx.Session.UserName,
				AdditionalData: models.CommonOperation{
					Message: "Saving Environment Failed",
					Failed:  true,
					Item:    env.Name,
					Total:   aLen,
					Current: count,
				},
			})
			// ---
			utils.Warning.Println("no SWE code on host:", env.Name)
		} else {
			r := DbGetRepo(ctx.Config.Hosts[hostname], ctx)

			codeInfoURL, errU = RemoteURLGetSVNInfoName(hostname, env.Name, r, ctx)
			//if errU != nil {
			//	utils.Trace.WithFields(logrus.Fields{"name": env.Name}).Trace("no SWE code in repo")
			//}
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

		if !failed {
			// Socket Broadcast ---
			ctx.Session.SendMsg(models.WSMessage{
				Broadcast: false,
				HostName:  hostname,
				Resource:  models.Environment,
				Operation: "sync",
				UserName:  ctx.Session.UserName,
				AdditionalData: models.CommonOperation{
					Message: "Environment Saved",
					Done:    true,
					Item:    env.Name,
					Total:   aLen,
					Current: count,
				},
			})
			// ---
		}

		count++
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

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments :: Done",
		Host:    hostname,
	}))

	return nil
}
