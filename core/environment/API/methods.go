package API

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Return Smart Proxy ID from foreman
func (Get) SmartProxyID(host string, ctx *user.GlobalCTX) int {

	// VARS
	var ID int
	var responseObj SmartProxies

	// =========
	response, err := utils.ForemanAPI("GET", host, "environments", "", ctx)
	if err != nil {
		ID = -1
	}

	err = json.Unmarshal(response.Body, &responseObj)
	if err != nil {
		ID = -1
	} else {
		ID = responseObj.Results[0].ID
	}

	return ID
}

// Return all environments by host
func (Get) All(host string, ctx *user.GlobalCTX) (Environments, error) {

	// VARS
	var result Environments

	// ======
	response, err := utils.ForemanAPI("GET", host, "environments", "", ctx)
	if err != nil {
		return Environments{}, err
	}

	err = json.Unmarshal(response.Body, &result)
	if err != nil {
		return Environments{}, err
	}

	return result, nil
}

// =====================================================================================================================
// INSERT
// =====================================================================================================================

// TODO: Need to implement importing new puppet classes from code for environments
//func ImportPuppetClasses(p SweUpdateParams, ctx *user.GlobalCTX) {
//	//POST /api/environments/:environment_id/smart_proxies/:id/import_puppetclasses
//	pID := ApiGetSmartProxy(p.Host, ctx)
//	eID := DbForemanID(p.Host, p.Environment, ctx)
//	pApi := SweUpdatePOSTParams{
//		DryRun: p.DryRun,
//		Except: p.Except,
//	}
//	fmt.Println(pID)
//	fmt.Println(eID)
//	fmt.Println(pApi)
//}
