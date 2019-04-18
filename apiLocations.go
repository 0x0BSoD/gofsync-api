package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

// ===============
// GET
// ===============
func locations(host string, cfg *models.Config) (models.Locations, error) {
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

func locationsByHG(host string, hgID int, lastID int64, cfg *models.Config) {
	var result models.Locations
	uri := fmt.Sprintf("hostgroups/%d/locations", hgID)
	bodyText, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, bodyText)
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, bodyText)
	}
	var locList []int64
	for _, loc := range result.Results {
		lId := checkLoc(host, loc.Name, cfg)
		locList = append(locList, int64(lId))
	}
	updateLocInHG(lastID, locList, cfg)
}
