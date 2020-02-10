package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"io/ioutil"
	"os"
	"strings"
)

// =====================================================================================================================
// NEW HG
func PushNewHG(data HWPostRes, host string, ctx *user.GlobalCTX) (string, error) {

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		HostName:  host,
		Resource:  models.HostGroup,
		Operation: "submit",
		UserName:  ctx.Session.UserName,
		AdditionalData: models.CommonOperation{
			Message:   "Submitting HostGroup to foreman",
			HostGroup: data.BaseInfo.Name,
		},
	})
	// ---

	jDataBase, _ := json.Marshal(POSTStructBase{HostGroup: data.BaseInfo})
	response, _ := utils.ForemanAPI("POST", host, "hostgroups", string(jDataBase), ctx)
	if response.StatusCode == 200 || response.StatusCode == 201 {
		if len(data.Overrides) > 0 {
			err := PushNewOverride(&data, host, ctx)
			if err != nil {
				return "", err
			}
		}
		if len(data.Parameters) > 0 {
			err := PushNewParameter(&data, response.Body, host, ctx)
			if err != nil {
				return "", err
			}
		}
		// Log
		utils.Actions.Printf("[hostgroup] |crated| %s : HG:%s on %s", ctx.Session.UserName, data.BaseInfo.Name, host)
		return fmt.Sprintf("[hostgroup] |crated| %s : HG:%s on %s", ctx.Session.UserName, data.BaseInfo.Name, host), nil
	}
	return "", fmt.Errorf(string(response.Body))
}

func PushNewParameter(data *HWPostRes, response []byte, host string, ctx *user.GlobalCTX) error {

	var rb HostGroupForeman
	err := json.Unmarshal(response, &rb)
	if err != nil {
		return err
	}
	count := 1
	aLen := len(data.Parameters)
	for _, p := range data.Parameters {

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			HostName:  host,
			Resource:  models.HostGroup,
			Operation: "submit",
			UserName:  ctx.Session.UserName,
			AdditionalData: models.CommonOperation{
				HostGroup: data.BaseInfo.Name,
				Message:   "Submitting HostGroup Parameter to foreman",
				Item:      fmt.Sprintf("%s=>%s", p.Name, p.Value),
				Total:     aLen,
				Current:   count,
			},
		})
		// ---

		objP := struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: p.Name, Value: p.Value}
		d := POSTStructParameter{HGParam: objP}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("hostgroups/%d/parameters", rb.ID)
		_, err := utils.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
		if err != nil {
			return err
		}
		utils.Actions.Printf("[created] |updated| %s : N:%sV:%s on %s", ctx.Session.UserName, p.Name, p.Value, host)
		count++
	}
	return nil
}
func PushNewOverride(data *HWPostRes, host string, ctx *user.GlobalCTX) error {
	count := 1
	aLen := len(data.Overrides)
	for _, ovr := range data.Overrides {

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			HostName:  host,
			Resource:  models.HostGroup,
			Operation: "submit",
			UserName:  ctx.Session.UserName,
			AdditionalData: models.CommonOperation{
				HostGroup: data.BaseInfo.Name,
				Message:   "Submitting Override Value to foreman",
				Item:      ovr.Value,
				Total:     aLen,
				Current:   count,
			},
		})
		// ---

		p := struct {
			Match string `json:"match"`
			Value string `json:"value"`
		}{Match: ovr.Match, Value: ovr.Value}
		d := POSTStructOvrVal{p}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
		_, err := utils.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
		if err != nil {
			return err
		}
		utils.Actions.Printf("[override] |created| %s : M:%sV:%s on %s", ctx.Session.UserName, ovr.Match, ovr.Value, host)
		count++
	}
	return nil
}

// UPDATE ==============================================================================================================
func UpdateHG(data HWPostRes, host string, ctx *user.GlobalCTX) (string, error) {

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		HostName:  host,
		Resource:  models.HostGroup,
		Operation: "submit",
		UserName:  ctx.Session.UserName,
		AdditionalData: models.CommonOperation{
			Message:   "Updating HostGroup on foreman",
			HostGroup: data.BaseInfo.Name,
		},
	})
	// ---

	jDataBase, _ := json.Marshal(POSTStructBase{HostGroup: data.BaseInfo})
	uri := fmt.Sprintf("hostgroups/%d", data.ExistId)
	response, err := utils.ForemanAPI("PUT", host, uri, string(jDataBase), ctx)
	if err == nil {
		if len(data.Overrides) > 0 {
			err := UpdateOverride(&data, host, ctx)
			if err != nil {
				return "", err
			}
		}

		if len(data.Parameters) > 0 {
			err := UpdateParameter(&data, response.Body, host, ctx)
			if err != nil {
				return "", err
			}
		}
	}

	// Log ============================
	utils.Actions.Printf("[hostgroup] |updated| %s : HG:%s on %s", ctx.Session.UserName, data.BaseInfo.Name, host)
	return fmt.Sprintf("[hostgroup] |updated| %s : HG:%s on %s", ctx.Session.UserName, data.BaseInfo.Name, host), nil
}

func UpdateOverride(data *HWPostRes, host string, ctx *user.GlobalCTX) error {
	count := 1
	aLen := len(data.Overrides)
	for _, ovr := range data.Overrides {

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			HostName:  host,
			Resource:  models.HostGroup,
			Operation: "submit",
			UserName:  ctx.Session.UserName,
			AdditionalData: models.CommonOperation{
				HostGroup: data.BaseInfo.Name,
				Message:   "Updating Override Value on foreman",
				Item:      ovr.Value,
				Total:     aLen,
				Current:   count,
			},
		})
		// ---

		p := struct {
			Match string `json:"match"`
			Value string `json:"value"`
		}{Match: ovr.Match, Value: ovr.Value}
		d := POSTStructOvrVal{OverrideValue: p}
		jDataOvr, _ := json.Marshal(d)

		if ovr.OvrForemanId != -1 {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values/%d", ovr.ScForemanId, ovr.OvrForemanId)

			resp, err := utils.ForemanAPI("PUT", host, uri, string(jDataOvr), ctx)
			if err != nil {
				return err
			}
			if resp.StatusCode == 404 {
				uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
				_, err := utils.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
				if err != nil {
					return err
				}
				utils.Actions.Printf("[override] |created| %s : M:%sV:%s on %s", ctx.Session.UserName, ovr.Match, ovr.Value, host)
			} else {
				utils.Actions.Printf("[override] |updated| %s : M:%sV:%s on %s", ctx.Session.UserName, ovr.Match, ovr.Value, host)
			}

		} else {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
			_, err := utils.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
			if err != nil {
				return err
			}
			utils.Actions.Printf("[override] |created| %s : M:%sV:%s on %s", ctx.Session.UserName, ovr.Match, ovr.Value, host)
		}
		count++
	}
	return nil
}
func UpdateParameter(data *HWPostRes, response []byte, host string, ctx *user.GlobalCTX) error {
	var rb HostGroupForeman
	err := json.Unmarshal(response, &rb)
	if err != nil {
		return err
	}
	count := 1
	aLen := len(data.Parameters)
	for _, p := range data.Parameters {

		// Socket Broadcast ---
		ctx.Session.SendMsg(models.WSMessage{
			Broadcast: false,
			HostName:  host,
			Resource:  models.HostGroup,
			Operation: "submit",
			UserName:  ctx.Session.UserName,
			AdditionalData: models.CommonOperation{
				HostGroup: data.BaseInfo.Name,
				Message:   "Updating HostGroup Parameter on foreman",
				Item:      fmt.Sprintf("%s=>%s", p.Name, p.Value),
				Total:     aLen,
				Current:   count,
			},
		})
		// ---

		objP := struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: p.Name, Value: p.Value}
		d := POSTStructParameter{HGParam: objP}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("/hostgroups/%d/parameters/%d", rb.ID, p.ForemanID)
		_, err := utils.ForemanAPI("PUT", host, uri, string(jDataOvr), ctx)
		if err != nil {
			return err
		}
		utils.Actions.Printf("[parameter] |updated| %s : N:%sV:%s on %s", ctx.Session.UserName, p.Name, p.Value, host)
		count++
	}
	return nil
}

// =====================================================================================================================
// Build object for POST to target Foreman
// Steps:
// 1. is exist
// 2. env
// 3. parent id on target host
// 4. get all locations for the target host
// 5. All puppet classes exist on target host
// 6. Smart class ids  on target host
// 7. overrides for smart classes
// 8. POST
func HGDataItem(sHost string, tHost string, hgId int, ctx *user.GlobalCTX) (HWPostRes, error) {

	// ---
	// Source Host Group
	hostGroupData := Get(hgId, ctx)

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: false,
		HostName:  sHost,
		Resource:  models.HostGroup,
		Operation: "submit",
		UserName:  ctx.Session.UserName,
		AdditionalData: models.CommonOperation{
			Message:   "Generating HostGroup JSON",
			HostGroup: hostGroupData.Name,
		},
	})
	// ---

	// Step 1. Check if Host Group exist on the host
	hostGroupExistBase := ID(ctx.Config.Hosts[tHost], hostGroupData.Name, ctx)
	tmp := HostGroupCheckName(tHost, hostGroupData.Name, ctx)
	hostGroupExist := tmp.ID

	// Step 2. Check Environment exist on the target host
	environmentExist := environment.ForemanID(ctx.Config.Hosts[tHost], hostGroupData.Environment, ctx)
	if environmentExist == -1 {
		return HWPostRes{}, fmt.Errorf("environment '%s' not exist on %s", hostGroupData.Environment, tHost)
	}

	// Step 3. Get parent Host Group ID on target host
	parentHGId := ForemanID(ctx.Config.Hosts[tHost], "SWE", ctx)
	if parentHGId == -1 {
		return HWPostRes{}, fmt.Errorf("parent Host Group 'SWE' not exist on %s", tHost)
	}

	// Step 4. Get all locations for the target host
	locationsIds := locations.DbAllForemanID(ctx.Config.Hosts[tHost], ctx)

	// Step 5. Check Puppet Classes on existing on the target host
	// and
	// Step 6. Get Smart Class data
	var PuppetClassesIds []int
	var SCOverrides []HostGroupOverrides
	for _, i := range hostGroupData.PuppetClasses {
		// Get Puppet Classes IDs for target Foreman
		for _, subclass := range i {
			targetPCData := puppetclass.DbByName(ctx.Config.Hosts[tHost], subclass.Subclass, ctx)
			//sourcePCData := getByNamePC(subclass.Subclass, sHost)

			// If we not have Puppet Class for target host
			if targetPCData.ID == 0 {
				//return HWPostRes{}, errors.New(fmt.Sprintf("Puppet Class '%s' not exist on %s", name, tHost))
			} else {

				// Build Target PC id's and SmartClasses
				PuppetClassesIds = append(PuppetClassesIds, targetPCData.ForemanId)
				var sourceScDataSet []smartclass.SCGetResAdv
				for _, pc := range hostGroupData.PuppetClasses {
					for _, subPc := range pc {
						for _, sc := range subPc.SmartClasses {
							// Get Smart Class data
							sourceScData := smartclass.GetSC(ctx.Config.Hosts[sHost], subclass.Subclass, sc.Name, ctx)
							// If source have overrides
							if sourceScData.OverrideValuesCount > 0 {
								sourceScDataSet = append(sourceScDataSet, sourceScData)
							}
						}
					}
				}

				// Step 7. Overrides for smart classes
				// Iterate the Source Smart classes and target Smart classes and if SC exist in both
				// check if we have overrides if true - add to result
				if len(targetPCData.SCIDs) > 0 {
					for _, scId := range utils.Integers(targetPCData.SCIDs) {
						targetSC := smartclass.GetSCData(scId, ctx)
						currScCount := 0
						for _, sourceSC := range sourceScDataSet {
							currScCount++
							if sourceSC.Name == targetSC.Name {
								srcOvr, _ := smartclass.GetOvrData(sourceSC.ID, hostGroupData.Name, targetSC.Name, ctx)
								targetOvr, trgErr := smartclass.GetOvrData(targetSC.ID, hostGroupData.Name, targetSC.Name, ctx)
								if srcOvr.SmartClassId != 0 {
									OverrideID := -1
									if trgErr == nil {
										OverrideID = targetOvr.ForemanID
									}
									SCOverrides = append(SCOverrides, HostGroupOverrides{
										OvrForemanId: OverrideID,
										ScForemanId:  targetSC.ForemanID,
										Match:        srcOvr.Match,
										Value:        srcOvr.Value,
									})
								}
							}
						}
					}
				} // if len()
			}
		} // for subclasses
	}
	return HWPostRes{
		BaseInfo: HostGroupBase{
			Name:           hostGroupData.Name,
			ParentId:       parentHGId,
			EnvironmentId:  environmentExist,
			LocationIds:    locationsIds,
			PuppetClassIds: PuppetClassesIds,
		},
		Parameters: hostGroupData.Params,
		Overrides:  SCOverrides,
		DBHGExist:  hostGroupExistBase,
		ExistId:    hostGroupExist,
	}, nil
}

func PostCheckHG(tHost string, hgId int, ctx *user.GlobalCTX) bool {
	// Source Host Group
	hostGroupData := Get(hgId, ctx)
	// Step 1. Check if Host Group exist on the host
	hostGroupExist := ID(ctx.Config.Hosts[tHost], hostGroupData.Name, ctx)
	res := false
	if hostGroupExist != -1 {
		res = true
	}
	return res
}

func SaveHGToJson(ctx *user.GlobalCTX) {
	for hostname, ID := range ctx.Config.Hosts {
		data := OnHost(ID, ctx)
		for _, d := range data {
			hgData := Get(d.ID, ctx)
			rJson, _ := json.MarshalIndent(hgData, "", "    ")
			path := fmt.Sprintf("/%s/%s/%s.json", ctx.Config.Git.Directory, hostname, hgData.Name)
			if _, err := os.Stat(ctx.Config.Git.Directory + "/" + hostname); os.IsNotExist(err) {
				err = os.Mkdir(ctx.Config.Git.Directory+"/"+hostname, 0777)
				if err != nil {
					utils.Error.Printf("Error on mkdir: %s", err)
				}
			}
			err := ioutil.WriteFile(path, rJson, 0644)
			if err != nil {
				utils.Error.Printf("Error on writing file: %s", err)
			}
		}
	}
}

func CommitJsonByHgID(hgID int, host string, ctx *user.GlobalCTX) error {
	hgData := Get(hgID, ctx)
	rJson, _ := json.MarshalIndent(hgData, "", "    ")
	path := fmt.Sprintf("/%s/%s/%s.json", ctx.Config.Git.Directory, host, hgData.Name)
	gitPath := fmt.Sprintf("%s/%s.json", host, hgData.Name)
	if _, err := os.Stat(ctx.Config.Git.Directory + "/" + host); os.IsNotExist(err) {
		err = os.Mkdir(ctx.Config.Git.Directory+"/"+host, 0777)
		if err != nil {
			utils.Error.Printf("Error on mkdir: %s", err)
			return err
		}
	}
	err := ioutil.WriteFile(path, rJson, 0644)
	if err != nil {
		utils.Error.Printf("Error on writing file: %s", err)
		return err
	}

	utils.AddToRepo(gitPath, ctx)
	utils.CommitRepo(1, false, ctx)

	return nil
}

func HGDataNewItem(host string, hostGroupJSON HGElem, ctx *user.GlobalCTX) (HWPostRes, error) {

	// VARS
	var puppetClassesIds []int
	var smartClassOverrides []HostGroupOverrides
	match := fmt.Sprintf("hostgroup=SWE/%s", hostGroupJSON.Name)

	// =====
	for _, puppetClass := range hostGroupJSON.PuppetClasses {
		for _, subclass := range puppetClass {
			foremanID := puppetclass.ForemanID(ctx.Config.Hosts[host], subclass.Subclass, ctx)
			puppetClassesIds = append(puppetClassesIds, foremanID)
			for _, sc := range subclass.Overrides {
				SmartClass := smartclass.GetSCData(sc.SmartClassId, ctx)
				smartClassOverrides = append(smartClassOverrides, HostGroupOverrides{
					OvrForemanId: sc.ForemanID,
					ScForemanId:  SmartClass.ForemanID,
					Match:        match,
					Value:        sc.Value,
				})
			}
		}
	}

	return HWPostRes{
		BaseInfo: HostGroupBase{
			Name:           hostGroupJSON.Name,
			PuppetClassIds: puppetClassesIds,
		},
		Overrides:  smartClassOverrides,
		Parameters: hostGroupJSON.Params,
	}, nil
}

func RTBuildObj(env string, ctx *user.GlobalCTX) map[string]string {
	// RT SWEs =================================================================================================
	var swes []RackTablesSWE
	result := make(map[string]string)

	var rtHost string
	if env == "stage" {
		rtHost = ctx.Config.RackTables.Stage
	} else if env == "prod" {
		rtHost = ctx.Config.RackTables.Production
	}

	body, err := utils.RackTablesAPI("GET", rtHost, "rchwswelookups/search?q=name~.*&fields=name,swestatus&format=json", "", ctx)
	if err != nil {
		utils.Error.Println(err)
	}
	err = json.Unmarshal(body.Body, &swes)
	if err != nil {
		utils.Warning.Printf("%q:\n %s\n", err, body.Body)
	}
	// RT SWEs =================================================================================================
	for _, swe := range swes {
		result[swe.Name] = swe.SweStatus
	}
	return result
}

func rmJSON(name, host string, ctx *user.GlobalCTX) {
	fName := fmt.Sprintf("%s.json", name)
	path := ctx.Config.Git.Directory + "/" + host + "/" + fName
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	} else {
		err := os.Remove(path)
		if err != nil {
			utils.Error.Println(err)
		}
	}
}

func CompareHGWorker(first, second HGElem) bool {
	if first.ForemanID != second.ForemanID {
		return false
	}

	if len(first.PuppetClasses) != len(second.PuppetClasses) {
		return false
	}

	for fClass, fp := range first.PuppetClasses {
		for sClass, sp := range second.PuppetClasses {

			if fClass == sClass {
				if len(fp) != len(sp) {
					fmt.Println("subclasses", len(fp), len(sp))
					return false
				}
				for _, subFirst := range fp {
					for _, subSecond := range sp {
						if subFirst.Subclass == subSecond.Subclass {
							if len(subFirst.SmartClasses) != len(subSecond.SmartClasses) {
								return false
							}
							if len(subFirst.Overrides) > 0 && len(subSecond.Overrides) > 0 {
								for _, fOvr := range subFirst.Overrides {
									for _, sOvr := range subSecond.Overrides {
										if fOvr.Parameter == sOvr.Parameter {
											fmt.Println(fOvr)
											fmt.Println(sOvr)
											i := strings.Trim(fOvr.Value, "\"")
											ii := strings.Trim(sOvr.Value, "\"")
											if i != ii {
												fmt.Println(i)
												fmt.Println(ii)
												return false
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return true
}
