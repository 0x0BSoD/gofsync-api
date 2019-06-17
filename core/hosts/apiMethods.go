package hosts

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"net/url"
	"strings"
	"sync"
)

// ===============================
// CHECKS
// ===============================

// ===============================
// GET
// ===============================
func ByHostgroup(host string, hg string, cfg *models.Config) models.Response {
	uri := fmt.Sprintf("hostgroups/%s/hosts?format=json&per_page=%d", hg, cfg.Api.GetPerPage)
	response, err := utils.ForemanAPI("GET", host, uri, "", cfg)
	if err != nil {
		logger.Error.Println(err)
	}
	return response
}

func ByHostgroupName(hgName string, params url.Values, cfg *models.Config) map[string][]models.Host {
	result := make(map[string][]models.Host)
	for _, host := range cfg.Hosts {
		id := hostgroups.CheckHGID(hgName, host, cfg)
		uri := "hostgroups/%d/hosts?format=json&per_page=%d"
		if id != -1 {
			if val, ok := params["changed"]; ok {
				p := strings.Trim(val[0], " ")
				if strings.HasPrefix(p, "<") || strings.HasPrefix(p, ">") {
					uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report+%s", id, cfg.Api.GetPerPage, url.QueryEscape(p))
				} else {
					uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report%%3D+%s", id, cfg.Api.GetPerPage, url.QueryEscape(p))
				}
				fmt.Println(uri)
			} else {
				uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d", id, cfg.Api.GetPerPage)
			}
			response, err := utils.ForemanAPI("GET", host, uri, "", cfg)
			if err != nil {
				logger.Error.Println(err)
			}
			if response.StatusCode == 404 {
				logger.Error.Println("not found")
			}

			var tmpResult models.Hosts
			err = json.Unmarshal(response.Body, &tmpResult)
			if err != nil {
				logger.Error.Println(err)
			}

			for _, i := range tmpResult.Results {
				result[host] = append(result[host], i)
			}
		}
	}
	return result
}

// Result struct
type HResult struct {
	sync.Mutex
	hosts map[string][]string
}

func (r *HResult) Add(foreman string, hostname string) {
	r.Lock()
	r.hosts[foreman] = append(r.hosts[foreman], hostname)
	r.Unlock()
}
func (r *HResult) Init() {
	r.hosts = make(map[string][]string)
}
func ByHostgroupNameHostNames(hgName string, params url.Values, cfg *models.Config) map[string][]string {
	var r HResult
	r.Init()
	// Create a new WorkQueue.
	wq := utils.New()
	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	var wg sync.WaitGroup
	for _, host := range cfg.Hosts {
		wg.Add(1)
		go func(h string) {
			wq <- func() {
				defer wg.Done()
				id := hostgroups.CheckHGID(hgName, h, cfg)
				uri := "hostgroups/%d/hosts?format=json&per_page=%d"
				if id != -1 {
					if val, ok := params["changed"]; ok {
						p := strings.Trim(val[0], " ")
						if strings.HasPrefix(p, "<") || strings.HasPrefix(p, ">") {
							uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report+%s", id, cfg.Api.GetPerPage, url.QueryEscape(p))
						} else {
							uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report%%3D+%s", id, cfg.Api.GetPerPage, url.QueryEscape(p))
						}
					} else {
						uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d", id, cfg.Api.GetPerPage)
					}
					response, err := utils.ForemanAPI("GET", h, uri, "", cfg)
					if err != nil {
						logger.Error.Println(err)
						utils.GetErrorContext(err)
					}
					if response.StatusCode == 404 {
						logger.Error.Println("not found")
					}

					var tmpResult models.Hosts
					err = json.Unmarshal(response.Body, &tmpResult)
					if err != nil {
						logger.Error.Println(err)
						utils.GetErrorContext(err)
					}

					for _, i := range tmpResult.Results {
						r.Add(h, i.Name)
					}
				}
			}
		}(host)
	}

	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)
	return r.hosts
}
