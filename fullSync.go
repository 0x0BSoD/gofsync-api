package main

import (
	"git.ringcentral.com/alexander.simonov/goFsync/logger"
	"sync"
)

// =================================================================
// RUN
// =================================================================
func fullSync() {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			// Locations ===
			locationsResult, err := locations(host)
			if err != nil {
				logger.Warning.Printf("Error on getting Locations:\n%q", err)
			}

			for _, loc := range locationsResult.Results {
				insertToLocations(host, loc.Name, loc.ID)
			}

			// Environments ===
			environmentsResult, err := environments(host)
			if err != nil {
				logger.Warning.Printf("Error on getting Environments:\n%q", err)
			}

			for _, env := range environmentsResult.Results {
				insertToEnvironments(host, env.Name, env.ID)
			}

			// Puppet classes ===
			getAllPCResult, err := getAllPC(host)
			if err != nil {
				logger.Warning.Printf("Error on getting Puppet classes:\n%q", err)
			}

			for className, subClasses := range getAllPCResult {
				for _, subClass := range subClasses {
					insertPC(host, className, subClass.Name, subClass.ID)
				}
			}

			// Smart classes ===
			smartClassesResult, err := smartClasses(host)
			if err != nil {
				logger.Warning.Printf("Error on getting Smart Classes and Overrides:\n%q", err)
			}

			for _, i := range smartClassesResult {
				lastID := insertSC(host, i)
				if lastID != -1 {
					// Getting data by Foreman Smart Class ID
					ovrResult := scOverridesById(host, i.ID)
					for _, ovr := range ovrResult {
						// Storing data by internal SmartClass ID
						insertSCOverride(lastID, ovr, i.ParameterType)
					}
				}
			}

			// Host groups ===
			hg := SWE{}
			hg.Get(host)

			// Match smart classes to puppet class ==
			smartClassByPC(host)

		}(host)
	}
	wg.Wait()
}
