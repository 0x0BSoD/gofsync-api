package environment

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
)

// ===============
// GET
// ===============
func Environments(host string, cfg *models.Config) (models.Environments, error) {

	var result models.Environments
	bodyText, err := utils.ForemanAPI("GET", host, "environments", "", cfg)
	if err != nil {
		return models.Environments{}, err
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return models.Environments{}, err
	}

	return result, nil
}
