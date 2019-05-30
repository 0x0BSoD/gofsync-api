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
	"os"
	"os/exec"
	"time"
)

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
func HGDataItem(sHost string, tHost string, hgId int, cfg *models.Config) (models.HWPostRes, error) {

	// Source Host Group
	// Socket Broadcast ---
	msg := models.Step{
		Host:    sHost,
		Actions: "Getting source host group data from db",
	}
	utils.BroadCastMsg(cfg, msg)
	// ---
	hostGroupData := GetHG(hgId, cfg)

	// Step 1. Check if Host Group exist on the host
	// Socket Broadcast ---
	msg = models.Step{
		Host:    tHost,
		Actions: "Getting target host group data from db",
	}
	utils.BroadCastMsg(cfg, msg)
	// ---
	hostGroupExistBase := CheckHG(hostGroupData.Name, tHost, cfg)
	tmp := HostGroupCheck(tHost, hostGroupData.Name, cfg)
	hostGroupExist := tmp.ID

	// Step 2. Check Environment exist on the target host
	// Socket Broadcast ---
	msg = models.Step{
		Host:    tHost,
		Actions: "Getting target environments from db",
	}
	utils.BroadCastMsg(cfg, msg)
	// ---
	environmentExist := environment.CheckPostEnv(tHost, hostGroupData.Environment, cfg)
	if environmentExist == -1 {
		return models.HWPostRes{}, errors.New(fmt.Sprintf("Environment '%s' not exist on %s", hostGroupData.Environment, tHost))
	}

	// Step 3. Get parent Host Group ID on target host
	// Socket Broadcast ---
	msg = models.Step{
		Host:    tHost,
		Actions: "Get parent Host Group ID on target host",
	}
	utils.BroadCastMsg(cfg, msg)
	// ---
	parentHGId := CheckHGID("SWE", tHost, cfg)
	if parentHGId == -1 {
		return models.HWPostRes{}, errors.New(fmt.Sprintf("Parent Host Group 'SWE' not exist on %s", tHost))
	}

	// Step 4. Get all locations for the target host
	// Socket Broadcast ---
	msg = models.Step{
		Host:    tHost,
		Actions: "Get all locations for the target host",
	}
	utils.BroadCastMsg(cfg, msg)
	// ---
	locationsIds := locations.DbAllForemanID(tHost, cfg)

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
			utils.BroadCastMsg(cfg, msg)
			// ---

			targetPCData := puppetclass.DbByName(subclass.Subclass, tHost, cfg)
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
							sourceScData := smartclass.GetSC(sHost, subclass.Subclass, sc.Name, cfg)
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
						targetSC := smartclass.GetSCData(scId, cfg)
						scLenght := len(sourceScDataSet)
						currScCount := 0
						for _, sourceSC := range sourceScDataSet {
							currScCount++
							if sourceSC.Name == targetSC.Name {
								srcOvr, _ := smartclass.GetOvrData(sourceSC.ID, hostGroupData.Name, targetSC.Name, cfg)
								targetOvr, trgErr := smartclass.GetOvrData(targetSC.ID, hostGroupData.Name, targetSC.Name, cfg)
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
									utils.BroadCastMsg(cfg, msg)
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
		Overrides: SCOverrides,
		DBHGExist: hostGroupExistBase,
		ExistId:   hostGroupExist,
	}, nil
}

func PostCheckHG(tHost string, hgId int, cfg *models.Config) bool {
	// Source Host Group
	hostGroupData := GetHG(hgId, cfg)
	// Step 1. Check if Host Group exist on the host
	hostGroupExist := CheckHG(hostGroupData.Name, tHost, cfg)
	res := false
	if hostGroupExist != -1 {
		res = true
	}
	return res
}

func SaveHGToJson(cfg *models.Config) {
	for _, host := range cfg.Hosts {
		data := GetHGList(host, cfg)
		for _, d := range data {
			hgData := GetHG(d.ID, cfg)
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

func Sync(host string, cfg *models.Config) {
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
	utils.BroadCastMsg(cfg, msg)
	// ---

	beforeUpdate := GetForemanIDs(host, cfg)
	var afterUpdate []int

	results := GetHostGroups(host, cfg)

	for idx, i := range results {
		// Socket Broadcast ---
		msg := models.Step{
			Host:    host,
			Actions: "Saving HostGroups",
			State:   fmt.Sprintf("HostGroup: %s %d/%d", i.Name, idx+1, len(results)),
		}
		utils.BroadCastMsg(cfg, msg)
		// ---
		sJson, _ := json.Marshal(i)
		sweStatus := GetFromRT(i.Name, host, cfg)
		lastId := Insert(i.Name, host, string(sJson), sweStatus, i.ID, cfg)
		afterUpdate = append(afterUpdate, i.ID)

		if lastId != -1 {
			puppetclass.ApiByHG(host, i.ID, lastId, cfg)
			HgParams(host, lastId, i.ID, cfg)
		}
	}

	for _, i := range beforeUpdate {
		if !utils.IntegerInSlice(i, afterUpdate) {
			fmt.Println("To Delete ", i, host)
		}
	}
}

func StoreHosts(cfg *models.Config) {
	for _, host := range cfg.Hosts {
		InsertHost(host, cfg)
	}
}
