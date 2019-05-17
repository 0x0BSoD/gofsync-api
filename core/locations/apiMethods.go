package locations

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
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
