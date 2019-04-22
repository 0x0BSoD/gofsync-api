package main

import (
	"git.ringcentral.com/alexander.simonov/goFsync/core/environment"
	"git.ringcentral.com/alexander.simonov/goFsync/core/hostgroups"
	"git.ringcentral.com/alexander.simonov/goFsync/core/locations"
	"git.ringcentral.com/alexander.simonov/goFsync/core/puppetclass"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"sync"
)

// =================================================================
// RUN
// =================================================================
func fullSync(cfg *models.Config) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			// Locations ===
			locationsResult, err := locations.Locations(host, cfg)
			if err != nil {
				logger.Warning.Printf("Error on getting Locations:\n%q", err)
			}
			for _, loc := range locationsResult.Results {
				locations.InsertToLocations(host, loc.Name, loc.ID, cfg)
			}

			// Environments ===
			environmentsResult, err := environment.Environments(host, cfg)
			if err != nil {
				logger.Warning.Printf("Error on getting Environments:\n%q", err)
			}
			for _, env := range environmentsResult.Results {
				environment.InsertToEnvironments(host, env.Name, env.ID, cfg)
			}

			// Puppet classes ===
			getAllPCResult, err := puppetclass.GetAllPC(host, cfg)
			if err != nil {
				logger.Warning.Printf("Error on getting Puppet classes:\n%q", err)
			}
			for className, subClasses := range getAllPCResult {
				for _, subClass := range subClasses {
					puppetclass.InsertPC(host, className, subClass.Name, subClass.ID, cfg)
				}
			}

			// Smart classes ===
			smartClassesResult, err := smartclass.GetAll(host, cfg)
			if err != nil {
				logger.Warning.Printf("Error on getting Smart Classes and Overrides:\n%q", err)
			}
			for _, i := range smartClassesResult {
				lastID := smartclass.InsertSC(host, i, cfg)
				if lastID != -1 {
					// Getting data by Foreman Smart Class ID
					ovrResult := smartclass.SCOverridesById(host, i.ID, cfg)
					for _, ovr := range ovrResult {
						// Storing data by internal SmartClass ID
						smartclass.InsertSCOverride(lastID, ovr, i.ParameterType, cfg)
					}
				}
			}

			// Host groups ===
			hostgroups.GetHostGroups(host, cfg)

			// Match smart classes to puppet class ==
			puppetclass.UpdateSCID(host, cfg)

		}(host)
	}
	wg.Wait()
}
