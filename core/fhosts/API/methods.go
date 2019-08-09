package API

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"net/url"
	"strings"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Get host names by host group name
func (Get) ByHostGroup(hgName string, params url.Values, ctx *user.GlobalCTX) map[string][]Host {
	result := make(map[string][]Host)
	for _, host := range ctx.Config.Hosts {
		// TODO: >
		//id := hostgroups.CheckHGID(hgName, host, ctx)
		id := -1
		uri := "hostgroups/%d/hosts?format=json&per_page=%d"
		if id != -1 {
			if val, ok := params["changed"]; ok {
				p := strings.Trim(val[0], " ")
				if strings.HasPrefix(p, "<") || strings.HasPrefix(p, ">") {
					uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report+%s", id, ctx.Config.Api.GetPerPage, url.QueryEscape(p))
				} else {
					uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report%%3D+%s", id, ctx.Config.Api.GetPerPage, url.QueryEscape(p))
				}
				fmt.Println(uri)
			} else {
				uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d", id, ctx.Config.Api.GetPerPage)
			}
			response, err := utils.ForemanAPI("GET", host, uri, "", ctx)
			if err != nil {
				utils.Error.Println(err)
			}
			if response.StatusCode == 404 {
				utils.Error.Println("not found")
			}

			var tmpResult Hosts
			err = json.Unmarshal(response.Body, &tmpResult)
			if err != nil {
				utils.Error.Println(err)
			}

			for _, i := range tmpResult.Results {
				result[host] = append(result[host], i)
			}
		}
	}
	return result
}
