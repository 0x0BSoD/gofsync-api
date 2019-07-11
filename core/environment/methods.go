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

	beforeUpdate := DbByHost(host, ctx)
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

		codeInfoDIR, err := RemoteDIRGetSVNInfoName(host, env.Name, ctx)
		if err != nil {
			logger.Warning.Println("no SWE code on host:", env.Name)
		}

		//codeInfoURL, err := RemoteURLGetSVNInfoName(host, env.Name, env.Repo, ctx)
		//if err != nil {
		//	logger.Warning.Println("no SWE code on host:", env.Name)
		//}
		// TODO: codeInfoURL

		DbInsert(host, env.Name, env.ID, codeInfoDIR, ctx)
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
	envs := DbByHost(host, ctx)
	for _, env := range envs {
		if strings.HasPrefix(env, "swe") {
			cmd := utils.CmdSvnDirInfo(env)
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

func RemoteDIRGetSVNInfoName(host, name string, ctx *user.GlobalCTX) (utils.SvnInfo, error) {
	var res utils.SvnInfo
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnDirInfo(name)
		var tmpRes []string
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
			return utils.SvnInfo{}, err
		}
		dataSplit := strings.Split(data, "\n")
		for _, s := range dataSplit {
			if s != "" {
				if s == "NIL" {
					logger.Warning.Println("no SWE code on host:", name)
					return utils.SvnInfo{}, utils.NewError("no SWE code on host: " + name)
				} else {
					tmpRes = append(tmpRes, s)
				}
			} else {
				continue
			}
		}

		if len(tmpRes) > 0 {
			joined := strings.Join(tmpRes, "\n")
			res = utils.ParseSvnInfo(joined)
		}
	}
	return res, nil
}

func RemoteURLGetSVNInfoName(host, name, url string, ctx *user.GlobalCTX) (utils.SvnInfo, error) {
	var res utils.SvnInfo
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnUrlInfo(name, url+name)
		var tmpRes []string
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
			return utils.SvnInfo{}, err
		}
		dataSplit := strings.Split(data, "\n")
		for _, s := range dataSplit {
			if s != "" {
				if s == "NIL" {
					logger.Warning.Println("no SWE code on host:", name)
					return utils.SvnInfo{}, utils.NewError("no SWE code on host: " + name)
				} else {
					tmpRes = append(tmpRes, s)
				}
			} else {
				continue
			}
		}

		if len(tmpRes) > 0 {
			joined := strings.Join(tmpRes, "\n")
			res = utils.ParseSvnInfo(joined)
		}
	}
	return res, nil
}

func RemoteGetSVNInfo(ctx *user.GlobalCTX) utils.AllEnvSvn {
	res := utils.AllEnvSvn{
		Info: make(map[string][]utils.SvnInfo),
	}
	for _, host := range ctx.Config.Hosts {
		envs := DbByHost(host, ctx)
		for _, env := range envs {
			if strings.HasPrefix(env, "swe") {
				cmd := utils.CmdSvnDirInfo(env)
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

func RemoteGetSVNDiff(host, name string, ctx *user.GlobalCTX) {
	//var res utils.SvnInfo
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnDiff(name)
		//var tmpRes []string
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
		}
		fmt.Println(data)
		//dataSplit := strings.Split(data, "\n")
		//for _, s := range dataSplit {
		//	if s != "" {
		//		if s == "NIL" {
		//			logger.Warning.Println("no SWE code on host:", name)
		//			return utils.SvnInfo{}, utils.NewError("no SWE code on host: " + name)
		//		} else {
		//			tmpRes = append(tmpRes, s)
		//		}
		//	} else {
		//		continue
		//	}
		//}
		//
		//if len(tmpRes) > 0 {
		//	joined := strings.Join(tmpRes, "\n")
		//	res = utils.ParseSvnInfo(joined)
		//}
	}
	//return res, nil
}
