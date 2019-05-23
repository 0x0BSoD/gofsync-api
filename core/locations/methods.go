package locations

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, cfg *models.Config) {
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Locations",
		Host:    host,
	}))

	beforeUpdate := DbAll(host, cfg)
	var afterUpdate []string

	locationsResult, err := ApiAll(host, cfg)
	if err != nil {
		logger.Warning.Printf("Error on getting Locations:\n%q", err)
	}

	sort.Slice(locationsResult.Results, func(i, j int) bool {
		return locationsResult.Results[i].ID < locationsResult.Results[j].ID
	})

	for _, loc := range locationsResult.Results {
		DbInsert(host, loc.Name, loc.ID, cfg)
		afterUpdate = append(afterUpdate, loc.Name)
	}
	sort.Strings(afterUpdate)

	for _, i := range beforeUpdate {
		if !utils.StringInSlice(i, afterUpdate) {
			DbDelete(host, i, cfg)
		}
	}
}
