package locations

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

func LocSync(host string, cfg *models.Config) {
	// Locations ===
	//==========================================================================================================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Locations",
		Host:    host,
	}))

	beforeUpdate := GetAllLocNames(host, cfg)
	var afterUpdate []string

	locationsResult, err := Locations(host, cfg)
	if err != nil {
		logger.Warning.Printf("Error on getting Locations:\n%q", err)
	}
	for _, loc := range locationsResult.Results {
		InsertToLocations(host, loc.Name, loc.ID, cfg)
		afterUpdate = append(afterUpdate, loc.Name)
	}

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i, afterUpdate) {
			DeleteLocation(host, i, cfg)
		}
	}
}
