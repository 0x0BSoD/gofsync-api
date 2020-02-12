package environment

import (
	"bytes"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/0x0bsod/CmdPusher"
	"github.com/0x0bsod/goLittleHelpers"
	"strconv"
	"strings"
	"time"
)

func RemoteGetGITInfo(hostname, name, commit string, ctx *user.GlobalCTX) (string, error) {
	if ID(ctx.Config.Hosts[hostname], name, ctx) != -1 {
		command := fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then cd %s && sudo git --no-pager show %s; else echo \"NIL\";  fi'", name, name, commit)
		data, err := cmdRunCommandGit(hostname, []string{command})
		if err != nil {
			return "", err
		}
		return data, nil
	}
	return "", nil
}

func RemoteGetGITEnvInfo(hostname, envName string, ctx *user.GlobalCTX) (map[string]map[string]Commits, error) {

	result := make(map[string]map[string]Commits)

	result[hostname] = make(map[string]Commits)
	if strings.HasPrefix(envName, "swe") {
		command := fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then cd %s && sudo git --no-pager log --decorate=short; else echo \"NIL\";  fi'", envName, envName)
		data, err := cmdRunCommandGit(hostname, []string{command})
		if err != nil {
			utils.Error.Println(err)
			return result, err
		}
		commits := parseLog(data)
		result[hostname][envName] = commits
	}
	return result, nil
}

func RemoteGetGITAllEnvInfo(hostname string, ctx *user.GlobalCTX) (map[string]map[string][]string, error) {

	result := make(map[string]map[string][]string)

	envs := DbGetByHost(ctx.Config.Hosts[hostname], ctx)
	result[hostname] = make(map[string][]string)
	for _, env := range envs {
		if strings.HasPrefix(env, "swe") {
			command := fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then cd %s && sudo git --no-pager log --decorate=short --pretty=oneline; else echo \"NIL\";  fi'", env, env)
			data, err := cmdRunCommandGit(hostname, []string{command})
			if err != nil {
				utils.Error.Println(err)
				return result, err
			}
			commits := strings.Split(data, "\n")
			for _, c := range commits {
				result[hostname][env] = append(result[hostname][env], c)
			}
		}
	}
	return result, nil
}

func RemoteGITPull(hostname, name string, ctx *user.GlobalCTX) (string, error) {
	if ID(ctx.Config.Hosts[hostname], name, ctx) != -1 {

		data, err := cmdRunCommandGit(hostname, []string{
			fmt.Sprintf("cd %s && bash -c 'sudo git pull'", name),
			"bash -c 'sudo chown -R puppet:puppet .'",
			"bash -c 'sudo chmod -R 755 .'",
		})

		if err != nil {
			utils.Error.Println(err)
			if err.Error() != "Already up-to-date" {
				DbSetUpdated(ctx.Config.Hosts[hostname], name, "error", ctx)
			}
			return "", err
		}

		DbSetUpdated(ctx.Config.Hosts[hostname], name, "ok", ctx)

		return data, nil
	} else {
		return "", fmt.Errorf("environment %s not exist", name)
	}
}

func RemoteGITClone(hostname, name string, ctx *user.GlobalCTX) (string, error) {
	envExist := ID(ctx.Config.Hosts[hostname], name, ctx)

	if envExist != -1 {

		data, err := cmdRunCommandGit(hostname, []string{
			fmt.Sprintf("bash -c 'sudo git clone --branch=%s https://i5XLLiXsxj4zpmbZ1ynL@git.ringcentral.com/archops/swe.git %s'", name, name),
			fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
			fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name),
		})

		fmt.Println(data)
		fmt.Println(err)

		if err != nil {
			DbSetUpdated(ctx.Config.Hosts[hostname], name, "error", ctx)
			return "", err
		}

		DbSetUpdated(ctx.Config.Hosts[hostname], name, "ok", ctx)
		return data, nil
	} else {
		return "", fmt.Errorf("environment %s not exist, env not exist: %d", name, envExist)
	}
}

func cmdRunCommandGit(host string, cmds []string) (string, error) {
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
		CurrentDir: "/etc/puppet/environments_git",
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

type Commit struct {
	Author  string
	Date    time.Time
	Comment string
	Ticket  string
	SvnID   string
}
type Commits map[string]Commit

type CommitDetail struct {
	Author  string
	Date    time.Time
	Comment string
	Ticket  string
	SvnID   string
	Diffs   []string
}

func parseDate(strDate string) time.Time {
	//Wed Dec 11 14:50:42 2019 +0000
	_split := strings.Split(strDate, " ")
	_time := strings.Split(_split[3], ":")

	var month time.Month
	switch _split[1] {
	case "Jan":
		month = time.January
	case "Feb":
		month = time.February
	case "Mar":
		month = time.March
	case "Apr":
		month = time.April
	case "May":
		month = time.May
	case "Jun":
		month = time.June
	case "Jul":
		month = time.July
	case "Aug":
		month = time.August
	case "Sep":
		month = time.September
	case "Oct":
		month = time.October
	case "Nov":
		month = time.November
	case "Dec":
		month = time.December

	}

	y, _ := strconv.Atoi(_split[4])
	d, _ := strconv.Atoi(_split[2])
	h, _ := strconv.Atoi(_time[0])
	m, _ := strconv.Atoi(_time[1])
	s, _ := strconv.Atoi(_time[2])

	date := time.Date(y, month, d, h, m, s, 0, time.UTC)

	return date
}

func parseCommit(strData string) CommitDetail {
	diffFlag := false
	var oneDiff []string
	var result CommitDetail
	var author string
	var date time.Time
	var comm string

	for _, l := range strings.Split(strData, "\n") {
		if len(l) > 0 {
			_split := strings.Split(l, " ")
			if strings.HasPrefix(_split[0], "diff") {
				diffFlag = true
				if len(oneDiff) > 0 {
					result.Diffs = append(result.Diffs, strings.Join(oneDiff, "\n"))
				}
				oneDiff = []string{}
				oneDiff = append(oneDiff, strings.Join(_split, " "))
			} else if diffFlag {
				oneDiff = append(oneDiff, strings.Join(_split, " "))
			} else {
				_split := strings.Split(strings.TrimSpace(l), " ")
				if len(_split) >= 2 {
					if strings.HasPrefix(_split[0], "commit") {
						diffFlag = false
					} else if strings.HasPrefix(_split[0], "Author") {
						author = _split[1]
					} else if strings.HasPrefix(_split[0], "Date") {
						date = parseDate(strings.Join(_split[3:], " "))
					} else if strings.HasPrefix(_split[0], "git-svn-id") {
						result.Author = author
						result.Date = date
						result.SvnID = _split[1]
						result.Ticket = strings.Replace(_split[0], ":", "", 1)
						result.Comment = comm
					} else {
						comm = strings.Join(_split[1:], " ")
					}
				}
			}
		}
	}
	// add in case diff only one or last
	result.Diffs = append(result.Diffs, strings.Join(oneDiff, "\n"))
	_ = goLittleHelpers.PrettyPrint(result)
	return result
}

func parseLog(strData string) Commits {

	fmt.Println(strData)

	var result Commits
	var hash string
	var author string
	var date time.Time
	var comm string
	result = make(map[string]Commit)

	for _, l := range strings.Split(strData, "\n") {
		if len(l) > 0 {
			_split := strings.Split(strings.TrimSpace(l), " ")
			if len(_split) >= 2 {
				if strings.HasPrefix(_split[0], "commit") {
					hash = _split[1]
				} else if strings.HasPrefix(_split[0], "Author") {
					author = _split[1]
				} else if strings.HasPrefix(_split[0], "Date") {
					date = parseDate(strings.Join(_split[3:], " "))
				} else if strings.HasPrefix(_split[0], "git-svn-id") {
					result[hash] = Commit{
						Author:  author,
						Date:    date,
						SvnID:   _split[1],
						Ticket:  strings.Replace(_split[0], ":", "", 1),
						Comment: comm,
					}
				} else {
					comm = strings.Join(_split[1:], " ")
				}
			}
		}
	}

	_ = goLittleHelpers.PrettyPrint(result)
	return result
}
