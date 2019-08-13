package smartclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/smartclass/API"
	"git.ringcentral.com/archops/goFsync/core/smartclass/DB"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, ctx *user.GlobalCTX) {

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes",
		Host:    host,
	}))

	// VARS
	var gAPI API.Get
	var gDB DB.Get
	var iDB DB.Insert
	var dDB DB.Delete

	// ====
	beforeUpdate := gDB.ForemanIDs(host, ctx)
	sort.Ints(beforeUpdate)

	var afterUpdate []int

	smartClassesResult, err := gAPI.All(host, ctx)
	if err != nil {
		utils.Warning.Printf("Error on getting Smart Classes and Overrides:\n%q", err)
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Storing Smart classes",
		Host:    host,
	}))

	for _, i := range smartClassesResult {
		afterUpdate = append(afterUpdate, i.ForemanID)
		sort.Ints(afterUpdate)
		iDB.Add(host, i, ctx)
	}

	for _, i := range beforeUpdate {
		if len(beforeUpdate) != len(afterUpdate) {

			fmt.Println(utils.PrintJsonStep(models.Step{
				Actions: fmt.Sprintf("Checking ...%d", i),
				Host:    host,
			}))

			if !utils.Search(afterUpdate, i) {

				fmt.Println(utils.PrintJsonStep(models.Step{
					Actions: fmt.Sprintf("Deleting ...%d", i),
					Host:    host,
				}))

				err = dDB.SmartClass(host, i, ctx)
				if err != nil {
					utils.Warning.Printf("Error on while deleteing Smart Class:\n%q", err)
				}
			}
		}
	}
}

func SmartClassInList(a string, list []DB.SmartClass) bool {
	for _, b := range list {
		if b.Name == a {
			return true
		}
	}
	return false
}
