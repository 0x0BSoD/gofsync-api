package locations

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/locations/API"
	"git.ringcentral.com/archops/goFsync/core/locations/DB"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, ctx *user.GlobalCTX) {

	// Step LOG to stdout ======================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Locations",
		Host:    host,
	}))
	// =========================================

	// VARS
	var (
		gDB  DB.Get
		iDB  DB.Insert
		dDB  DB.Delete
		gAPI API.Get
	)

	// from DB
	beforeUpdate, _ := gDB.All(host, ctx)
	var afterUpdate []string

	// from foreman
	locationsResult, err := gAPI.All(host, ctx)
	if err != nil {
		logger.Warning.Printf("Error on getting Locations:\n%q", err)
		utils.GetErrorContext(err)
	}
	sort.Slice(locationsResult.Results, func(i, j int) bool {
		return locationsResult.Results[i].ID < locationsResult.Results[j].ID
	})

	// store
	for _, loc := range locationsResult.Results {
		iDB.Add(host, loc.Name, loc.ID, ctx)
		afterUpdate = append(afterUpdate, loc.Name)
	}
	sort.Strings(afterUpdate)

	// delete if don't have any errors
	if err == nil {
		for _, i := range beforeUpdate {
			if !utils.StringInSlice(i, afterUpdate) {
				dDB.ByName(host, i, ctx)
			}
		}
	}
}
