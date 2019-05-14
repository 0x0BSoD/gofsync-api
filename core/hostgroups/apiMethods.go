package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/core/puppetclass"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"sort"
)

// ===============================
// CHECKS
// ===============================
func HostGroupCheck(host string, hostGroupName string, cfg *models.Config) models.HgError {

	var r models.HostGroups

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	body, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(body.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q, hostGroupJson", err)
		}
		if len(r.Results) > 0 {
			return models.HgError{
				ID:        r.Results[0].ID,
				HostGroup: hostGroupName,
				Host:      host,
				Error:     "found",
			}
		}
	}
	return models.HgError{
		ID:        -1,
		HostGroup: hostGroupName,
		Host:      host,
		Error:     "not found",
	}
}

// ===============================
// GET
// ===============================
// Just get HostGroup info by name
func HostGroupJson(host string, hostGroupName string, cfg *models.Config) (models.HGElem, models.HgError) {

	var r models.HostGroups

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	body, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(body.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q, hostGroupJson", err)
		}

		resPc := make(map[string][]models.PuppetClassesWeb)
		puppetClass := puppetclass.GetPCByHgJson(host, r.Results[0].ID, cfg)
		for pcName, subClasses := range puppetClass {
			for _, subClass := range subClasses {
				scData := smartclass.SCByPCJson(host, subClass.ID, cfg)
				var scp []string
				var overrides []models.SCOParams
				for _, i := range scData {
					if !utils.StringInSlice(i.Parameter, scp) {
						scp = append(scp, i.Parameter)
						if i.OverrideValuesCount > 0 {
							sco := smartclass.SCOverridesById(host, i.ID, cfg)
							for _, j := range sco {
								match := fmt.Sprintf("hostgroup=SWE/%s", r.Results[0].Name)
								if j.Match == match {
									jsonVal, _ := json.Marshal(j.Value)
									overrides = append(overrides, models.SCOParams{
										Match:     j.Match,
										Value:     string(jsonVal),
										Parameter: i.Parameter,
									})
								}
							}
						}
					}
				}
				resPc[pcName] = append(resPc[pcName], models.PuppetClassesWeb{
					Subclass:     subClass.Name,
					SmartClasses: scp,
					Overrides:    overrides,
				})
			}
		}
		dbId := r.Results[0].ID
		tmpDbId := CheckHG(r.Results[0].Name, host, cfg)
		if tmpDbId != -1 {
			dbId = tmpDbId
		}

		if len(r.Results) > 0 {

			base := models.HGElem{
				ID:            dbId,
				ForemanID:     r.Results[0].ID,
				Name:          r.Results[0].Name,
				Environment:   r.Results[0].EnvironmentName,
				ParentId:      r.Results[0].Ancestry,
				PuppetClasses: resPc,
			}

			return base, models.HgError{}
		}
	}
	return models.HGElem{}, models.HgError{
		HostGroup: hostGroupName,
		Host:      host,
		Error:     "not found",
	}
}

// ===================================
// Get SWE from Foreman
func GetHostGroups(host string, cfg *models.Config) []models.HostGroup {
	var r models.HostGroups
	uri := fmt.Sprintf("hostgroups?format=json&per_page=%d&search=label+~+SWE", cfg.Api.GetPerPage)
	body, err := utils.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err = json.Unmarshal(body.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q:\n %s\n", err, body.Body)
		}

		var resultsContainer []models.HostGroup

		if r.Total > cfg.Api.GetPerPage {
			pagesRange := utils.Pager(r.Total, cfg.Api.GetPerPage)
			for i := 1; i <= pagesRange; i++ {
				uri := fmt.Sprintf("hostgroups?format=json&page=%d&per_page=%d&search=label+~+SWE", i, cfg.Api.GetPerPage)
				body, err := utils.ForemanAPI("GET", host, uri, "", cfg)
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
		return []models.HostGroup{}
	}
}

// Get SWE Parameters from Foreman
func HgParams(host string, dbID int, sweID int, cfg *models.Config) {
	var r models.HostGroupPContainer
	uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&per_page=%d", sweID, cfg.Api.GetPerPage)
	body, err := utils.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err = json.Unmarshal(body.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q:\n %s\n", err, body.Body)
		}

		if r.Total > cfg.Api.GetPerPage {
			pagesRange := utils.Pager(r.Total, cfg.Api.GetPerPage)
			for i := 1; i <= pagesRange; i++ {

				uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&page=%d&per_page=%d", sweID, i, cfg.Api.GetPerPage)
				body, err := utils.ForemanAPI("GET", host, uri, "", cfg)
				if err == nil {
					err = json.Unmarshal(body.Body, &r)
					if err != nil {
						logger.Error.Printf("%q:\n %s\n", err, body.Body)
					}
					for _, j := range r.Results {
						InsertHGP(dbID, j.Name, j.Value, j.Priority, cfg)
					}
				}
			}
		} else {
			for _, i := range r.Results {
				InsertHGP(dbID, i.Name, i.Value, i.Priority, cfg)
			}
		}
	} else {
		logger.Error.Printf("Error on getting HG Params, %s", err)
	}
}

// Dump HostGroup info by name
func HostGroup(host string, hostGroupName string, cfg *models.Config) {
	var r models.HostGroups
	// Socket Broadcast ---
	msg := models.Step{
		Host:    host,
		Actions: "Getting host group from Foreman",
	}
	utils.BroadCastMsg(cfg, msg)
	// ---
	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	body, err := utils.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(body.Body, &r)
		if err != nil {
			logger.Warning.Printf("%q:\n %s\n", err, body.Body)
		}

		for _, i := range r.Results {

			sJson, _ := json.Marshal(i)
			// Socket Broadcast ---
			msg := models.Step{
				Host:    host,
				Actions: "Saving host group",
			}
			utils.BroadCastMsg(cfg, msg)
			// ---
			lastId := InsertHG(i.Name, host, string(sJson), i.ID, cfg)
			// Socket Broadcast ---
			msg = models.Step{
				Host:    host,
				Actions: "Getting Puppet Classes from Foreman",
			}
			utils.BroadCastMsg(cfg, msg)
			// ---
			scpIds := puppetclass.GetPCByHg(host, i.ID, lastId, cfg)

			// Socket Broadcast ---
			msg = models.Step{
				Host:    host,
				Actions: "Getting Host group parameters from Foreman",
			}
			utils.BroadCastMsg(cfg, msg)
			// ---
			HgParams(host, lastId, i.ID, cfg)

			for _, scp := range scpIds {
				scpData := smartclass.SCByPCJsonV2(host, scp, cfg)
				// Socket Broadcast ---
				msg = models.Step{
					Host:    host,
					Actions: "Getting Smart classes from Foreman",
					State:   scpData.Name,
				}
				utils.BroadCastMsg(cfg, msg)
				// ---
				for _, scParam := range scpData.SmartClassParameters {
					// Socket Broadcast ---
					msg = models.Step{
						Host:    host,
						Actions: "Getting Smart class parameters from Foreman",
						State:   scParam.Name,
					}
					utils.BroadCastMsg(cfg, msg)
					// ---
					scpSummary := smartclass.SCByFId(host, scParam.ID, cfg)
					scId := smartclass.InsertSC(host, scpSummary, cfg)
					if scpSummary.OverrideValuesCount > 0 {
						ovrs := smartclass.SCOverridesById(host, scParam.ID, cfg)
						for _, ovr := range ovrs {
							match := fmt.Sprintf("hostgroup=SWE/%s", i.Name)
							if ovr.Match == match {
								// Socket Broadcast ---
								msg = models.Step{
									Host:    host,
									Actions: "Getting Override from Foreman",
									State:   scParam.Name,
								}
								utils.BroadCastMsg(cfg, msg)
								// ---
								smartclass.InsertSCOverride(scId, ovr, scpSummary.ParameterType, cfg)
							}
						}
					}
				}
			}
		}
	} else {
		logger.Error.Printf("Error on getting HG, %s", err)
	}
	// Socket Broadcast ---
	msg = models.Step{
		Host:    host,
		Actions: "Update done.",
	}
	utils.BroadCastMsg(cfg, msg)
	// ---
}

func DeleteHG(host string, hgId int, cfg *models.Config) error {
	data := GetHG(hgId, cfg)
	uri := fmt.Sprintf("hostgroups/%d", data.ForemanID)
	resp, err := logger.ForemanAPI("DELETE", host, uri, "", cfg)
	logger.Trace.Printf("Response on DELETE HG: %q", resp)

	if err != nil {
		logger.Error.Printf("Error on DELETE HG: %s, uri: %s", err, uri)
		return err
	} else {
		DeleteHGbyId(hgId, cfg)
	}
	return nil
}
