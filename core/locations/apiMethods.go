package locations

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// ===============
// GET
// ===============
func ApiAll(host string, ctx *user.GlobalCTX) (Locations, error) {
	var result Locations
	bodyText, err := utils.ForemanAPI("GET", host, "locations", "", ctx)
	if err != nil {
		return Locations{}, err
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return Locations{}, err
	}
	return result, nil
}
