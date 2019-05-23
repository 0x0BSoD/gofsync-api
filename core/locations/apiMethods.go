package locations

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/models"
	logger "git.ringcentral.com/archops/goFsync/utils"
)

// ===============
// GET
// ===============
func ApiAll(host string, cfg *models.Config) (models.Locations, error) {
	var result models.Locations
	bodyText, err := logger.ForemanAPI("GET", host, "locations", "", cfg)
	if err != nil {
		return models.Locations{}, err
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return models.Locations{}, err
	}
	return result, nil
}
