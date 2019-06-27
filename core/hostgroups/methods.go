package hostgroups

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

// =====================================================================================================================
// NEW HG
func PushNewHG(data models.HWPostRes, host string, ss *models.Session) (string, error) {
	jDataBase, _ := json.Marshal(models.POSTStructBase{HostGroup: data.BaseInfo})
	response, _ := logger.ForemanAPI("POST", host, "hostgroups", string(jDataBase), ss.Config)
	if response.StatusCode == 200 || response.StatusCode == 201 {
		if len(data.Overrides) > 0 {
			err := PushNewOverride(&data, host, ss)
			if err != nil {
				return "", err
			}
			logger.Info.Printf("crated overrides for HG || %s : %s on %s", ss.UserName, data.BaseInfo.Name, host)
		}
		fmt.Println(data.Parameters)
		fmt.Println(len(data.Parameters))
		if len(data.Parameters) > 0 {
			err := PushNewParameter(&data, response.Body, host, ss)
			if err != nil {
				return "", err
			}
			logger.Info.Printf("crated parameters for HG || %s : %s on %s", ss.UserName, data.BaseInfo.Name, host)
		}
		// Log
		return fmt.Sprintf("crated HG || %s : %s on %s", ss.UserName, data.BaseInfo.Name, host), nil
	}
	return "", utils.NewError(string(response.Body))
}
func PushNewParameter(data *models.HWPostRes, response []byte, host string, ss *models.Session) error {
	var rb models.HostGroup
	err := json.Unmarshal(response, &rb)
	if err != nil {
		return err
	}

	fmt.Println(string(response))
	fmt.Println(rb)

	for _, p := range data.Parameters {

		// Socket Broadcast ---
		msg := models.Step{
			Host:    host,
			Actions: "Submitting parameters",
			State:   fmt.Sprintf("Parameter: %s", p.Name),
		}
		utils.BroadCastMsg(ss, msg)
		// ---

		objP := struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: p.Name, Value: p.Value}
		d := models.POSTStructParameter{HGParam: objP}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("hostgroups/%d/parameters", rb.ID)
		resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ss.Config)
		if err != nil {
			return err
		}
		logger.Info.Println(string(resp.Body), resp.RequestUri)
	}
	return nil
}
func PushNewOverride(data *models.HWPostRes, host string, ss *models.Session) error {
	for _, ovr := range data.Overrides {

		// Socket Broadcast ---
		msg := models.Step{
			Host:    host,
			Actions: "Submitting overrides",
			State:   fmt.Sprintf("Parameter: %s", ovr.Value),
		}
		utils.BroadCastMsg(ss, msg)
		// ---

		p := struct {
			Match string `json:"match"`
			Value string `json:"value"`
		}{Match: ovr.Match, Value: ovr.Value}
		d := models.POSTStructOvrVal{p}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
		resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ss.Config)
		if err != nil {
			return err
		}
		logger.Info.Println(string(resp.Body), resp.RequestUri)
	}
	return nil
}

// UPDATE ==============================================================================================================
func UpdateHG(data models.HWPostRes, host string, ss *models.Session) (string, error) {
	jDataBase, _ := json.Marshal(models.POSTStructBase{HostGroup: data.BaseInfo})
	uri := fmt.Sprintf("hostgroups/%d", data.ExistId)
	response, err := logger.ForemanAPI("PUT", host, uri, string(jDataBase), ss.Config)
	if err == nil {
		if len(data.Overrides) > 0 {
			err := UpdateOverride(&data, host, ss)
			if err != nil {
				return "", err
			}
			logger.Info.Printf("updated overrides for HG || %s : %s on %s", ss.UserName, data.BaseInfo.Name, host)
		}

		if len(data.Parameters) > 0 {
			err := UpdateParameter(&data, response.Body, host, ss)
			if err != nil {
				return "", err
			}
		}
	}

	// Log ============================
	logger.Info.Printf("updated HG || %s : %s on %s", ss.UserName, data.BaseInfo.Name, host)

	// Socket Broadcast ---
	msg := models.Step{
		Host:    host,
		Actions: "Uploading Done!",
	}
	utils.BroadCastMsg(ss, msg)
	// ---

	return fmt.Sprintf("updated HG || %s : %s on %s", ss.UserName, data.BaseInfo.Name, host), nil
}
func UpdateOverride(data *models.HWPostRes, host string, ss *models.Session) error {
	for _, ovr := range data.Overrides {

		// Socket Broadcast ---
		msg := models.Step{
			Host:    host,
			Actions: "Updating overrides",
			State:   fmt.Sprintf("Parameter: %s", ovr.Value),
		}
		utils.BroadCastMsg(ss, msg)
		// ---

		p := struct {
			Match string `json:"match"`
			Value string `json:"value"`
		}{Match: ovr.Match, Value: ovr.Value}
		d := models.POSTStructOvrVal{OverrideValue: p}
		jDataOvr, _ := json.Marshal(d)

		if ovr.OvrForemanId != -1 {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values/%d", ovr.ScForemanId, ovr.OvrForemanId)

			resp, err := logger.ForemanAPI("PUT", host, uri, string(jDataOvr), ss.Config)
			if err != nil {
				return err
			}
			if resp.StatusCode == 404 {
				uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
				resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ss.Config)
				if err != nil {
					return err
				}
				logger.Info.Printf("%s : created Override ForemanID: %d on %s", ss.UserName, ovr.ScForemanId, host)
				logger.Trace.Println(string(resp.Body))
			}
			logger.Info.Printf("%s : updated Override ForemanID: %d, Value: %s on %s", ss.UserName, ovr.ScForemanId, ovr.Value, host)

		} else {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
			resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ss.Config)
			if err != nil {
				return err
			}
			logger.Info.Printf("%s : created Override ForemanID: %d on %s", ss.UserName, ovr.ScForemanId, host)
			logger.Trace.Println(string(resp.Body))
		}
	}
	return nil
}
func UpdateParameter(data *models.HWPostRes, response []byte, host string, ss *models.Session) error {
	var rb models.HostGroup
	err := json.Unmarshal(response, rb)
	if err != nil {
		return err
	}
	for _, p := range data.Parameters {
		// Socket Broadcast ---
		msg := models.Step{
			Host:    host,
			Actions: "Submitting parameters",
			State:   fmt.Sprintf("Parameter: %s", p.Name),
		}
		utils.BroadCastMsg(ss, msg)
		// ---

		objP := struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: p.Name, Value: p.Value}
		d := models.POSTStructParameter{HGParam: objP}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("/api/hostgroups/%d/parameters/%d", rb.ID, p.ForemanID)
		resp, err := logger.ForemanAPI("PUT", host, uri, string(jDataOvr), ss.Config)
		if err != nil {
			return err
		}
		logger.Info.Println(string(resp.Body), resp.RequestUri)
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
func HGDataItem(sHost string, tHost string, hgId int, ss *models.Session) (models.HWPostRes, error) {

	// Source Host Group
	// Socket Broadcast ---
	msg := models.Step{
		Host:    sHost,
		Actions: "Getting source host group data from db",
	}
	utils.BroadCastMsg(ss, msg)
	// ---
	hostGroupData := GetHG(hgId, ss)

	// Step 1. Check if Host Group exist on the host
	// Socket Broadcast ---
	msg = models.Step{
		Host:    tHost,
		Actions: "Getting target host group data from db",
	}
	utils.BroadCastMsg(ss, msg)
	// ---
	hostGroupExistBase := CheckHG(hostGroupData.Name, tHost, ss)
	tmp := HostGroupCheck(tHost, hostGroupData.Name, ss)
	hostGroupExist := tmp.ID

	// Step 2. Check Environment exist on the target host
	// Socket Broadcast ---
	msg = models.Step{
		Host:    tHost,
		Actions: "Getting target environments from db",
	}
	utils.BroadCastMsg(ss, msg)
	// ---

	log.Println("==============================================")
	log.Println(hgId, sHost, tHost, hostGroupData)
	log.Println("==============================================")

	environmentExist := environment.CheckPostEnv(tHost, hostGroupData.Environment, ss)
	if environmentExist == -1 {
		return models.HWPostRes{}, errors.New(fmt.Sprintf("Environment '%s' not exist on %s", hostGroupData.Environment, tHost))
	}

	// Step 3. Get parent Host Group ID on target host
	// Socket Broadcast ---
	msg = models.Step{
		Host:    tHost,
		Actions: "Get parent Host Group ID on target host",
	}
	utils.BroadCastMsg(ss, msg)
	// ---
	parentHGId := CheckHGID("SWE", tHost, ss)
	if parentHGId == -1 {
		return models.HWPostRes{}, errors.New(fmt.Sprintf("Parent Host Group 'SWE' not exist on %s", tHost))
	}

	// Step 4. Get all locations for the target host
	// Socket Broadcast ---
	msg = models.Step{
		Host:    tHost,
		Actions: "Get all locations for the target host",
	}
	utils.BroadCastMsg(ss, msg)
	// ---
	locationsIds := locations.DbAllForemanID(tHost, ss)

	// Step 5. Check Puppet Classes on existing on the target host
	// and
	// Step 6. Get Smart Class data
	var PuppetClassesIds []int
	var SCOverrides []models.HostGroupOverrides
	for pcName, i := range hostGroupData.PuppetClasses {
		// Get Puppet Classes IDs for target Foreman
		subclassLen := len(i)
		currentCounter := 0
		for _, subclass := range i {

			// Socket Broadcast ---
			currentCounter++
			msg = models.Step{
				Host:    tHost,
				Actions: "Get Puppet and Smart Class data",
				State:   fmt.Sprintf("Puppet Class: %s, Smart Class: %s", pcName, subclass.Subclass),
				Counter: currentCounter,
				Total:   subclassLen,
			}
			utils.BroadCastMsg(ss, msg)
			// ---

			targetPCData := puppetclass.DbByName(subclass.Subclass, tHost, ss)
			//sourcePCData := getByNamePC(subclass.Subclass, sHost)

			// If we not have Puppet Class for target host
			if targetPCData.ID == 0 {
				//return HWPostRes{}, errors.New(fmt.Sprintf("Puppet Class '%s' not exist on %s", name, tHost))
			} else {

				// Build Target PC id's and SmartClasses
				PuppetClassesIds = append(PuppetClassesIds, targetPCData.ForemanId)
				var sourceScDataSet []models.SCGetResAdv
				for _, pc := range hostGroupData.PuppetClasses {
					for _, subPc := range pc {
						for _, sc := range subPc.SmartClasses {
							// Get Smart Class data
							sourceScData := smartclass.GetSC(sHost, subclass.Subclass, sc.Name, ss)
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
						targetSC := smartclass.GetSCData(scId, ss)
						scLenght := len(sourceScDataSet)
						currScCount := 0
						for _, sourceSC := range sourceScDataSet {
							currScCount++
							if sourceSC.Name == targetSC.Name {
								srcOvr, _ := smartclass.GetOvrData(sourceSC.ID, hostGroupData.Name, targetSC.Name, ss)
								targetOvr, trgErr := smartclass.GetOvrData(targetSC.ID, hostGroupData.Name, targetSC.Name, ss)
								if srcOvr.SmartClassId != 0 {

									OverrideID := -1

									if trgErr == nil {
										OverrideID = targetOvr.OverrideId
									}

									//fmt.Println("Match: ", srcOvr.Match)
									//fmt.Println("Value: ", srcOvr.Value)
									//fmt.Println("OverrideId: ", OverrideID)
									//fmt.Println("Parameter: ", srcOvr.Parameter)
									//fmt.Println("SmartClassId: ", targetSC.ForemanId)
									//fmt.Println("============================================")
									// Socket Broadcast ---
									msg = models.Step{
										Host:    tHost,
										Actions: "Getting overrides",
										State:   fmt.Sprintf("Parameter: %s", srcOvr.Parameter),
										Counter: currScCount,
										Total:   scLenght,
									}
									utils.BroadCastMsg(ss, msg)
									// ---

									SCOverrides = append(SCOverrides, models.HostGroupOverrides{
										OvrForemanId: OverrideID,
										ScForemanId:  targetSC.ForemanId,
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

	return models.HWPostRes{
		BaseInfo: models.HostGroupBase{
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

func PostCheckHG(tHost string, hgId int, ss *models.Session) bool {
	// Source Host Group
	hostGroupData := GetHG(hgId, ss)
	// Step 1. Check if Host Group exist on the host
	hostGroupExist := CheckHG(hostGroupData.Name, tHost, ss)
	res := false
	if hostGroupExist != -1 {
		res = true
	}
	return res
}

func SaveHGToJson(ss *models.Session) {
	for _, host := range ss.Config.Hosts {
		data := GetHGList(host, ss)
		for _, d := range data {
			hgData := GetHG(d.ID, ss)
			rJson, _ := json.MarshalIndent(hgData, "", "    ")
			path := fmt.Sprintf("/opt/goFsync/HG/%s/%s.json", host, hgData.Name)
			if _, err := os.Stat("/opt/goFsync/HG/" + host); os.IsNotExist(err) {
				err = os.Mkdir("/opt/goFsync/HG/"+host, 0777)
				if err != nil {
					logger.Error.Printf("Error on mkdir: %s", err)
				}
			}
			err := ioutil.WriteFile(path, rJson, 0644)
			if err != nil {
				logger.Error.Printf("Error on writing file: %s", err)
			}
		}
	}

	var out bytes.Buffer
	commitMessage := fmt.Sprintf("Auto commit. Date: %s", time.Now())

	cmd := exec.Command("bash", "HG/lazygit.sh", commitMessage)
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		logger.Error.Println(err)
	}
}

func HGDataNewItem(sHost string, data models.HGElem, ss *models.Session) (models.HWPostRes, error) {
	var PuppetClassesIds []int
	var SCOverrides []models.HostGroupOverrides
	for _, i := range data.PuppetClasses {
		for _, subclass := range i {

			targetPCData := puppetclass.DbByName(subclass.Subclass, sHost, ss)
			//sourcePCData := getByNamePC(subclass.Subclass, sHost)

			// If we not have Puppet Class for target host
			if targetPCData.ID == 0 {
				//return HWPostRes{}, errors.New(fmt.Sprintf("Puppet Class '%s' not exist on %s", name, tHost))
			} else {

				// Build Target PC id's and SmartClasses
				PuppetClassesIds = append(PuppetClassesIds, targetPCData.ForemanId)
				var sourceScDataSet []models.SCGetResAdv
				for _, pc := range data.PuppetClasses {
					for _, subPc := range pc {
						for _, sc := range subPc.SmartClasses {
							// Get Smart Class data
							sourceScData := smartclass.GetSC(sHost, subclass.Subclass, sc.Name, ss)
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
						targetSC := smartclass.GetSCData(scId, ss)
						targetOvr, trgErr := smartclass.GetOvrData(targetSC.ID, data.SourceName, targetSC.Name, ss)
						if targetOvr.SmartClassId != 0 {

							OverrideID := -1

							if trgErr == nil {
								OverrideID = targetOvr.OverrideId
							}

							oldMatch := fmt.Sprintf("hostgroup=SWE/%s", data.SourceName)

							if targetOvr.Match == oldMatch {
								targetOvr.Match = fmt.Sprintf("hostgroup=SWE/%s", data.Name)
							}

							//fmt.Println("Match: ", targetOvr.Match)
							//fmt.Println("Value: ", targetOvr.Value)
							//fmt.Println("OverrideId: ", OverrideID)
							//fmt.Println("Parameter: ", targetOvr.Parameter)
							//fmt.Println("SmartClassId: ", targetSC.ForemanId)
							//fmt.Println("============================================")

							SCOverrides = append(SCOverrides, models.HostGroupOverrides{
								OvrForemanId: OverrideID,
								ScForemanId:  targetSC.ForemanId,
								Match:        targetOvr.Match,
								Value:        targetOvr.Value,
							})
						}
						//}
						//}
					}
				} // if len()
			}
		} // for subclasses
	}
	return models.HWPostRes{
		BaseInfo: models.HostGroupBase{
			Name:           data.Name,
			PuppetClassIds: PuppetClassesIds,
		},
		Overrides:  SCOverrides,
		Parameters: data.Params,
	}, nil
}

func Sync(host string, ss *models.Session) {
	// Host groups ===
	//==========================================================================================================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Filling HostGroups",
		Host:    host,
	}))

	// Socket Broadcast ---
	msg := models.Step{
		Host:    host,
		Actions: "Getting HostGroups",
		State:   "",
	}
	utils.BroadCastMsg(ss, msg)
	// ---

	beforeUpdate := GetForemanIDs(host, ss)
	var afterUpdate []int

	results := GetHostGroups(host, ss)

	// RT SWEs =================================================================================================
	swes := utils.RTbuildObj(HostEnv(host, ss), ss.Config)

	for idx, i := range results {
		// Socket Broadcast ---
		msg := models.Step{
			Host:    host,
			Actions: "Saving HostGroups",
			State:   fmt.Sprintf("HostGroup: %s %d/%d", i.Name, idx+1, len(results)),
		}
		utils.BroadCastMsg(ss, msg)
		// ---
		sJson, _ := json.Marshal(i)

		sweStatus := GetFromRT(i.Name, swes)
		fmt.Printf("Get: %s\tStatus:%s\n", i.Name, sweStatus)

		lastId := Insert(i.Name, host, string(sJson), sweStatus, i.ID, ss)
		afterUpdate = append(afterUpdate, i.ID)

		if lastId != -1 {
			puppetclass.ApiByHG(host, i.ID, lastId, ss)
			HgParams(host, lastId, i.ID, ss)
		}
	}

	for _, i := range beforeUpdate {
		if !utils.Search(afterUpdate, i) {
			fmt.Println("To Delete ", i, host)
		}
	}
}

func StoreHosts(cfg *models.Config) {
	for _, host := range cfg.Hosts {
		InsertHost(host, cfg)
	}
}

//func Compare(cfg *models.Session) {
//	HGList := GetHGList(ss.Config.MasterHost, ss)
//	for _, i := range HGList {
//		for _, h := range ss.Config.Hosts {
//			if h != ss.Config.MasterHost {
//				ch := CheckHG(i.Name, h, ss)
//				state := "nope"
//				if ch != -1 {
//					state = "1"
//				} else {
//					state = "0"
//				}
//				fmt.Println(i.ID)
//				fmt.Println(i.Name)
//				fmt.Println(i.Status)
//				fmt.Println(h, " ================================")
//				insertState(i.Name, h, state, ss)
//			}
//		}
//
//	}
//}
