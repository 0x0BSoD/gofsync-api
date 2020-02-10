package environment

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/0x0bsod/CmdPusher"
	"strings"
	"sync"
)

func RemoteGetSVNInfoHost(hostname string, ctx *user.GlobalCTX) []SvnDirInfo {
	var res []SvnDirInfo

	envs := DbGetByHost(ctx.Config.Hosts[hostname], ctx)

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
			//utils.Error.Println(err)
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
			//utils.Error.Println(err)
			return SvnDirInfo{}, err
		}

		err = xml.Unmarshal([]byte(data), &info)
		if err != nil {
			//utils.Error.Println(err)
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
		envs := DbGetByHost(ID, ctx)
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
						Resource:  models.Environment,
						HostName:  host,
						Operation: "svnCheck",
						UserName:  ctx.Session.UserName,
						AdditionalData: models.CommonOperation{
							Message: "Checking Environment Code",
							Item:    name,
						},
					})
					// ---

					codeInfoDIR, err := RemoteDIRGetSVNInfoName(host, name, ctx)
					if err != nil {
						utils.Warning.Println("no SWE code on host:", name)
						DbSetUpdated(hostID, name, "error", ctx)

						// Socket Broadcast ---
						ctx.Session.SendMsg(models.WSMessage{
							Broadcast: false,
							Resource:  models.Environment,
							HostName:  host,
							Operation: "svnCheck",
							UserName:  ctx.Session.UserName,
							AdditionalData: models.CommonOperation{
								Message: "Code Status",
								State:   "error",
								Failed:  true,
								Item:    name,
							},
						})
						// ---

					} else {

						r := DbGetRepo(hostID, ctx)
						codeInfoURL, err := RemoteURLGetSVNInfoName(host, name, r, ctx)
						if err != nil {
							utils.Warning.Println("no SWE code on host:", name)
							DbSetUpdated(hostID, name, "error", ctx)

							// Socket Broadcast ---
							ctx.Session.SendMsg(models.WSMessage{
								Broadcast: false,
								Resource:  models.Environment,
								HostName:  host,
								Operation: "svnCheck",
								UserName:  ctx.Session.UserName,
								AdditionalData: models.CommonOperation{
									Message: "Code Status",
									State:   "error",
									Failed:  true,
									Item:    name,
								},
							})
							// ---

						} else {
							state := compareInfo(codeInfoDIR, codeInfoURL)

							// Socket Broadcast ---
							ctx.Session.SendMsg(models.WSMessage{
								Broadcast: false,
								Resource:  models.Environment,
								HostName:  host,
								Operation: "svnCheck",
								UserName:  ctx.Session.UserName,
								AdditionalData: models.CommonOperation{
									Message: "Code Status",
									State:   state,
									Item:    name,
								},
							})
							// ---

							if state == "ok" {
								DbSetUpdated(hostID, name, state, ctx)
							} else if state == "outdated" {

								// Socket Broadcast ---
								ctx.Session.SendMsg(models.WSMessage{
									Broadcast: false,
									Resource:  models.Environment,
									HostName:  host,
									Operation: "svnCheck",
									UserName:  ctx.Session.UserName,
									AdditionalData: models.CommonOperation{
										Message: "Running 'svn up'",
										State:   "svnUpdate",
										Item:    name,
									},
								})
								// ---

								_, err := RemoteSVNUpdate(host, name, ctx)
								if err != nil {
									utils.Warning.Println("swe update error:", name)
									// Socket Broadcast ---
									ctx.Session.SendMsg(models.WSMessage{
										Broadcast: false,
										Resource:  models.Environment,
										HostName:  host,
										Operation: "svnCheck",
										UserName:  ctx.Session.UserName,
										AdditionalData: models.CommonOperation{
											Message: "Running 'svn up' failed",
											State:   err.Error(),
											Failed:  true,
											Item:    name,
										},
									})
									// ---
								}
							} else if state == "absent" {
								// Socket Broadcast ---
								ctx.Session.SendMsg(models.WSMessage{
									Broadcast: false,
									Resource:  models.Environment,
									HostName:  host,
									Operation: "svnCheck",
									UserName:  ctx.Session.UserName,
									AdditionalData: models.CommonOperation{
										Message: "Running 'svn co'",
										State:   "svnCheckout",
										Item:    name,
									},
								})
								// ---
								url := DbGetRepo(hostID, ctx)
								_, err := RemoteSVNCheckout(host, name, url, ctx)
								if err != nil {
									utils.Warning.Println("swe checkout error:", name)
									// Socket Broadcast ---
									ctx.Session.SendMsg(models.WSMessage{
										Broadcast: false,
										Resource:  models.Environment,
										HostName:  host,
										Operation: "svnCheck",
										UserName:  ctx.Session.UserName,
										AdditionalData: models.CommonOperation{
											Message: "Running 'svn co' failed",
											State:   err.Error(),
											Failed:  true,
											Item:    name,
										},
									})
									// ---
								}
							}
						}
					}

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Resource:  models.Environment,
						HostName:  host,
						Operation: "svnCheck",
						UserName:  ctx.Session.UserName,
						AdditionalData: models.CommonOperation{
							Message: "Checking Environment Code done",
							Item:    name,
							Done:    true,
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
}

// =====================================================================================================================
//	HELPERS
// =====================================================================================================================
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
		errStr := bErr.String()
		outStr := bOut.String()
		return "", fmt.Errorf(outStr + "\n" + errStr + "\n" + err.Error())
	}
	_ = client.Close()
	outStr := bOut.String()

	return outStr, nil
}
