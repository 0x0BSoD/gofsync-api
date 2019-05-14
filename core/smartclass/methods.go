package smartclass

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"sort"
)

func Sync(host string, cfg *models.Config) {
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes",
		Host:    host,
	}))

	beforeUpdate := GetForemanIDs(host, cfg)
	sort.Ints(beforeUpdate)

	var afterUpdate []int
	var afterUpdateOvr []int

	smartClassesResult, err := GetAll(host, cfg)
	if err != nil {
		logger.Warning.Printf("Error on getting Smart Classes and Overrides:\n%q", err)
	}

	for _, i := range smartClassesResult {
		afterUpdate = append(afterUpdate, i.ID)
		sort.Ints(afterUpdate)

		lastID := InsertSC(host, i, cfg)
		if lastID != -1 {
			beforeUpdateOvr := GetOverrodesForemanIDs(int(lastID), cfg)
			afterUpdateOvr = []int{}
			sort.Ints(beforeUpdateOvr)
			// Getting data by Foreman Smart Class ID
			ovrResult := SCOverridesById(host, i.ID, cfg)
			for _, ovr := range ovrResult {
				// Storing data by internal SmartClass ID
				afterUpdateOvr = append(afterUpdateOvr, ovr.ID)
				InsertSCOverride(lastID, ovr, i.ParameterType, cfg)
			}
			sort.Ints(afterUpdateOvr)
			for _, i := range beforeUpdateOvr {
				if !utils.IntegerInSlice(i, afterUpdateOvr) {
					DeleteOverride(int(lastID), i, cfg)
				}
			}
		}
	}

	fmt.Println("SC Before: ", len(beforeUpdate))
	fmt.Println("SC After: ", len(afterUpdate))

	for i := range beforeUpdate {
		if len(beforeUpdate) != len(afterUpdate) {
			if !utils.IntegerInSlice(i, afterUpdate) {
				fmt.Println(beforeUpdate[i])
				DeleteSmartClass(host, i, cfg)
			}
		}
	}
}
