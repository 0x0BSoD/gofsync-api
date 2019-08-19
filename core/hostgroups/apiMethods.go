package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

// ===============================
// CHECKS
// ===============================
func HostGroupCheck(host string, hostGroupName string, ctx *user.GlobalCTX) HgError {

	var r HostGroups

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	body, _ := logger.ForemanAPI("GET", host, uri, "", ctx)
	err := json.Unmarshal(body.Body, &r)
	if err != nil {
		logger.Warning.Printf("%q, hostGroupJson", err)
	}
	if body.StatusCode == 200 && len(r.Results) > 0 {
		return HgError{
			ID:        r.Results[0].ID,
			HostGroup: hostGroupName,
			Host:      host,
			Error:     "found",
		}
	} else if body.StatusCode == 404 {
		return HgError{
			ID:        -1,
			HostGroup: hostGroupName,
			Host:      host,
			Error:     "not found",
		}
	} else {
		return HgError{
			ID:        -1,
			HostGroup: hostGroupName,
			Host:      host,
			Error:     fmt.Sprintf("error %d", body.StatusCode),
		}
	}

}

// ===============================
// GET
// ===============================
// Just get HostGroup info by name
func HostGroupJson(host string, hostGroupName string, ctx *user.GlobalCTX) (HGElem, HgError) {

	var r HostGroups

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	body, err := logger.ForemanAPI("GET", host, uri, "", ctx)
	if err == nil {
		err := json.Unmarshal(body.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q, hostGroupJson", err)
		}

		resPc := make(map[string][]puppetclass.PuppetClassesWeb)
		puppetClass := puppetclass.ApiByHGJson(host, r.Results[0].ID, ctx)
		for pcName, subClasses := range puppetClass {
			for _, subClass := range subClasses {
				scData := smartclass.SCByPCJson(host, subClass.ID, ctx)
				var scp []smartclass.SmartClass
				var overrides []smartclass.SCOParams
				for _, i := range scData {
					if !StringInMap(i.Parameter, scp) {
						scp = append(scp, smartclass.SmartClass{
							Id:        -1,
							ForemanId: i.ID,
							Name:      i.Parameter,
						})
						if i.OverrideValuesCount > 0 {
							sco := smartclass.SCOverridesById(host, i.ID, ctx)
							for _, j := range sco {
								match := fmt.Sprintf("hostgroup=SWE/%s", r.Results[0].Name)
								if j.Match == match {
									jsonVal, _ := json.Marshal(j.Value)
									overrides = append(overrides, smartclass.SCOParams{
										Match:     j.Match,
										Value:     string(jsonVal),
										Parameter: i.Parameter,
									})
								}
							}
						}
					}
				}
				resPc[pcName] = append(resPc[pcName], puppetclass.PuppetClassesWeb{
					Subclass:     subClass.Name,
					SmartClasses: scp,
					Overrides:    overrides,
				})
			}
		}
		dbId := r.Results[0].ID
		tmpDbId := ID(r.Results[0].Name, host, ctx)
		if tmpDbId != -1 {
			dbId = tmpDbId
		}

		if len(r.Results) > 0 {

			base := HGElem{
				ID:            dbId,
				ForemanID:     r.Results[0].ID,
				Name:          r.Results[0].Name,
				Environment:   r.Results[0].EnvironmentName,
				ParentId:      r.Results[0].Ancestry,
				PuppetClasses: resPc,
			}

			return base, HgError{}
		}
	}
	return HGElem{}, HgError{
		HostGroup: hostGroupName,
		Host:      host,
		Error:     "not found",
	}
}

func StringInMap(a string, list []smartclass.SmartClass) bool {
	for _, b := range list {
		if b.Name == a {
			return true
		}
	}
	return false
}

func GetFromRT(name string, swes map[string]string) string {
	if val, ok := swes[name]; ok {
		return val
	}
	return "nope"
}

// ===================================
// Get SWE from Foreman
func GetHostGroups(host string, ctx *user.GlobalCTX) []HostGroupForeman {
	var r HostGroups
	uri := fmt.Sprintf("hostgroups?format=json&per_page=%d&search=label+~+SWE", ctx.Config.Api.GetPerPage)
	body, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err == nil {
		err = json.Unmarshal(body.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q:\n %s\n", err, body.Body)
		}

		var resultsContainer []HostGroupForeman

		if r.Total > ctx.Config.Api.GetPerPage {
			pagesRange := utils.Pager(r.Total, ctx.Config.Api.GetPerPage)
			for i := 1; i <= pagesRange; i++ {
				uri := fmt.Sprintf("hostgroups?format=json&page=%d&per_page=%d&search=label+~+SWE", i, ctx.Config.Api.GetPerPage)
				body, err := utils.ForemanAPI("GET", host, uri, "", ctx)
				if err == nil {
					err = json.Unmarshal(body.Body, &r)
					if err != nil {
						logger.Warning.Printf("%q:\n %s\n", err, body.Body)
					}
					resultsContainer = append(resultsContainer, r.Results...)
				}
			}
		} else {
			resultsContainer = append(resultsContainer, r.Results...)
		}

		sort.Slice(resultsContainer, func(i, j int) bool {
			return resultsContainer[i].ID < resultsContainer[j].ID
		})

		return resultsContainer
	} else {
		logger.Error.Printf("Error on getting HG, %s", err)
		return []HostGroupForeman{}
	}
}

// Get SWE Parameters from Foreman
func HgParams(host string, dbID int, sweID int, ctx *user.GlobalCTX) {
	var r HostGroupPContainer
	uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&per_page=%d", sweID, ctx.Config.Api.GetPerPage)
	body, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err == nil {
		err = json.Unmarshal(body.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q:\n %s\n", err, body.Body)
		}

		if r.Total > ctx.Config.Api.GetPerPage {
			pagesRange := utils.Pager(r.Total, ctx.Config.Api.GetPerPage)
			for i := 1; i <= pagesRange; i++ {

				uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&page=%d&per_page=%d", sweID, i, ctx.Config.Api.GetPerPage)
				body, err := utils.ForemanAPI("GET", host, uri, "", ctx)
				if err == nil {
					err = json.Unmarshal(body.Body, &r)
					if err != nil {
						logger.Error.Printf("%q:\n %s\n", err, body.Body)
					}
					for _, j := range r.Results {
						InsertParameters(dbID, j, ctx)
					}
				}
			}
		} else {
			for _, i := range r.Results {
				InsertParameters(dbID, i, ctx)
			}
		}
	} else {
		logger.Error.Printf("Error on getting HG Params, %s", err)
	}
}

// Dump HostGroup info by name
func HostGroup(host string, hostGroupName string, ctx *user.GlobalCTX) int {
	var r HostGroups
	lastId := -1

	// Socket Broadcast ---
	data := models.Step{
		Host:    host,
		Actions: "Getting host group from Foreman",
	}
	msg, _ := json.Marshal(data)
	ctx.Session.SendMsg(msg)
	// ---

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	response, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err == nil {
		err := json.Unmarshal(response.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q:\n %s\n", err, response.Body)
		}
		swes := RTBuildObj(PuppetHostEnv(host, ctx), ctx)
		for _, i := range r.Results {

			sJson, _ := json.Marshal(i)

			// Socket Broadcast ---
			data := models.Step{
				Host:    host,
				Actions: "Saving host group",
			}
			msg, _ := json.Marshal(data)
			ctx.Session.SendMsg(msg)
			// ---

			sweStatus := GetFromRT(i.Name, swes)

			lastId = Insert(i.Name, host, string(sJson), sweStatus, i.ID, ctx)

			// Socket Broadcast ---
			data = models.Step{
				Host:    host,
				Actions: "Getting Puppet Classes from Foreman",
			}
			msg, _ = json.Marshal(data)
			ctx.Session.SendMsg(msg)
			// ---

			scpIds := puppetclass.ApiByHG(host, i.ID, lastId, ctx)

			// Socket Broadcast ---
			data = models.Step{
				Host:    host,
				Actions: "Getting Host group parameters from Foreman",
			}
			msg, _ = json.Marshal(data)
			ctx.Session.SendMsg(msg)
			// ---

			HgParams(host, lastId, i.ID, ctx)

			for _, scp := range scpIds {
				scpData := smartclass.SCByPCJsonV2(host, scp, ctx)

				// Socket Broadcast ---
				data := models.Step{
					Host:    host,
					Actions: "Getting Smart classes from Foreman",
					State:   scpData.Name,
				}
				msg, _ := json.Marshal(data)
				ctx.Session.SendMsg(msg)
				// ---

				for _, scParam := range scpData.SmartClassParameters {

					// Socket Broadcast ---
					data := models.Step{
						Host:    host,
						Actions: "Getting Smart class parameters from Foreman",
						State:   scParam.Parameter,
					}
					msg, _ := json.Marshal(data)
					ctx.Session.SendMsg(msg)
					// ---

					scpSummary := smartclass.SCByFId(host, scParam.ID, ctx)
					smartclass.InsertSC(host, scpSummary, ctx)
				}
			}
		}
	} else {
		logger.Error.Printf("Error on getting HG, %s", err)
	}

	//// Socket Broadcast ---
	//data := models.Step{
	//	Host:    host,
	//	Actions: "Update done.",
	//}
	//msg, _ := json.Marshal(data)
	//ctx.Session.SendMsg(msg)
	//// ---

	return lastId
}

//func DeleteHG(host string, hgId int, ctx *user.GlobalCTX) error {
//	data := GetHG(hgId, ctx)
//	uri := fmt.Sprintf("hostgroups/%d", data.ForemanID)
//	resp, err := logger.ForemanAPI("DELETE", host, uri, "", ctx)
//	logger.Trace.Printf("Response on DELETE HG: %q", resp)
//
//	if err != nil {
//		logger.Error.Printf("Error on DELETE HG: %s, uri: %s", err, uri)
//		return err
//	} else {
//		DeleteHGbyId(hgId, ctx)
//	}
//	return nil
//}
