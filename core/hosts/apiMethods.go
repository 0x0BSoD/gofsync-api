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

func ByHostgroupNameHostNames(hgName string, params url.Values, cfg *models.Config) map[string][]string {
	result := make(map[string][]string)
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
				result[host] = append(result[host], i.Name)
			}
		}
	}
	return result
}
