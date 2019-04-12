package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/logger"
	"log"
)

// ===============================
// TYPES & VARS
// ===============================
// For Getting SWE from RackTables
//type RTSWE struct {
//	Name      string `json:"name"`
//	BaseTpl   string `json:"basetpl"`
//	OsVersion string `json:"osversion"`
//	SWEStatus string `json:"swestatus"`
//}
//type RTSWES []RTSWE

// For Getting SWE from Foreman
type SWE struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Title               string `json:"title"`
	SubnetID            int    `json:"subnet_id"`
	SubnetName          string `json:"subnet_name"`
	OperatingSystemID   int    `json:"operatingsystem_id"`
	OperatingSystemName string `json:"operatingsystem_name"`
	DomainID            int    `json:"domain_id"`
	DomainName          string `json:"domain_name"`
	EnvironmentID       int    `json:"environment_id"`
	EnvironmentName     string `json:"environment_name"`
	ComputeProfileId    int    `json:"compute_profile_id"`
	ComputeProfileName  string `json:"compute_profile_name"`
	Ancestry            string `json:"ancestry,omitempty"`
	PuppetProxyId       int    `json:"puppet_proxy_id"`
	PuppetCaProxyId     int    `json:"puppet_ca_proxy_id"`
	PTableId            int    `json:"ptable_id"`
	PTableName          string `json:"ptable_name"`
	MediumId            int    `json:"medium_id"`
	MediumName          string `json:"medium_name"`
	ArchitectureId      int    `json:"architecture_id"`
	ArchitectureName    int    `json:"architecture_name"`
	RealmId             int    `json:"realm_id"`
	RealmName           string `json:"realm_name"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
type SWEContainer struct {
	Results  []SWE  `json:"results"`
	Total    int    `json:"total"`
	SubTotal int    `json:"subtotal"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	Search   string `json:"search"`
}

//  Host Group parameters
type HostGroupPContainer struct {
	Results  []HostGroupP `json:"results"`
	Total    int          `json:"total"`
	SubTotal int          `json:"subtotal"`
	Page     int          `json:"page"`
	PerPage  int          `json:"per_page"`
	Search   string       `json:"search"`
}
type HostGroupP struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Priority int    `json:"priority"`
}
type HostGroupS struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}
type errs struct {
	ID        int    `json:"id"`
	HostGroup string `json:"host_group"`
	Host      string `json:"host"`
	Error     string `json:"error"`
}

// ===============================
// CHECKS
// ===============================
func hostGroupCheck(host string, hostGroupName string) errs {

	var r SWEContainer

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	body, err := ForemanAPI("GET", host, uri, "")
	if err == nil {
		err := json.Unmarshal(body, &r)
		if err != nil {
			logger.Warning.Printf("%q, hostGroupJson", err)
		}
		if len(r.Results) > 0 {
			return errs{
				ID:        r.Results[0].ID,
				HostGroup: hostGroupName,
				Host:      host,
				Error:     "found",
			}
		}
	}
	return errs{
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
func hostGroupJson(host string, hostGroupName string) (HGElem, errs) {

	var r SWEContainer

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	body, err := ForemanAPI("GET", host, uri, "")
	if err == nil {
		err := json.Unmarshal(body, &r)
		if err != nil {
			logger.Warning.Printf("%q, hostGroupJson", err)
		}

		resPc := make(map[string][]PuppetClassesWeb)
		pc := getPCByHgJson(host, r.Results[0].ID)
		for pcName, subClasses := range pc {
			for _, subClass := range subClasses {
				scData := smartClassByPCJson(host, subClass.ID)
				var scp []string
				var ovrs []SCOParams
				for _, i := range scData {
					if !stringInSlice(i.Parameter, scp) {
						scp = append(scp, i.Parameter)
						if i.OverrideValuesCount > 0 {
							sco := scOverridesById(host, i.ID)
							for _, j := range sco {
								match := fmt.Sprintf("hostgroup=SWE/%s", r.Results[0].Name)
								if j.Match == match {
									jsonVal, _ := json.Marshal(j.Value)
									ovrs = append(ovrs, SCOParams{
										Match:     j.Match,
										Value:     string(jsonVal),
										Parameter: i.Parameter,
									})
								}
							}
						}
					}
				}
				resPc[pcName] = append(resPc[pcName], PuppetClassesWeb{
					Subclass:     subClass.Name,
					SmartClasses: scp,
					Overrides:    ovrs,
				})
			}
		}
		dbId := r.Results[0].ID
		tmpDbId := checkHG(r.Results[0].Name, host)
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

			return base, errs{}
		}
	}
	return HGElem{}, errs{
		HostGroup: hostGroupName,
		Host:      host,
		Error:     "not found",
	}
}

// Get SWE from RackTables
//func (swe RTSWE) Get(host string) RTSWES {
//	var r RTSWES
//	body := RTAPI("GET", host,
//		"api/rchwswelookups/search?q=name~.*&fields=name,osversion,basetpl,swestatus&format=json")
//
//	err := json.Unmarshal(body, &r)
//	if err != nil {
//		//log.Printf("%q:\n %s\n", err, body)
//		return []RTSWE{}
//	}
//	return r
//}

// ===================================
// Get SWE from Foreman
func (swe SWE) Get(host string) {
	var r SWEContainer

	uri := fmt.Sprintf("hostgroups?format=json&per_page=%d&search=label+~+SWE", globConf.PerPage)
	body, err := ForemanAPI("GET", host, uri, "")
	if err == nil {
		//log.Printf("%q:\n %s\n", err, body)

		err = json.Unmarshal(body, &r)
		if err != nil {
			log.Fatalf("%q:\n %s\n", err, body)
		}

		var resultsContainer []SWE

		if r.Total > globConf.PerPage {
			pagesRange := Pager(r.Total)
			for i := 1; i <= pagesRange; i++ {
				uri := fmt.Sprintf("hostgroups?format=json&page=%d&per_page=%d&search=label+~+SWE", i, globConf.PerPage)
				body, err := ForemanAPI("GET", host, uri, "")
				if err == nil {
					err = json.Unmarshal(body, &r)
					if err != nil {
						log.Fatalf("%q:\n %s\n", err, body)
					}
					resultsContainer = append(resultsContainer, r.Results...)
				}
			}
		} else {
			resultsContainer = append(resultsContainer, r.Results...)
		}

		for _, i := range resultsContainer {
			sJson, _ := json.Marshal(i)
			lastId := insertHG(i.Name, host, string(sJson), i.ID)
			if lastId != -1 {
				getPCByHg(host, i.ID, lastId)
				hgParams(host, lastId, i.ID)
			}
		}
	} else {
		log.Printf("Error on getting HG, %s", err)
	}
}

// Get SWE Parameters from Foreman
func hgParams(host string, dbID int64, sweID int) {
	var r HostGroupPContainer
	uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&per_page=%d", sweID, globConf.PerPage)
	body, err := ForemanAPI("GET", host, uri, "")
	if err == nil {
		err = json.Unmarshal(body, &r)
		if err != nil {
			log.Fatalf("%q:\n %s\n", err, body)
		}

		if r.Total > globConf.PerPage {
			pagesRange := Pager(r.Total)
			for i := 1; i <= pagesRange; i++ {

				uri := fmt.Sprintf("hostgroups/%d/parameters?format=json&page=%d&per_page=%d", sweID, i, globConf.PerPage)
				body, err := ForemanAPI("GET", host, uri, "")
				if err == nil {
					err = json.Unmarshal(body, &r)
					if err != nil {
						log.Fatalf("%q:\n %s\n", err, body)
					}
					for _, j := range r.Results {
						insertHGP(dbID, j.Name, j.Value, j.Priority)
					}
				}
			}
		} else {
			for _, i := range r.Results {
				insertHGP(dbID, i.Name, i.Value, i.Priority)
			}
		}
	} else {
		log.Printf("Error on getting HG Params, %s", err)
	}
}

// Dump HostGroup info by name
func hostGroup(host string, hostGroupName string) {

	var r SWEContainer

	uri := fmt.Sprintf("hostgroups?search=name+=+%s", hostGroupName)
	body, err := ForemanAPI("GET", host, uri, "")
	if err == nil {
		err := json.Unmarshal(body, &r)
		if err != nil {
			logger.Warning.Printf("%q:\n %s\n", err, body)
		}

		fmt.Println(host, r)

		for _, i := range r.Results {

			sJson, _ := json.Marshal(i)
			lastId := insertHG(i.Name, host, string(sJson), i.ID)
			scpIds := getPCByHg(host, i.ID, lastId)

			hgParams(host, lastId, i.ID)

			for _, scp := range scpIds {
				scpData := smartClassByPCJsonV2(host, scp)
				for _, scParam := range scpData.SmartClassParameters {
					scpSummary := smartClassByFId(host, scParam.ID)
					scId := insertSC(host, scpSummary)
					if scpSummary.OverrideValuesCount > 0 {
						ovrs := scOverridesById(host, scParam.ID)
						for _, ovr := range ovrs {
							match := fmt.Sprintf("hostgroup=SWE/%s", i.Name)
							if ovr.Match == match {
								insertSCOverride(scId, ovr, scpSummary.ParameterType)
							}
						}
					}
				}
			}
		}
	} else {
		logger.Error.Printf("Error on getting HG, %s", err)
	}
}

func deleteHG(host string, hgId int) error {
	data := getHG(hgId)
	uri := fmt.Sprintf("hostgroups/%d", data.ForemanID)
	resp, err := ForemanAPI("DELETE", host, uri, "")

	logger.Info.Printf("Response on DELETE HG: %s", resp)

	if err != nil {
		logger.Error.Printf("Error on DELETE HG: %s, uri: %s", err, uri)
		return err
	} else {
		deleteHGbyId(hgId)
	}
	return nil
}
