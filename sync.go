package main

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/core/environment"
	"git.ringcentral.com/alexander.simonov/goFsync/core/hostgroups"
	"git.ringcentral.com/alexander.simonov/goFsync/core/locations"
	"git.ringcentral.com/alexander.simonov/goFsync/core/puppetclass"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	"sync"
)

// =================================================================
// RUN
// =================================================================
func locSync(cfg *models.Config) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			locations.Sync(host, cfg)
		}(host)
	}
	wg.Wait()
}

func envSync(cfg *models.Config) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			environment.Sync(host, cfg)
		}(host)
	}
	wg.Wait()
}

func puppetClassSync(cfg *models.Config) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			puppetclass.Sync(host, cfg)
		}(host)
	}
	wg.Wait()
}

func smartClassSync(cfg *models.Config) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			smartclass.Sync(host, cfg)
		}(host)
	}
	wg.Wait()
	fmt.Println("DONE")
}

func fullSync(cfg *models.Config) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			// Locations ===
			//==========================================================================================================
			locations.Sync(host, cfg)

			// Environments ===
			//==========================================================================================================
			environment.Sync(host, cfg)

			// Puppet classes ===
			//==========================================================================================================
			puppetclass.Sync(host, cfg)

			// Smart classes ===
			//==========================================================================================================
			smartclass.Sync(host, cfg)

			// Host groups ===
			//==========================================================================================================
			fmt.Println(utils.PrintJsonStep(models.Step{
				Actions: "Filling HostGroups",
				Host:    host,
			}))
			hostgroups.GetHostGroups(host, cfg)

			// Match smart classes to puppet class ==
			//==========================================================================================================
			fmt.Println(utils.PrintJsonStep(models.Step{
				Actions: "Match smart classes to puppet class ID's",
				Host:    host,
			}))
			puppetclass.UpdateSCID(host, cfg)

		}(host)
	}
	wg.Wait()
}
