package hostgroups

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"io/ioutil"
	"os"
)

// =====================================================================================================================
// NEW HG
func PushNewHG(data HWPostRes, host string, ctx *user.GlobalCTX) (string, error) {
	jDataBase, _ := json.Marshal(POSTStructBase{HostGroup: data.BaseInfo})
	fmt.Println(string(jDataBase))
	response, _ := logger.ForemanAPI("POST", host, "hostgroups", string(jDataBase), ctx)
	if response.StatusCode == 200 || response.StatusCode == 201 {
		if len(data.Overrides) > 0 {
			err := PushNewOverride(&data, host, ctx)
			if err != nil {
				return "", err
			}
			logger.Info.Printf("crated overrides for HG || %s : %s on %s", ctx.Session.UserName, data.BaseInfo.Name, host)
		}
		if len(data.Parameters) > 0 {
			err := PushNewParameter(&data, response.Body, host, ctx)
			if err != nil {
				return "", err
			}
			logger.Info.Printf("crated parameters for HG || %s : %s on %s", ctx.Session.UserName, data.BaseInfo.Name, host)
		}
		// Log
		return fmt.Sprintf("crated HG || %s : %s on %s", ctx.Session.UserName, data.BaseInfo.Name, host), nil
	}
	return "", utils.NewError(string(response.Body))
}

func PushNewParameter(data *HWPostRes, response []byte, host string, ctx *user.GlobalCTX) error {
	var rb HostGroupForeman
	err := json.Unmarshal(response, &rb)
	if err != nil {
		return err
	}
	for _, p := range data.Parameters {

		// Socket Broadcast ---
		if ctx.Session.PumpStarted {
			data := models.Step{
				Host:    host,
				Actions: "Submitting parameters",
				State:   fmt.Sprintf("Parameter: %s", p.Name),
			}
			msg, _ := json.Marshal(data)
			ctx.Session.SendMsg(msg)
		}
		// ---

		objP := struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: p.Name, Value: p.Value}
		d := POSTStructParameter{HGParam: objP}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("hostgroups/%d/parameters", rb.ID)
		resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
		if err != nil {
			return err
		}
		logger.Info.Println(string(resp.Body), resp.RequestUri)
	}
	return nil
}
func PushNewOverride(data *HWPostRes, host string, ctx *user.GlobalCTX) error {
	for _, ovr := range data.Overrides {

		// Socket Broadcast ---
		if ctx.Session.PumpStarted {
			data := models.Step{
				Host:    host,
				Actions: "Submitting overrides",
				State:   fmt.Sprintf("Parameter: %s", ovr.Value),
			}
			msg, _ := json.Marshal(data)
			ctx.Session.SendMsg(msg)
		}
		// ---

		p := struct {
			Match string `json:"match"`
			Value string `json:"value"`
		}{Match: ovr.Match, Value: ovr.Value}
		d := POSTStructOvrVal{p}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
		resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
		if err != nil {
			return err
		}
		logger.Info.Println(string(resp.Body), resp.RequestUri)
	}
	return nil
}

// UPDATE ==============================================================================================================
func UpdateHG(data HWPostRes, host string, ctx *user.GlobalCTX) (string, error) {
	jDataBase, _ := json.Marshal(POSTStructBase{HostGroup: data.BaseInfo})
	uri := fmt.Sprintf("hostgroups/%d", data.ExistId)
	response, err := logger.ForemanAPI("PUT", host, uri, string(jDataBase), ctx)
	if err == nil {
		if len(data.Overrides) > 0 {
			err := UpdateOverride(&data, host, ctx)
			if err != nil {
				return "", err
			}
			logger.Info.Printf("updated overrides for HG || %s : %s on %s", ctx.Session.UserName, data.BaseInfo.Name, host)
		}

		if len(data.Parameters) > 0 {
			err := UpdateParameter(&data, response.Body, host, ctx)
			if err != nil {
				return "", err
			}
		}
	}

	// Log ============================
	logger.Info.Printf("updated HG || %s : %s on %s", ctx.Session.UserName, data.BaseInfo.Name, host)
	return fmt.Sprintf("updated HG || %s : %s on %s", ctx.Session.UserName, data.BaseInfo.Name, host), nil
}
func UpdateOverride(data *HWPostRes, host string, ctx *user.GlobalCTX) error {
	for _, ovr := range data.Overrides {

		// Socket Broadcast ---
		if ctx.Session.PumpStarted {
			data := models.Step{
				Host:    host,
				Actions: "Updating overrides",
				State:   fmt.Sprintf("Parameter: %s", ovr.Value),
			}
			msg, _ := json.Marshal(data)
			ctx.Session.SendMsg(msg)
		}
		// ---

		p := struct {
			Match string `json:"match"`
			Value string `json:"value"`
		}{Match: ovr.Match, Value: ovr.Value}
		d := POSTStructOvrVal{OverrideValue: p}
		jDataOvr, _ := json.Marshal(d)

		if ovr.OvrForemanId != -1 {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values/%d", ovr.ScForemanId, ovr.OvrForemanId)

			resp, err := logger.ForemanAPI("PUT", host, uri, string(jDataOvr), ctx)
			if err != nil {
				return err
			}
			if resp.StatusCode == 404 {
				uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
				resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
				if err != nil {
					return err
				}
				logger.Info.Printf("%s : created Override ForemanID: %d on %s", ctx.Session.UserName, ovr.ScForemanId, host)
				logger.Trace.Println(string(resp.Body))
			} else {
				logger.Info.Printf("%s : updated Override ForemanID: %d, Value: %s on %s", ctx.Session.UserName, ovr.ScForemanId, ovr.Value, host)
			}

		} else {
			uri := fmt.Sprintf("smart_class_parameters/%d/override_values", ovr.ScForemanId)
			resp, err := logger.ForemanAPI("POST", host, uri, string(jDataOvr), ctx)
			if err != nil {
				return err
			}
			logger.Info.Printf("%s : created Override ForemanID: %d on %s", ctx.Session.UserName, ovr.ScForemanId, host)
			logger.Trace.Println(string(resp.Body))
		}
	}
	return nil
}
func UpdateParameter(data *HWPostRes, response []byte, host string, ctx *user.GlobalCTX) error {
	var rb HostGroupForeman
	err := json.Unmarshal(response, &rb)
	if err != nil {
		return err
	}
	for _, p := range data.Parameters {
		// Socket Broadcast ---
		if ctx.Session.PumpStarted {
			data := models.Step{
				Host:    host,
				Actions: "Submitting parameters",
				State:   fmt.Sprintf("Parameter: %s", p.Name),
			}
			msg, _ := json.Marshal(data)
			ctx.Session.SendMsg(msg)
		}
		// ---

		objP := struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: p.Name, Value: p.Value}
		d := POSTStructParameter{HGParam: objP}
		jDataOvr, _ := json.Marshal(d)
		uri := fmt.Sprintf("/hostgroups/%d/parameters/%d", rb.ID, p.ForemanID)
		resp, err := logger.ForemanAPI("PUT", host, uri, string(jDataOvr), ctx)
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
func HGDataItem(sHost string, tHost string, hgId int, ctx *user.GlobalCTX) (HWPostRes, error) {

	// Source Host Group

	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    sHost,
			Actions: "Getting source host group data from db",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
	// ---
	hostGroupData := Get(hgId, ctx)

	// Step 1. Check if Host Group exist on the host

	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    tHost,
			Actions: "Getting target host group data from db",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
	// ---
	hostGroupExistBase := ID(hostGroupData.Name, tHost, ctx)
	tmp := HostGroupCheck(tHost, hostGroupData.Name, ctx)
	hostGroupExist := tmp.ID

	// Step 2. Check Environment exist on the target host
	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    tHost,
			Actions: "Getting target environments from db",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
	// ---

	environmentExist := environment.DbForemanID(tHost, hostGroupData.Environment, ctx)
	if environmentExist == -1 {
		return HWPostRes{}, errors.New(fmt.Sprintf("Environment '%s' not exist on %s", hostGroupData.Environment, tHost))
	}

	// Step 3. Get parent Host Group ID on target host
	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    tHost,
			Actions: "Get parent Host Group ID on target host",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
	// ---
	parentHGId := FID("SWE", tHost, ctx)
	if parentHGId == -1 {
		return HWPostRes{}, errors.New(fmt.Sprintf("Parent Host Group 'SWE' not exist on %s", tHost))
	}

	// Step 4. Get all locations for the target host
	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    tHost,
			Actions: "Get all locations for the target host",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
	// ---
	locationsIds := locations.DbAllForemanID(tHost, ctx)

	// Step 5. Check Puppet Classes on existing on the target host
	// and
	// Step 6. Get Smart Class data
	var PuppetClassesIds []int
	var SCOverrides []HostGroupOverrides
	for pcName, i := range hostGroupData.PuppetClasses {
		// Get Puppet Classes IDs for target Foreman
		subclassLen := len(i)
		currentCounter := 0
		for _, subclass := range i {

			// Socket Broadcast ---
			if ctx.Session.PumpStarted {
				currentCounter++
				data := models.Step{
					Host:    tHost,
					Actions: "Get Puppet and Smart Class data",
					State:   fmt.Sprintf("Puppet Class: %s, Smart Class: %s", pcName, subclass.Subclass),
					Counter: currentCounter,
					Total:   subclassLen,
				}
				msg, _ := json.Marshal(data)
				ctx.Session.SendMsg(msg)
			}
			// ---

			targetPCData := puppetclass.DbByName(subclass.Subclass, tHost, ctx)
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
							sourceScData := smartclass.GetSC(sHost, subclass.Subclass, sc.Name, ctx)
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
						scLenght := len(sourceScDataSet)
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

									// Socket Broadcast ---
									if ctx.Session.PumpStarted {
										data := models.Step{
											Host:    tHost,
											Actions: "Getting overrides",
											State:   fmt.Sprintf("Parameter: %s", srcOvr.Parameter),
											Counter: currScCount,
											Total:   scLenght,
										}
										msg, _ := json.Marshal(data)
										ctx.Session.SendMsg(msg)
									}
									// ---

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
	hostGroupExist := ID(hostGroupData.Name, tHost, ctx)
	res := false
	if hostGroupExist != -1 {
		res = true
	}
	return res
}

func SaveHGToJson(ctx *user.GlobalCTX) {
	for _, host := range ctx.Config.Hosts {
		data := OnHost(host, ctx)
		for _, d := range data {
			hgData := Get(d.ID, ctx)
			rJson, _ := json.MarshalIndent(hgData, "", "    ")
			path := fmt.Sprintf("/%s/%s/%s.json", ctx.Config.Git.Directory, host, hgData.Name)
			if _, err := os.Stat(ctx.Config.Git.Directory + "/" + host); os.IsNotExist(err) {
				err = os.Mkdir(ctx.Config.Git.Directory+"/"+host, 0777)
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

	utils.CommitRepo(1, true, ctx)
}

func CommitJsonByHgID(hgID int, host string, ctx *user.GlobalCTX) error {
	hgData := Get(hgID, ctx)
	rJson, _ := json.MarshalIndent(hgData, "", "    ")
	path := fmt.Sprintf("/%s/%s/%s.json", ctx.Config.Git.Directory, host, hgData.Name)
	gitPath := fmt.Sprintf("%s/%s.json", host, hgData.Name)
	if _, err := os.Stat(ctx.Config.Git.Directory + "/" + host); os.IsNotExist(err) {
		err = os.Mkdir(ctx.Config.Git.Directory+"/"+host, 0777)
		if err != nil {
			logger.Error.Printf("Error on mkdir: %s", err)
			return err
		}
	}
	err := ioutil.WriteFile(path, rJson, 0644)
	if err != nil {
		logger.Error.Printf("Error on writing file: %s", err)
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

	// =====
	for _, puppetClass := range hostGroupJSON.PuppetClasses {
		for _, subclass := range puppetClass {
			foremanID := puppetclass.ForemanID(subclass.Subclass, host, ctx)
			fmt.Println("== ", subclass.Subclass, " == ", foremanID)
			puppetClassesIds = append(puppetClassesIds, foremanID)

			for _, sc := range subclass.Overrides {

				SmartClass := smartclass.GetSCData(sc.SmartClassId, ctx)

				smartClassOverrides = append(smartClassOverrides, HostGroupOverrides{
					OvrForemanId: sc.ForemanID,
					ScForemanId:  SmartClass.ForemanID,
					Match:        sc.Match,
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

func Sync(host string, ctx *user.GlobalCTX) {
	// Host groups ===
	//==========================================================================================================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Filling HostGroups",
		Host:    host,
	}))

	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    host,
			Actions: "Getting HostGroups",
			State:   "",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
	// ---

	beforeUpdate := FIDs(host, ctx)
	var afterUpdate []int

	results := GetHostGroups(host, ctx)

	// RT SWEs =================================================================================================
	swes := RTBuildObj(PuppetHostEnv(host, ctx), ctx)

	for idx, i := range results {
		// Socket Broadcast ---
		if ctx.Session.PumpStarted {
			data := models.Step{
				Host:    host,
				Actions: "Saving HostGroups",
				State:   fmt.Sprintf("HostGroup: %s %d/%d", i.Name, idx+1, len(results)),
			}
			msg, _ := json.Marshal(data)
			ctx.Session.SendMsg(msg)
		}
		// ---
		sJson, _ := json.Marshal(i)

		sweStatus := GetFromRT(i.Name, swes)
		fmt.Printf("Get: %s\tStatus:%s\n", i.Name, sweStatus)

		lastId := Insert(i.Name, host, string(sJson), sweStatus, i.ID, ctx)
		afterUpdate = append(afterUpdate, i.ID)

		if lastId != -1 {
			puppetclass.ApiByHG(host, i.ID, lastId, ctx)
			HgParams(host, lastId, i.ID, ctx)
		}
	}

	for _, i := range beforeUpdate {
		if !utils.Search(afterUpdate, i) {
			fmt.Println("Deleting ... ", i, host)
			name := Name(i, host, ctx)
			Delete(i, host, ctx)
			rmJSON(name, host, ctx)
		}
	}
}

func StoreHosts(cfg *models.Config) {
	for _, host := range cfg.Hosts {
		InsertHost(host, cfg)
	}
}

//func Compare(cfg *models.Session) {
//	HGList := GetHGList(ctx.MasterHost, ctx)
//	for _, i := range HGList {
//		for _, h := range ctx.Hosts {
//			if h != ctx.MasterHost {
//				ch := CheckHG(i.Name, h, ctx)
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
//				insertState(i.Name, h, state, ctx)
//			}
//		}
//
//	}
//}

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
