package environment

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
	"strings"
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

func RemoteGetSVNInfoHost(host string, ctx *user.GlobalCTX) []utils.SvnInfo {
	var res []utils.SvnInfo
	envs := DbAll(host, ctx)
	for _, env := range envs {
		if strings.HasPrefix(env, "swe") {
			cmd := utils.CmdSvnInfo(env)
			var tmpRes []string
			data, err := utils.CallCMDs(host, cmd)
			if err != nil {
				logger.Error.Println(err)
			}
			dataSplit := strings.Split(data, "\n")
			for _, s := range dataSplit {
				if s != "" {
					if s == "NIL" {
						logger.Warning.Println("no SWE code on host:", env)
					} else {
						tmpRes = append(tmpRes, s)
					}
				} else {
					continue
				}
			}

			if len(tmpRes) > 0 {
				joined := strings.Join(tmpRes, "\n")
				res = append(res, utils.ParseSvnInfo(joined))
			}
		}
	}
	return res
}

func RemoteGetSVNInfoName(host, name string, ctx *user.GlobalCTX) []utils.SvnInfo {
	var res []utils.SvnInfo
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnInfo(name)
		var tmpRes []string
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
		}
		dataSplit := strings.Split(data, "\n")
		for _, s := range dataSplit {
			if s != "" {
				if s == "NIL" {
					logger.Warning.Println("no SWE code on host:", name)
				} else {
					tmpRes = append(tmpRes, s)
				}
			} else {
				continue
			}
		}

		if len(tmpRes) > 0 {
			joined := strings.Join(tmpRes, "\n")
			res = append(res, utils.ParseSvnInfo(joined))
		}
	}
	return res
}

func RemoteGetSVNInfo(ctx *user.GlobalCTX) utils.AllEnvSvn {
	res := utils.AllEnvSvn{
		Info: make(map[string][]utils.SvnInfo),
	}
	for _, host := range ctx.Config.Hosts {
		envs := DbAll(host, ctx)
		for _, env := range envs {
			if strings.HasPrefix(env, "swe") {
				cmd := utils.CmdSvnInfo(env)
				var tmpRes []string
				data, err := utils.CallCMDs(host, cmd)
				if err != nil {
					logger.Error.Println(err)
				}
				dataSplit := strings.Split(data, "\n")
				for _, s := range dataSplit {
					if s != "" {
						if s == "NIL" {
							logger.Warning.Println("no SWE code on host:", env)
						} else {
							tmpRes = append(tmpRes, s)
						}
					} else {
						continue
					}
				}

				if len(tmpRes) > 0 {
					joined := strings.Join(tmpRes, "\n")
					res.Info[host] = append(res.Info[host], utils.ParseSvnInfo(joined))
				}
			}
		}
	}
	return res
}
