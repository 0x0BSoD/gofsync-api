package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var SVNDIR = "/home/asimonov/Projects/DINS/puppet"

type SvnStatus struct {
	State       string `json:"state"`
	Path        string `json:"path"`
	Environment string `json:"environment"`
}

type SvnInfo struct {
	Path                string `json:"path"`
	WorkingCopyRootPath string `json:"working_copy_root_path"`
	URL                 string `json:"url"`
	RelativeURL         string `json:"relative_url"`
	RepoRoot            string `json:"repo_root"`
	RepoUUID            string `json:"repo_uuid"`
	Revision            string `json:"revision"`
	NodeKind            string `json:"node_kind"`
	Schedule            string `json:"schedule"`
	LastAuthor          string `json:"last_author"`
	LastRev             string `json:"last_rev"`
	LastDate            string `json:"last_date"`
}

func GetInfo() {
	err := os.Chdir(SVNDIR)
	if err != nil {
		Error.Println(err)
	}

	cmd := exec.Command("svn", "info")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	byNewLine := strings.Split(string(stdoutStderr), "\n")
	var res SvnInfo
	for _, i := range byNewLine {
		if strings.Contains(i, ":") {
			j := strings.Split(string(i), ":")

			switch strings.Trim(j[0], " ") {
			case "Path":
				res.Path = strings.Trim(j[1], " ")
			case "Working Copy Root Path":
				res.WorkingCopyRootPath = strings.Trim(j[1], " ")
			case "URL":
				res.URL = strings.Trim(j[1], " ")
			case "Relative URL":
				res.RelativeURL = strings.Trim(j[1], " ")
			case "Repository Root":
				res.RepoRoot = strings.Trim(j[1], " ")
			case "Repository UUID":
				res.RepoUUID = strings.Trim(j[1], " ")
			case "Revision":
				res.Revision = strings.Trim(j[1], " ")
			case "Node Kind":
				res.NodeKind = strings.Trim(j[1], " ")
			case "Schedule":
				res.Schedule = strings.Trim(j[1], " ")
			case "Last Changed Author":
				res.LastAuthor = strings.Trim(j[1], " ")
			case "Last Changed Rev":
				res.LastRev = strings.Trim(j[1], " ")
			case "Last Changed Date":
				res.LastDate = strings.Trim(j[1], " ")
			}
		}
	}
	msg, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(msg))

}
