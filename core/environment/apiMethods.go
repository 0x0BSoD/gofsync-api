package environment

import (
	"encoding/json"
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
