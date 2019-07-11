package utils

import (
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

type AllEnvSvn struct {
	Info map[string][]SvnInfo
}

func ParseSvnInfo(stdout string) SvnInfo {
	byNewLine := strings.Split(string(stdout), "\n")
	var res SvnInfo
	for _, i := range byNewLine {
		if strings.Contains(i, ":") {
			j := strings.Split(string(i), ":")

			switch strings.Trim(j[0], " ") {
			case "Path":
				res.Path = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Working Copy Root Path":
				res.WorkingCopyRootPath = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "URL":
				res.URL = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Relative URL":
				res.RelativeURL = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Repository Root":
				res.RepoRoot = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Repository UUID":
				res.RepoUUID = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Revision":
				res.Revision = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Node Kind":
				res.NodeKind = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Schedule":
				res.Schedule = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Last Changed Author":
				res.LastAuthor = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Last Changed Rev":
				res.LastRev = strings.Trim(strings.Join(j[1:], ":"), " ")
			case "Last Changed Date":
				res.LastDate = strings.Trim(strings.Join(j[1:], ":"), " ")
			}
		}
	}
	return res
}
