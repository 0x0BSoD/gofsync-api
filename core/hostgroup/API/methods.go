package API

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/hostgroup/DB"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Get SWE from Foreman
func (Get) All(host string, ctx *user.GlobalCTX) []DB.HostGroup {

	// VARS
	var r HostGroups

	// ========
	uri := fmt.Sprintf("hostgroups?format=json&per_page=%d&search=label+~+SWE", ctx.Config.Api.GetPerPage)
	body, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err != nil {
		utils.Error.Printf("Error on getting HG, %s", err)
		return []DB.HostGroup{}
	}
	err = json.Unmarshal(body.Body, &r)
	if err != nil {
		utils.Warning.Printf("%q:\n %s\n", err, body.Body)
	}

	var result []DB.HostGroup

	if r.Total > ctx.Config.Api.GetPerPage {
		pagesRange := utils.Pager(r.Total, ctx.Config.Api.GetPerPage)
		for i := 1; i <= pagesRange; i++ {
			uri := fmt.Sprintf("hostgroups?format=json&page=%d&per_page=%d&search=label+~+SWE", i, ctx.Config.Api.GetPerPage)
			body, err := utils.ForemanAPI("GET", host, uri, "", ctx)
			if err == nil {
				err = json.Unmarshal(body.Body, &r)
				if err != nil {
					utils.Warning.Printf("%q:\n %s\n", err, body.Body)
				}
				result = append(result, r.Results...)
			}
		}
	} else {
		result = append(result, r.Results...)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ForemanID < result[j].ForemanID
	})

	return result
}

// Get SWE Parameters from Foreman
func (Get) Parameters(host string, dbID int, hgID int, ctx *user.GlobalCTX) []DB.HostGroupParameter {

	// VARS
	var response Parameters
	var result []DB.HostGroupParameter

	// =======
	uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&per_page=%d", hgID, ctx.Config.Api.GetPerPage)
	body, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err == nil {
		utils.Error.Printf("Error on getting HG Params, %s", err)
	}
	err = json.Unmarshal(body.Body, &response)
	if err != nil {
		utils.Warning.Printf("%q:\n %s\n", err, body.Body)
	}

	if response.Total > ctx.Config.Api.GetPerPage {
		pagesRange := utils.Pager(response.Total, ctx.Config.Api.GetPerPage)
		for i := 1; i <= pagesRange; i++ {

			uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&page=%d&per_page=%d", hgID, i, ctx.Config.Api.GetPerPage)
			body, err := utils.ForemanAPI("GET", host, uri, "", ctx)
			if err == nil {
				err = json.Unmarshal(body.Body, &response)
				if err != nil {
					utils.Error.Printf("%q:\n %s\n", err, body.Body)
				}
				for _, j := range response.Results {
					result = append(result, j)
				}
			}
		}
	} else {
		for _, i := range response.Results {
			result = append(result, i)
		}
	}

	return result
}
