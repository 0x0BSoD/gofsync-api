package API

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

func (Get) All(host string, ctx *user.GlobalCTX) (Locations, error) {

	// VARS
	var result Locations

	// =======
	response, err := utils.ForemanAPI("GET", host, "locations", "", ctx)
	if err != nil {
		return Locations{}, err
	}

	err = json.Unmarshal(response.Body, &result)
	if err != nil {
		return Locations{}, err
	}

	return result, nil
}
