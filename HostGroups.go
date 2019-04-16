package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

// HostGroupBase Structure for post
type HostGroupBase struct {
	DBHGExist      int    `json:"dbhg_exist"`
	ExistId        int    `json:"exist_id"`
	ParentId       int    `json:"parent_id"`
	Name           string `json:"name"`
	EnvironmentId  int    `json:"environment_id"`
	PuppetclassIds []int  `json:"puppetclass_ids"`
	LocationIds    []int  `json:"location_ids"`
}
type HostGroupOverrides struct {
	OvrForemanId int    `json:"ovr_foreman_id"`
	ScForemanId  int    `json:"sc_foreman_id"`
	Match        string `json:"match"`
	Value        string `json:"value"`
}
type HWPostRes struct {
	BaseInfo   HostGroupBase        `json:"hostgroup"`
	Overrides  []HostGroupOverrides `json:"override_value"`
	NotExistPC []int                `json:"not_exist_pc"`
}

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
func postHG(sHost string, tHost string, hgId int) (HWPostRes, error) {

	// Source Host Group
	hostGroupData := getHG(hgId)

	// Step 1. Check if Host Group exist on the host
	// --> we trust frontend that <--
	hostGroupExistBase := checkHG(hostGroupData.Name, tHost)
	tmp := hostGroupCheck(tHost, hostGroupData.Name)
	hostGroupExist := tmp.ID
	//if hostGroupExist != -1 {
	//	log.Fatalf("Host Group '%s' already exist on %s", hostGroupData.Name, tHost)
	//}

	// Step 2. Check Environment exist on the target host
	environmentExist := checkPostEnv(tHost, hostGroupData.Environment)
	if environmentExist == -1 {
		return HWPostRes{}, errors.New(fmt.Sprintf("Environment '%s' not exist on %s", hostGroupData.Environment, tHost))
	}

	// Step 3. Get parent Host Group ID on target host
	parentHGId := checkHGID("SWE", tHost)
	if parentHGId == -1 {
		return HWPostRes{}, errors.New(fmt.Sprintf("Parent Host Group 'SWE' not exist on %s", tHost))
	}

	// Step 4. Get all locations for the target host
	locationsIds := getAllLocations(tHost)

	// Step 5. Check Puppet Classes on existing on the target host
	// and
	// Step 6. Get Smart Class data
	var PuppetClassesIds []int
	var SCOverrides []HostGroupOverrides
	for _, i := range hostGroupData.PuppetClasses {
		// Get Puppet Classes IDs for target Foreman
		for _, subclass := range i {
			targetPCData := getByNamePC(subclass.Subclass, tHost)
			//sourcePCData := getByNamePC(subclass.Subclass, sHost)

			// If we not have Puppet Class for target host
			if targetPCData.ID == 0 {
				//return HWPostRes{}, errors.New(fmt.Sprintf("Puppet Class '%s' not exist on %s", name, tHost))
			} else {

				// Build Target PC id's and SmartClasses
				PuppetClassesIds = append(PuppetClassesIds, targetPCData.ForemanId)
				var sourceScDataSet []SCGetResAdv
				for _, pc := range hostGroupData.PuppetClasses {
					for _, subPc := range pc {
						for _, scName := range subPc.SmartClasses {
							// Get Smart Class data
							sourceScData := getSC(sHost, subclass.Subclass, scName)
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
					for _, scId := range Integers(targetPCData.SCIDs) {
						targetSC := getSCData(scId)
						for _, sourceSC := range sourceScDataSet {
							if sourceSC.Name == targetSC.Name {
								srcOvr, _ := getOvrData(sourceSC.ID, hostGroupData.Name, targetSC.Name)
								targetOvr, trgErr := getOvrData(targetSC.ID, hostGroupData.Name, targetSC.Name)
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

									SCOverrides = append(SCOverrides, HostGroupOverrides{
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

	return HWPostRes{
		BaseInfo: HostGroupBase{
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

func postCheckHG(tHost string, hgId int) bool {
	// Source Host Group
	hostGroupData := getHG(hgId)
	// Step 1. Check if Host Group exist on the host
	hostGroupExist := checkHG(hostGroupData.Name, tHost)
	res := false
	if hostGroupExist != -1 {
		res = true
	}

	return res
}

func saveHGToJson() {
	for _, host := range globConf.Hosts {
		data := getHGList(host)
		for _, d := range data {
			hgData := getHG(d.ID)
			rJson, _ := json.MarshalIndent(hgData, "", "    ")
			path := fmt.Sprintf("/opt/goFsync/HG/%s/%s.json", host, hgData.Name)
			if _, err := os.Stat("HG/" + host); os.IsNotExist(err) {
				err = os.Mkdir("HG/"+host, 0777)
				if err != nil {
					log.Fatalf("Error on mkdir: %s", err)
				}
			}
			//fmt.Println("Storing to: ", path)
			err := ioutil.WriteFile(path, rJson, 0644)
			if err != nil {
				log.Fatalf("Error on writing file: %s", err)
			}

		}
	}

	var out bytes.Buffer
	commitMessage := fmt.Sprintf("Auto commit. Date: %s", time.Now())

	cmd := exec.Command("bash", "HG/lazygit.sh", commitMessage)
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println(out.String())

}
