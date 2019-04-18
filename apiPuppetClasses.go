package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

// ===============
// GET
// ===============
// Get all Puppet Classes and insert to base
func getAllPC(host string, cfg *models.Config) (map[string][]models.PuppetClass, error) {

	var pcResult models.PuppetClasses
	var result = make(map[string][]models.PuppetClass)

	// check items count
	uri := fmt.Sprintf("puppetclasses?format=json&per_page=%d", cfg.Api.GetPerPage)
	bodyText, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(bodyText.Body, &pcResult)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, bodyText)
		}

		if pcResult.Total > cfg.Api.GetPerPage {
			pagesRange := utils.Pager(pcResult.Total, cfg.Api.GetPerPage)
			for i := 1; i <= pagesRange; i++ {
				uri := fmt.Sprintf("puppetclasses?format=json&page=%d&per_page=%d", i, cfg.Api.GetPerPage)
				bodyText, err := logger.ForemanAPI("GET", host, uri, "", cfg)
				if err == nil {
					err := json.Unmarshal(bodyText.Body, &pcResult)
					if err != nil {
						return result, err
					}

					for className, cl := range pcResult.Results {
						for _, subClass := range cl {
							result[className] = append(result[className], subClass)
						}
					}
				}
			}
		} else {
			for className, cl := range pcResult.Results {
				for _, subClass := range cl {
					result[className] = append(result[className], subClass)
				}
			}
		}
	}

	return result, nil
}

// Get Puppet Classes by host group and insert to Host Group
func getPCByHg(host string, hgID int, bdId int64, cfg *models.Config) []int {
	var result models.PuppetClasses
	var foremanSCIds []int

	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	bodyText, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(bodyText.Body, &result)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, bodyText)
		}
		var pcIDs []int64
		for className, cl := range result.Results {
			for _, subclass := range cl {
				foremanSCIds = append(foremanSCIds, subclass.ID)
				lastId := insertPC(host, className, subclass.Name, subclass.ID, cfg)
				if lastId != -1 {
					pcIDs = append(pcIDs, lastId)
				}
			}
		}
		updatePCinHG(bdId, pcIDs, cfg)
	}
	return foremanSCIds
}

// Just get Puppet Classes by host group
func getPCByHgJson(host string, hgID int, cfg *models.Config) map[string][]models.PuppetClass {

	var result models.PuppetClasses
	//var pcIDs []int64

	uri := fmt.Sprintf("hostgroups/%d/puppetclasses", hgID)
	bodyText, err := logger.ForemanAPI("GET", host, uri, "", cfg)
	if err == nil {
		err := json.Unmarshal(bodyText.Body, &result)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, bodyText)
		}
	} else {
		logger.Warning.Printf("%q: getPCByHgJson", err)
	}
	return result.Results
}
