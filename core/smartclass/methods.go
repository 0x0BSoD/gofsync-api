package smartclass

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

func Sync(host string, ss *models.Session) {
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Smart classes",
		Host:    host,
	}))

	beforeUpdate := GetForemanIDs(host, ss)
	sort.Ints(beforeUpdate)

	var afterUpdate []int

	smartClassesResult, err := GetAll(host, ss)
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
		InsertSC(host, i, ss)
	}

	//fmt.Println(utils.PrintJsonStep(models.Step{
	//	Actions: "Deleting Smart classes",
	//	Host:    host,
	//}))
	//for i := range beforeUpdate {
	//	if len(beforeUpdate) != len(afterUpdate) {
	//		if !utils.Search(afterUpdate, i) {
	//			DeleteSmartClass(host, i, cfg)
	//		}
	//	}
	//}
}
