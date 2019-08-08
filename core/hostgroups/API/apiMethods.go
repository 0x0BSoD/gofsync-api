package API

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

// ===============================
// GET
// ===============================

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
func HostGroupPPPP(host string, hostGroupName string, ctx *user.GlobalCTX) int {
	var r HostGroups
	lastId := -1

	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    host,
			Actions: "Getting host group from Foreman",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
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
			if ctx.Session.PumpStarted {
				data := models.Step{
					Host:    host,
					Actions: "Saving host group",
				}
				msg, _ := json.Marshal(data)
				ctx.Session.SendMsg(msg)
			}
			// ---

			sweStatus := GetFromRT(i.Name, swes)

			lastId = Insert(i.Name, host, string(sJson), sweStatus, i.ID, ctx)

			// Socket Broadcast ---
			if ctx.Session.PumpStarted {
				data := models.Step{
					Host:    host,
					Actions: "Getting Puppet Classes from Foreman",
				}
				msg, _ := json.Marshal(data)
				ctx.Session.SendMsg(msg)
			}
			// ---

			scpIds := puppetclass.ApiByHG(host, i.ID, lastId, ctx)

			// Socket Broadcast ---
			if ctx.Session.PumpStarted {
				data := models.Step{
					Host:    host,
					Actions: "Getting Host group parameters from Foreman",
				}
				msg, _ := json.Marshal(data)
				ctx.Session.SendMsg(msg)
			}
			// ---

			HgParams(host, lastId, i.ID, ctx)

			for _, scp := range scpIds {
				scpData := smartclass.SCByPCJsonV2(host, scp, ctx)

				// Socket Broadcast ---
				if ctx.Session.PumpStarted {
					data := models.Step{
						Host:    host,
						Actions: "Getting Smart classes from Foreman",
						State:   scpData.Name,
					}
					msg, _ := json.Marshal(data)
					ctx.Session.SendMsg(msg)
				}
				// ---

				for _, scParam := range scpData.SmartClassParameters {

					// Socket Broadcast ---
					if ctx.Session.PumpStarted {
						data := models.Step{
							Host:    host,
							Actions: "Getting Smart class parameters from Foreman",
							State:   scParam.Parameter,
						}
						msg, _ := json.Marshal(data)
						ctx.Session.SendMsg(msg)
					}
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
