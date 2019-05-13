package smartclass

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

func Sync(host string, cfg *models.Config) {
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes",
		Host:    host,
	}))

	beforeUpdate := GetForemanIDs(host, cfg)
	var afterUpdate []int

	smartClassesResult, err := GetAll(host, cfg)
	if err != nil {
		logger.Warning.Printf("Error on getting Smart Classes and Overrides:\n%q", err)
	}

	for _, i := range smartClassesResult {
		afterUpdate = append(afterUpdate, i.ID)
		lastID := InsertSC(host, i, cfg)
		if lastID != -1 {
			// Getting data by Foreman Smart Class ID
			ovrResult := SCOverridesById(host, i.ID, cfg)
			for _, ovr := range ovrResult {
				// Storing data by internal SmartClass ID
				InsertSCOverride(lastID, ovr, i.ParameterType, cfg)
			}
		}
	}

	for _, i := range beforeUpdate {
		if !utils.IntegerInSlice(i, afterUpdate) {
			DeleteSmartClass(host, i, cfg)
		}
	}
}
