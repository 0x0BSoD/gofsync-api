package environment

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/0x0bsod/CmdPusher"
	"sort"
	"strings"
	"sync"
)

func Sync(hostname string, ctx *user.GlobalCTX) {

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

	beforeUpdate := DbByHost(ctx.Config.Hosts[hostname], ctx)

	environmentsResult, err := ApiAll(hostname, ctx)
	if err != nil {
		utils.Warning.Printf("Error on getting Environments:\n%q", err)
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
				Host:   hostname,
				Status: "saving",
				Item:   env.Name,
			},
		})
		// ---

		codeInfoDIR, errD := RemoteDIRGetSVNInfoName(hostname, env.Name, ctx)
		if errD != nil {
			utils.Warning.Println("no SWE code on host:", env.Name)
		}

		r := DbGetRepo(ctx.Config.Hosts[hostname], ctx)

		codeInfoURL, errU := RemoteURLGetSVNInfoName(hostname, env.Name, r, ctx)
		if errU != nil {
			utils.Trace.Println("no SWE code in repo:", env.Name)
		}

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
	}
	sort.Strings(afterUpdate)

	if aLen != bLen {
		for _, i := range beforeUpdate {
			if !utils.StringInSlice(i, afterUpdate) {
				DbDelete(ctx.Config.Hosts[hostname], i, ctx)
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

}

func compareInfo(dir SvnDirInfo, url SvnUrlInfo) string {
	var state string

	fmt.Println("DIR:", dir.Entry.Commit.Revision)
	fmt.Println("URL:", url.Entry.Commit.Revision)

	if dir == (SvnDirInfo{}) || url == (SvnUrlInfo{}) {
		state = "error"
	} else if dir == (SvnDirInfo{}) {
		state = "absent"
	} else if dir.Entry.Commit.Revision != url.Entry.Commit.Revision {
		state = "outdated"
	} else {
		state = "ok"
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
	//errStr := bErr.String()

	//fmt.Printf("%s STD ===================\n", host)
	//fmt.Println(outStr)
	//fmt.Println("ERR ===================")
	//fmt.Println(errStr)

	return outStr, nil
}

func RemoteGetSVNInfoHost(hostname string, ctx *user.GlobalCTX) []SvnDirInfo {
	var res []SvnDirInfo

	envs := DbByHost(ctx.Config.Hosts[hostname], ctx)

	for _, env := range envs {
		if strings.HasPrefix(env, "swe") {
			var info SvnDirInfo

			command := fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn info --xml ./\"%s\"; else echo \"NIL\";  fi'", env, env)
			data, err := cmdRunCommand(hostname, []string{command})
			if err != nil {
				utils.Error.Println(err)
			}

			err = xml.Unmarshal([]byte(data), &info)
			if err != nil {
				utils.Error.Println(err)
				return []SvnDirInfo{}
			}

			res = append(res, info)
		}
	}
	return res
}

func RemoteGetSVNLog(hostname, name, url string, ctx *user.GlobalCTX) SvnLog {
	if ID(ctx.Config.Hosts[hostname], name, ctx) != -1 {

		command := fmt.Sprintf("bash -c 'sudo svn log --xml \"%s\"'", url+name)
		data, err := cmdRunCommand(hostname, []string{command})
		if err != nil {
			utils.Error.Println(err)
			return SvnLog{}
		}

		var logs SvnLog
		err = xml.Unmarshal([]byte(data), &logs)
		if err != nil {
			utils.Error.Println(err)
			return SvnLog{}
		}
		return logs
	}
	return SvnLog{}
}

func RemoteSVNUpdate(hostname, name string, ctx *user.GlobalCTX) (string, error) {
	if ID(ctx.Config.Hosts[hostname], name, ctx) != -1 {

		data, err := cmdRunCommand(hostname, []string{
			fmt.Sprintf("bash -c 'sudo svn update \"%s\"'", name),
			fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
			fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name)})

		if err != nil {
			utils.Error.Println(err)
			DbSetUpdated(ctx.Config.Hosts[hostname], name, "error", ctx)
			return "", fmt.Errorf("error on update: %s", name)
		}

		DbSetUpdated(ctx.Config.Hosts[hostname], name, "ok", ctx)

		return data, nil
	} else {
		return "", fmt.Errorf("environment %s not exist", name)
	}
}

func RemoteSVNCheckout(hostname, name, url string, ctx *user.GlobalCTX) (string, error) {
	envExist := ID(ctx.Config.Hosts[hostname], name, ctx)

	if envExist != -1 {

		data, err := cmdRunCommand(hostname, []string{
			fmt.Sprintf("bash -c 'sudo svn checkout \"%s\"'", url+name),
			fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
			fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name),
		})

		if err != nil {
			utils.Error.Println(err)
			DbSetUpdated(ctx.Config.Hosts[hostname], name, "error", ctx)
			return "", fmt.Errorf("error on update: %s", name)
		}

		DbSetUpdated(ctx.Config.Hosts[hostname], name, "ok", ctx)
		return data, nil
	} else {
		return "", fmt.Errorf("environment %s not exist, env not exist: %d", name, envExist)
	}
}

func RemoteDIRGetSVNInfoName(hostname, name string, ctx *user.GlobalCTX) (SvnDirInfo, error) {
	var info SvnDirInfo

	if ID(ctx.Config.Hosts[hostname], name, ctx) != -1 {
		command := fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn info --xml ./\"%s\"; else echo \"NIL\";  fi'", name, name)
		data, err := cmdRunCommand(hostname, []string{command})
		if err != nil {
			utils.Error.Println(err)
			return SvnDirInfo{}, err
		}

		err = xml.Unmarshal([]byte(data), &info)
		if err != nil {
			return SvnDirInfo{}, err
		}
	}
	return info, nil
}

func RemoteURLGetSVNInfoName(hostname, name, url string, ctx *user.GlobalCTX) (SvnUrlInfo, error) {
	var info SvnUrlInfo

	if ID(ctx.Config.Hosts[hostname], name, ctx) != -1 {
		command := fmt.Sprintf("bash -c 'sudo svn info --xml \"%s\"'", url+name)
		data, err := cmdRunCommand(hostname, []string{command})
		if err != nil {
			utils.Error.Println(err)
			return SvnUrlInfo{}, err
		}

		err = xml.Unmarshal([]byte(data), &info)
		if err != nil {
			return SvnUrlInfo{}, err
		}
	}
	return info, nil
}

func RemoteGetSVNInfo(ctx *user.GlobalCTX) (AllEnvSvn, error) {
	res := AllEnvSvn{
		Info: make(map[string][]SvnDirInfo),
	}
	for hostname, ID := range ctx.Config.Hosts {
		envs := DbByHost(ID, ctx)
		for _, env := range envs {
			if strings.HasPrefix(env, "swe") {
				var info SvnDirInfo

				command := fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn info --xml ./\"%s\"; else echo \"NIL\";  fi'", env, env)
				data, err := cmdRunCommand(hostname, []string{command})
				if err != nil {
					utils.Error.Println(err)
					return AllEnvSvn{}, err
				}

				err = xml.Unmarshal([]byte(data), &info)
				if err != nil {
					utils.Error.Println(err)
					return AllEnvSvn{}, err
				}
				res.Info[hostname] = append(res.Info[hostname], info)
			}
		}
	}
	return res, nil
}

func RemoteSVNBatch(body map[string][]string, ctx *user.GlobalCTX) {
	// Create a new WorkQueue.
	wq := utils.New()
	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	var wg sync.WaitGroup

	for hostname, envs := range body {
		hostID := ctx.Config.Hosts[hostname]
		wg.Add(1)
		go func(envs []string, host string, hostID int) {
			wq <- func() {
				defer wg.Done()
				for _, name := range envs {

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "svnCheck",
						Data: models.Step{
							Host:  host,
							Item:  name,
							State: "checking",
						},
					})
					// ---

					var state string
					codeInfoDIR, err := RemoteDIRGetSVNInfoName(host, name, ctx)
					if err != nil {
						utils.Warning.Println("no SWE code on host:", name)
						state = "error"
						DbSetUpdated(hostID, name, state, ctx)
					}

					r := DbGetRepo(hostID, ctx)

					codeInfoURL, err := RemoteURLGetSVNInfoName(host, name, r, ctx)
					if err != nil {
						utils.Warning.Println("no SWE code on host:", name)
						state = "error"
						DbSetUpdated(hostID, name, state, ctx)
					}

					if state != "error" {
						fmt.Println(host, name)
						state = compareInfo(codeInfoDIR, codeInfoURL)
					}

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "svnCheck",
						Data: models.Step{
							Host:  host,
							Item:  name,
							State: state,
						},
					})
					// ---

					//if state == "outdated" {
					//	r, err := RemoteSVNUpdate(host, name, ctx)
					//	if err != nil {
					//		utils.Warning.Println("swe update error:", name)
					//	}
					//	fmt.Println(r)
					//} else if state == "absent" {
					//	url := DbGetRepo(hostID, ctx)
					//	r, err := RemoteSVNCheckout(host, name, url, ctx)
					//	if err != nil {
					//		utils.Warning.Println("swe checkout error:", name)
					//	}
					//	fmt.Println(r)
					//}

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "svnCheck",
						Data: models.Step{
							Host:  host,
							Item:  name,
							State: "done",
						},
					})
					// ---
				}
			}
		}(envs, hostname, hostID)
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
