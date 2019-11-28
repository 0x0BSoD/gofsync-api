package environment

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"git.ringcentral.com/archops/CmdPusher"
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
		repo := DbGetRepo(host, ctx)
		DbInsert(host, env.Name, repo, state, env.ID, codeInfoDIR, ctx)
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

func compareInfo(dir SvnDirInfo, url SvnUrlInfo) string {
	var state string
	if dir == (SvnDirInfo{}) {
		state = "absent"
	} else {
		if dir.Entry.Commit.Revision != url.Entry.Commit.Revision {
			state = "outdated"
		} else {
			state = "ok"
		}
	}
	return state
}

// =====================================================================================================================
//	SVN
// =====================================================================================================================
func cmdRunCommand(host string, cmds []string) (string, error) {
	var client = CmdPusher.Client{
		Host:     host,
		Port:     "22",
		User:     "swe_checker",
		AuthKey:  fmt.Sprintf("./ssh_keys/%s_rsa", strings.Split(host, "-")[0]),
		Insecure: true,
	}

	var bOut bytes.Buffer
	var bErr bytes.Buffer

	cmd := &CmdPusher.Cmd{
		Commands:   cmds,
		CurrentDir: "/etc/puppet/environments",
		StdOut:     &bOut,
		StdErr:     &bErr,
	}

	err := client.Connect()
	if err != nil {
		return "", err
	}

	err = client.Run(cmd)
	if err != nil {
		return "", err
	}

	_ = client.Close()
	//if err != nil {
	//	return err
	//}

	outStr := bOut.String()
	errStr := bErr.String()

	fmt.Println("STD ===================")
	fmt.Println(outStr)
	fmt.Println("ERR ===================")
	fmt.Println(errStr)

	return outStr, nil
}

func RemoteGetSVNInfoHost(host string, ctx *user.GlobalCTX) []SvnDirInfo {
	var res []SvnDirInfo

	envs := DbByHost(host, ctx)

	for _, env := range envs {
		if strings.HasPrefix(env, "swe") {
			var info SvnDirInfo

			command := fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn info --xml ./\"%s\"; else echo \"NIL\";  fi'", env, env)
			data, err := cmdRunCommand(host, []string{command})
			if err != nil {
				logger.Error.Println(err)
			}

			err = xml.Unmarshal([]byte(data), &info)
			if err != nil {
				logger.Error.Println(err)
				return []SvnDirInfo{}
			}

			res = append(res, info)
		}
	}
	return res
}

func RemoteGetSVNLog(host, name, url string, ctx *user.GlobalCTX) SvnLog {
	envExist := ID(host, name, ctx)

	if envExist != -1 {

		command := fmt.Sprintf("bash -c 'sudo svn log --xml \"%s\"'", url+name)
		data, err := cmdRunCommand(host, []string{command})
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

		data, err := cmdRunCommand(host, []string{
			fmt.Sprintf("bash -c 'sudo svn update \"%s\"'", name),
			fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
			fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name)})

		fmt.Println(fmt.Sprintf("bash -c 'sudo svn update \"%s\"'", name),
			fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
			fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name))

		if err != nil {
			logger.Error.Println(err)
			DbSetUpdated("error", host, name, ctx)
			return "", fmt.Errorf("error on update: %s", name)
		}

		DbSetUpdated("ok", host, name, ctx)

		return data, nil
	} else {
		return "", fmt.Errorf("environment %s not exist", name)
	}
}

func RemoteSVNCheckout(host, name, url string, ctx *user.GlobalCTX) (string, error) {
	envExist := ID(host, name, ctx)

	if envExist != -1 {

		data, err := cmdRunCommand(host, []string{
			fmt.Sprintf("bash -c 'sudo svn checkout \"%s\"'", url+name),
			fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
			fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name),
		})

		fmt.Println(fmt.Sprintf("bash -c 'sudo svn checkout \"%s\"'", url+name),
			fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
			fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name))

		if err != nil {
			logger.Error.Println(err)
			DbSetUpdated("error", host, name, ctx)
			return "", fmt.Errorf("error on update: %s", name)
		}

		DbSetUpdated("ok", host, name, ctx)
		return data, nil
	} else {
		return "", fmt.Errorf("environment %s not exist, env not exist: %d", name, envExist)
	}
}

func RemoteDIRGetSVNInfoName(host, name string, ctx *user.GlobalCTX) (SvnDirInfo, error) {
	var info SvnDirInfo
	envExist := ID(host, name, ctx)

	if envExist != -1 {

		command := fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn info --xml ./\"%s\"; else echo \"NIL\";  fi'", name, name)
		data, err := cmdRunCommand(host, []string{command})
		if err != nil {
			logger.Error.Println(err)
			return SvnDirInfo{}, err
		}

		err = xml.Unmarshal([]byte(data), &info)
		if err != nil {
			return SvnDirInfo{}, err
		}
	}
	return info, nil
}

func RemoteURLGetSVNInfoName(host, name, url string, ctx *user.GlobalCTX) (SvnUrlInfo, error) {
	var info SvnUrlInfo
	envExist := ID(host, name, ctx)

	if envExist != -1 {

		command := fmt.Sprintf("bash -c 'sudo svn info --xml \"%s\"'", url+name)
		data, err := cmdRunCommand(host, []string{command})
		if err != nil {
			logger.Error.Println(err)
			return SvnUrlInfo{}, err
		}

		//cmd := utils.CmdSvnUrlInfo(url + name)
		//fmt.Println(cmd)
		//response, err := utils.CallCMDs(host, cmd)
		//fmt.Println(response)
		//fmt.Print("\n\n")
		//if err != nil {
		//	return SvnUrlInfo{}, err
		//}

		err = xml.Unmarshal([]byte(data), &info)
		if err != nil {
			return SvnUrlInfo{}, err
		}
	}
	return info, nil
}

func RemoteGetSVNInfo(ctx *user.GlobalCTX) AllEnvSvn {
	res := AllEnvSvn{
		Info: make(map[string][]SvnDirInfo),
	}
	for _, host := range ctx.Config.Hosts {
		envs := DbByHost(host, ctx)
		for _, env := range envs {
			if strings.HasPrefix(env, "swe") {
				var info SvnDirInfo
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
							Host:  host,
							Item:  env,
							State: "checking",
						},
					})
					// ---
					var state string
					codeInfoDIR, err := RemoteDIRGetSVNInfoName(host, env, ctx)
					if err != nil {
						logger.Warning.Println("no SWE code on host:", env)
						state = "error"
						DbSetUpdated(state, host, env, ctx)
					}

					r := DbGetRepo(host, ctx)

					codeInfoURL, err := RemoteURLGetSVNInfoName(host, env, r, ctx)
					if err != nil {
						logger.Warning.Println("no SWE code on host:", env)
						state = "error"
						DbSetUpdated(state, host, env, ctx)
					}

					if state != "error" {
						state = compareInfo(codeInfoDIR, codeInfoURL)
					}

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "svnCheck",
						Data: models.Step{
							Host:  host,
							Item:  env,
							State: state,
						},
					})
					// ---

					if state == "outdated" {
						r, err := RemoteSVNUpdate(host, env, ctx)
						if err != nil {
							logger.Warning.Println("swe update error:", env)
						}
						fmt.Println(r)
					} else if state == "absent" {
						url := DbGetRepo(host, ctx)
						r, err := RemoteSVNCheckout(host, env, url, ctx)
						if err != nil {
							logger.Warning.Println("swe checkout error:", env)
						}
						fmt.Println(r)
					}

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "svnCheck",
						Data: models.Step{
							Host:  host,
							Item:  env,
							State: "done",
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
