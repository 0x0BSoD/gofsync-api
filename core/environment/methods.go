package environment

import (
	"encoding/json"
	"encoding/xml"
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

		r := DbGetRepo(host, ctx)

		codeInfoURL, err := RemoteURLGetSVNInfoName(host, env.Name, r, ctx)
		if err != nil {
			logger.Warning.Println("no SWE code on host:", env.Name)
		}

		state := compareInfo(codeInfoDIR, codeInfoURL)

		DbInsert(host, env.Name, state, env.ID, codeInfoDIR, ctx)
		afterUpdate = append(afterUpdate, env.Name)
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i, afterUpdate) {
			DbDelete(host, i, ctx)
		}
	}
}

func compareInfo(dir, svn SvnInfo) string {
	var state string
	if dir == (SvnInfo{}) {
		state = "absent"
	} else {
		if dir.Entry.Commit.Revision != svn.Entry.Commit.Revision {
			state = "outdated"
		} else {
			state = "ok"
		}
	}
	return state
}
func RemoteGetSVNInfoHost(host string, ctx *user.GlobalCTX) []SvnInfo {
	var res []SvnInfo
	envs := DbByHost(host, ctx)
	for _, env := range envs {
		if strings.HasPrefix(env, "swe") {
			var info SvnInfo
			cmd := utils.CmdSvnDirInfo(env)
			data, err := utils.CallCMDs(host, cmd)
			if err != nil {
				logger.Error.Println(err)
			}

			err = xml.Unmarshal([]byte(data), &info)
			if err != nil {
				logger.Error.Println(err)
				return []SvnInfo{}
			}

			res = append(res, info)
		}
	}
	return res
}

func RemoteGetSVNLog(host, name, url string, ctx *user.GlobalCTX) SvnLog {
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnLog(url + name)
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
			return SvnLog{}
		}
		var logs SvnLog
		err = xml.Unmarshal([]byte(data), &logs)
		if err != nil {
			logger.Error.Println(err)
			return SvnLog{}
		}
		return logs
	}
	return SvnLog{}
}

func RemoteSVNUpdate(host, name string, ctx *user.GlobalCTX) {
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnUpdate(name)
		fmt.Println(cmd)
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
		}
		fmt.Println(data)
		DbSetUpdated("ok", host, name, ctx)
	}
}

func RemoteSVNCheckout(host, name, url string, ctx *user.GlobalCTX) {
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnCheckout(url + name)
		fmt.Println(cmd)
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
		}
		DbSetUpdated("ok", host, name, ctx)
		fmt.Println(data)
	}
}

func RemoteDIRGetSVNInfoName(host, name string, ctx *user.GlobalCTX) (SvnInfo, error) {
	var info SvnInfo
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnDirInfo(name)
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
			return SvnInfo{}, err
		}

		err = xml.Unmarshal([]byte(data), &info)
		if err != nil {
			logger.Error.Println(err)
			return SvnInfo{}, err
		}

	}
	return info, nil
}

func RemoteURLGetSVNInfoName(host, name, url string, ctx *user.GlobalCTX) (SvnInfo, error) {
	var info SvnInfo
	envExist := DbID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnUrlInfo(url + name)
		data, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
			return SvnInfo{}, err
		}

		err = xml.Unmarshal([]byte(data), &info)
		if err != nil {
			logger.Error.Println(err)
			return SvnInfo{}, err
		}

	}

	return info, nil
}

type AllEnvSvn struct {
	Info map[string][]SvnInfo `json:"info"`
}

func RemoteGetSVNInfo(ctx *user.GlobalCTX) AllEnvSvn {
	res := AllEnvSvn{
		Info: make(map[string][]SvnInfo),
	}
	for _, host := range ctx.Config.Hosts {
		envs := DbByHost(host, ctx)
		for _, env := range envs {
			if strings.HasPrefix(env, "swe") {
				var info SvnInfo
				cmd := utils.CmdSvnDirInfo(env)
				data, err := utils.CallCMDs(host, cmd)
				if err != nil {
					logger.Error.Println(err)
				}

				err = xml.Unmarshal([]byte(data), &info)
				if err != nil {
					logger.Error.Println(err)
					return AllEnvSvn{}
				}
				res.Info[host] = append(res.Info[host], info)
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
