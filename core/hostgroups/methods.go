package hostgroups

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/core/environment"
	"git.ringcentral.com/alexander.simonov/goFsync/core/locations"
	"git.ringcentral.com/alexander.simonov/goFsync/core/puppetclass"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
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
	hostGroupData := GetHG(hgId, cfg)

	// Step 1. Check if Host Group exist on the host
	// --> we trust frontend that <--
	hostGroupExistBase := CheckHG(hostGroupData.Name, tHost, cfg)
	tmp := HostGroupCheck(tHost, hostGroupData.Name, cfg)
	hostGroupExist := tmp.ID
	//if hostGroupExist != -1 {
	//	log.Fatalf("Host Group '%s' already exist on %s", hostGroupData.Name, tHost)
	//}

	// Step 2. Check Environment exist on the target host
	environmentExist := environment.CheckPostEnv(tHost, hostGroupData.Environment, cfg)
	if environmentExist == -1 {
		return models.HWPostRes{}, errors.New(fmt.Sprintf("Environment '%s' not exist on %s", hostGroupData.Environment, tHost))
	}

	// Step 3. Get parent Host Group ID on target host
	parentHGId := CheckHGID("SWE", tHost, cfg)
	if parentHGId == -1 {
		return models.HWPostRes{}, errors.New(fmt.Sprintf("Parent Host Group 'SWE' not exist on %s", tHost))
	}

	// Step 4. Get all locations for the target host
	locationsIds := locations.GetAllLocations(tHost, cfg)

	// Step 5. Check Puppet Classes on existing on the target host
	// and
	// Step 6. Get Smart Class data
	var PuppetClassesIds []int
	var SCOverrides []models.HostGroupOverrides
	for _, i := range hostGroupData.PuppetClasses {
		// Get Puppet Classes IDs for target Foreman
		for _, subclass := range i {
			targetPCData := puppetclass.GetByNamePC(subclass.Subclass, tHost, cfg)
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
						for _, scName := range subPc.SmartClasses {
							// Get Smart Class data
							sourceScData := smartclass.GetSC(sHost, subclass.Subclass, scName, cfg)
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
						for _, sourceSC := range sourceScDataSet {
							if sourceSC.Name == targetSC.Name {
								srcOvr, _ := smartclass.GetOvrData(sourceSC.ID, hostGroupData.Name, targetSC.Name, cfg)
								targetOvr, trgErr := smartclass.GetOvrData(targetSC.ID, hostGroupData.Name, targetSC.Name, cfg)
								if srcOvr.SmartClassId != 0 {

									OverrideID := -1

									if trgErr == nil {
										OverrideID = targetOvr.OverrideId
									}

									fmt.Println("Match: ", srcOvr.Match)
									fmt.Println("Value: ", srcOvr.Value)
									fmt.Println("OverrideId: ", OverrideID)
									fmt.Println("Parameter: ", srcOvr.Parameter)
									fmt.Println("SmartClassId: ", targetSC.ForemanId)
									fmt.Println("============================================")

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
			DBHGExist:      hostGroupExistBase,
			Name:           hostGroupData.Name,
			ParentId:       parentHGId,
			ExistId:        hostGroupExist,
			EnvironmentId:  environmentExist,
			LocationIds:    locationsIds,
			PuppetclassIds: PuppetClassesIds,
		},
		Overrides: SCOverrides,
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
			//fmt.Println("Storing to: ", path)
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
