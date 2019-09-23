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
	"time"
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
	if err != nil {
		return HGElem{}, HgError{
			HostGroup: hostGroupName,
			Host:      host,
			Error:     "not found",
		}
	}

	err = json.Unmarshal(body.Body, &r)
	if err != nil {
		logger.Warning.Printf("%q, hostGroupJson", err)
		return HGElem{}, HgError{
			HostGroup: hostGroupName,
			Host:      host,
			Error:     "not found",
		}
	}

	puppetClass := puppetclass.ApiByHGJson(host, r.Results[0].ID, ctx)
	resPc := make(map[string][]puppetclass.PuppetClassesWeb, len(puppetClass))

	ts := time.Now()
	fmt.Println("For started")
	for pcName, subClasses := range puppetClass {
		for _, subClass := range subClasses {
			scData := smartclass.SCByPCJson(host, subClass.ID, ctx)
			scp := make([]smartclass.SmartClass, 0, len(scData))
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
						overridesInner := make([]smartclass.SCOParams, 0, len(sco))
						for _, j := range sco {
							match := fmt.Sprintf("hostgroup=SWE/%s", r.Results[0].Name)
							if j.Match == match {
								jsonVal, _ := json.Marshal(j.Value)
								overridesInner = append(overridesInner, smartclass.SCOParams{
									Match:     j.Match,
									Value:     string(jsonVal),
									Parameter: i.Parameter,
								})
							}
						}
						overrides = append(overrides, overridesInner...)
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

	fmt.Println("For done, ", time.Since(ts))

	dbId := r.Results[0].ID
	tmpDbId := ID(r.Results[0].Name, host, ctx)
	if tmpDbId != -1 {
		dbId = tmpDbId
	}

	if len(r.Results) > 0 {
		return HGElem{
			ID:            dbId,
			ForemanID:     r.Results[0].ID,
			Name:          r.Results[0].Name,
			Environment:   r.Results[0].EnvironmentName,
			ParentId:      r.Results[0].Ancestry,
			PuppetClasses: resPc,
		}, HgError{}
	} else {
		return HGElem{}, HgError{
			HostGroup: hostGroupName,
			Host:      host,
			Error:     "not found",
		}
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
func HostGroup(host string, hostGroupName string, ctx *user.GlobalCTX) (int, error) {
	var r HostGroups
	lastId := -1

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		Operation: "getHG",
		Data: models.Step{
			Host:  host,
			State: "running",
		},
	})
	// ---

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	response, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err == nil && response.StatusCode != 500 {
		err := json.Unmarshal(response.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q:\n %s\n", err, response.Body)
		}
		swes := RTBuildObj(PuppetHostEnv(host, ctx), ctx)
		for _, i := range r.Results {

			sJson, _ := json.Marshal(i)

			// Socket Broadcast ---
			ctx.Session.SendMsg(models.WSMessage{
				Broadcast: false,
				Operation: "getHG",
				Data: models.Step{
					Host:  host,
					State: "saving",
				},
			})
			// ---

			sweStatus := GetFromRT(i.Name, swes)

			lastId = Insert(i.Name, host, string(sJson), sweStatus, i.ID, ctx)

			// Socket Broadcast ---
			ctx.Session.SendMsg(models.WSMessage{
				Broadcast: false,
				Operation: "getPC",
				Data: models.Step{
					Host:  host,
					State: "running",
				},
			})
			// ---

			scpIds := puppetclass.ApiByHG(host, i.ID, lastId, ctx)

			// Socket Broadcast ---
			ctx.Session.SendMsg(models.WSMessage{
				Broadcast: false,
				Operation: "getHGParameters",
				Data: models.Step{
					Host:  host,
					State: "running",
				},
			})
			// ---

			HgParams(host, lastId, i.ID, ctx)

			for _, scp := range scpIds {
				scpData := smartclass.SCByPCJsonV2(host, scp, ctx)

				// Socket Broadcast ---
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: false,
					Operation: "getSC",
					Data: models.Step{
						Host:  host,
						State: "running",
					},
				})
				// ---

				for _, scParam := range scpData.SmartClassParameters {

					// Socket Broadcast ---
					ctx.Session.SendMsg(models.WSMessage{
						Broadcast: false,
						Operation: "getSC",
						Data: models.Step{
							Host:  host,
							Item:  scParam.Parameter,
							State: "saving",
						},
					})
					// ---

					scpSummary := smartclass.SCByFId(host, scParam.ID, ctx)
					smartclass.InsertSC(host, scpSummary, ctx)
				}
			}
		}

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			Operation: "done",
		})
		// ---

	} else {
		logger.Error.Printf("Error on getting HG, %s", fmt.Errorf(string(response.Body)))
		return 0, fmt.Errorf(string(response.Body))
	}

	return lastId, nil
}
