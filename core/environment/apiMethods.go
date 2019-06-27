package environment

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
)

// ===============
// GET
// ===============
func ApiAll(host string, s *models.Session) (models.Environments, error) {

	var result models.Environments
	bodyText, err := utils.ForemanAPI("GET", host, "environments", "", s.Config)
	if err != nil {
		return models.Environments{}, err
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return models.Environments{}, err
	}

	return result, nil
}
