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
		scId := InsertSC(host, i, cfg)
		if i.OverrideValuesCount > 0 {
			fmt.Println("SmartClass ID: ", scId)
			beforeUpdateOvr := GetForemanIDsBySCid(scId, cfg)
			for _, ovr := range i.OverrideValues {
				fmt.Println(ovr)
				afterUpdateOvr = append(afterUpdateOvr, ovr.ID)
				InsertSCOverride(scId, ovr, i.ParameterType, cfg)
			}

			sort.Ints(afterUpdateOvr)

			for _, j := range beforeUpdateOvr {
				if !utils.IntegerInSlice(j, afterUpdateOvr) {
					fmt.Println("SC: ", string(i.PuppetClass.Name))

					fmt.Println("Before: ", len(beforeUpdateOvr))
					fmt.Println("After: ", len(afterUpdateOvr))
					fmt.Println("+=+=+==+=+++====+====+")

					DeleteOverride(scId, j, cfg)
				}
			}
			afterUpdateOvr = nil
			//	fmt.Println("========")
		}
	}

	for i := range beforeUpdate {
		if len(beforeUpdate) != len(afterUpdate) {
			if !utils.IntegerInSlice(i, afterUpdate) {
				fmt.Println("SC: ", string(i))

				fmt.Println("Before: ", len(beforeUpdate))
				fmt.Println("After: ", len(afterUpdate))
				fmt.Println("+=+=+==+=+++====+====+")
				DeleteSmartClass(host, i, cfg)
			}
		}
	}
}
