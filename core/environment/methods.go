package environment

import (
	"encoding/xml"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
	"strings"
	"sync"
)

func Sync(host string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments :: Start",
		Host:    host,
	}))

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    host,
			Actions: "environment",
			Status:  ctx.Session.UserName,
			State:   "started",
		},
	})

	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "getEnv",
		Data: models.Step{
			Host:  host,
			State: "running",
		},
	})
	// ---

	beforeUpdate := DbByHost(host, ctx)

	environmentsResult, err := ApiAll(host, ctx)
	if err != nil {
		logger.Warning.Printf("Error on getting Environments:\n%q", err)
	}

	sort.Slice(environmentsResult.Results, func(i, j int) bool {
		return environmentsResult.Results[i].ID < environmentsResult.Results[j].ID
	})

	aLen := len(environmentsResult.Results)
	bLen := len(beforeUpdate)

	var afterUpdate = make([]string, 0, aLen)

	for _, env := range environmentsResult.Results {

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			Operation: "getEnv",
			Data: models.Step{
				Host:   host,
				Status: "saving",
				Item:   env.Name,
			},
		})
		// ---

		codeInfoDIR, errD := RemoteDIRGetSVNInfoName(host, env.Name, ctx)
		if errD != nil {
			logger.Warning.Println("no SWE code on host:", env.Name)
		}

		r := DbGetRepo(host, ctx)

		codeInfoURL, errU := RemoteURLGetSVNInfoName(host, env.Name, r, ctx)
		if errU != nil {
			logger.Trace.Println("no SWE code in repo:", env.Name)
		}

		state := "absent"
		if errD == nil && errU == nil {
			state = compareInfo(codeInfoDIR, codeInfoURL)
		}

		DbInsert(host, env.Name, state, env.ID, codeInfoDIR, ctx)
		afterUpdate = append(afterUpdate, env.Name)
	}
	sort.Strings(afterUpdate)

	if aLen != bLen {
		for _, i := range beforeUpdate {
			if !utils.StringInSlice(i, afterUpdate) {
				DbDelete(host, i, ctx)
			}
		}
	}

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "done",
	})
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:    host,
			Actions: "environment",
			Status:  ctx.Session.UserName,
			State:   "done",
		},
	})
	// ---

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Environments :: Done",
		Host:    host,
	}))

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
	envExist := ID(host, name, ctx)
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

func RemoteSVNUpdate(host, name string, ctx *user.GlobalCTX) (string, error) {
	envExist := ID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnUpdate(name)
		out, err := utils.CallCMDs(host, cmd)
		logger.Info.Println(out)
		if err != nil {
			logger.Error.Println(err)
		}
		DbSetUpdated("ok", host, name, ctx)
		return out, nil
	} else {
		return "", fmt.Errorf("environment %s not exist", name)
	}
}

func RemoteSVNCheckout(host, name, url string, ctx *user.GlobalCTX) (string, error) {
	envExist := ID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnCheckout(url+name, name)
		out, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
		}
		DbSetUpdated("ok", host, name, ctx)
		return out, nil
	} else {
		return "", fmt.Errorf("environment %s not exist, env not exist: %d", name, envExist)
	}
}

func RemoteDIRGetSVNInfoName(host, name string, ctx *user.GlobalCTX) (SvnInfo, error) {
	var info SvnInfo
	envExist := ID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnDirInfo(name)

		response, err := utils.CallCMDs(host, cmd)
		if err != nil {
			return SvnInfo{}, err
		}

		err = xml.Unmarshal([]byte(response), &info)
		if err != nil {
			return SvnInfo{}, err
		}
	}
	return info, nil
}

func RemoteURLGetSVNInfoName(host, name, url string, ctx *user.GlobalCTX) (SvnInfo, error) {
	var info SvnInfo
	envExist := ID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnUrlInfo(url + name)

		response, err := utils.CallCMDs(host, cmd)
		if err != nil {
			return SvnInfo{}, err
		}

		err = xml.Unmarshal([]byte(response), &info)
		if err != nil {
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

func RemoteSVNBatch(body map[string][]string, ctx *user.GlobalCTX) {

	// Create a new WorkQueue.
	wq := utils.New()
	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	var wg sync.WaitGroup

	for host, envs := range body {
		wg.Add(1)
		go func(envs []string, host string) {
			wq <- func() {
				defer wg.Done()
				for _, env := range envs {

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "svnCheck",
						Data: models.Step{
							Host:   host,
							Item:   env,
							Status: "checking",
						},
					})
					// ---

					codeInfoDIR, err := RemoteDIRGetSVNInfoName(host, env, ctx)
					if err != nil {
						logger.Warning.Println("no SWE code on host:", env)
					}

					r := DbGetRepo(host, ctx)

					codeInfoURL, err := RemoteURLGetSVNInfoName(host, env, r, ctx)
					if err != nil {
						logger.Warning.Println("no SWE code on host:", env)
					}

					state := compareInfo(codeInfoDIR, codeInfoURL)

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "svnCheck",
						Data: models.Step{
							Host:    host,
							Actions: env,
							State:   state,
						},
					})
					// ---

					//fmt.Println(host, env, state)

					if state == "outdated" {
						r, _ := RemoteSVNUpdate(host, env, ctx)
						fmt.Println(r)
					} else if state == "absent" {
						url := DbGetRepo(host, ctx)
						r, _ := RemoteSVNCheckout(host, env, url, ctx)
						fmt.Println(r)
					}

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "svnCheck",
						Data: models.Step{
							Host:    host,
							Actions: env,
							State:   "done",
						},
					})
				}
			}
		}(envs, host)
	}
	// Wait for all the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "done",
	})
	// ---
}

// TODO: ~ sometime
func RemoteGetSVNDiff(host, name string, ctx *user.GlobalCTX) {
	//var res utils.SvnInfo
	envExist := ID(host, name, ctx)
	if envExist != -1 {
		cmd := utils.CmdSvnDiff(name)
		//var tmpRes []string
		_, err := utils.CallCMDs(host, cmd)
		if err != nil {
			logger.Error.Println(err)
		}
		//fmt.Println(data)
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
