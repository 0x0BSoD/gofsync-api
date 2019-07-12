package environment

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// ===============
// GET
// ===============
func ApiAll(host string, ctx *user.GlobalCTX) (Environments, error) {

	var result Environments
	bodyText, err := utils.ForemanAPI("GET", host, "environments", "", ctx)
	if err != nil {
		return Environments{}, err
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return Environments{}, err
	}

	return result, nil
}

func ApiGetSmartProxy(host string, ctx *user.GlobalCTX) int {
	var result SmartProxies
	bodyText, err := utils.ForemanAPI("GET", host, "environments", "", ctx)
	if err != nil {
		return -1
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return -1
	}

	return result.Results[0].ID
}

// ===============
// POST
// ===============
func ImportPuppetClasses(p SweUpdateParams, ctx *user.GlobalCTX) {
	//POST /api/environments/:environment_id/smart_proxies/:id/import_puppetclasses
	pID := ApiGetSmartProxy(p.Host, ctx)
	eID := DbForemanID(p.Host, p.Environment, ctx)
	pApi := SweUpdatePOSTParams{
		DryRun: p.DryRun,
		Except: p.Except,
	}
	fmt.Println(pID)
	fmt.Println(eID)
	fmt.Println(pApi)
}
