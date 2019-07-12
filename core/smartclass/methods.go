package smartclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, ctx *user.GlobalCTX) {
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes",
		Host:    host,
	}))

	beforeUpdate := GetForemanIDs(host, ctx)
	sort.Ints(beforeUpdate)

	var afterUpdate []int

	smartClassesResult, err := GetAll(host, ctx)
	if err != nil {
		logger.Warning.Printf("Error on getting Smart Classes and Overrides:\n%q", err)
	}

	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Storing Smart classes",
		Host:    host,
	}))

	for _, i := range smartClassesResult {
		afterUpdate = append(afterUpdate, i.ID)
		sort.Ints(afterUpdate)
		InsertSC(host, i, ctx)
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
				DeleteSmartClass(host, i, ctx)
			}
		}
	}
}
