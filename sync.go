package main

import (
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/models"
	"sync"
)

// =================================================================
// RUN
// =================================================================
func locSync(cfg *models.Config) {
	//var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		//wg.Add(1)
		//go func(host string) {
		//	defer wg.Done()
		locations.Sync(host, cfg)
		//}(host)
	}
	//wg.Wait()
}

func envSync(cfg *models.Config) {
	//var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		//wg.Add(1)
		//go func(host string) {
		//	defer wg.Done()
		environment.Sync(host, cfg)
		//}(host)
	}
	//wg.Wait()
}

func puppetClassSync(cfg *models.Config) {
	//var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		//wg.Add(1)
		//go func(host string) {
		//	defer wg.Done()
		puppetclass.Sync(host, cfg)
		//}(host)
	}
	//wg.Wait()
}

func smartClassSync(cfg *models.Config) {
	//var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		//wg.Add(1)
		//go func(host string) {
		//	defer wg.Done()
		smartclass.Sync(host, cfg)
		//}(host)
	}
	//wg.Wait()
}

func hostGroupsSync(cfg *models.Config) {
	//var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		//wg.Add(1)
		//go func(host string) {
		//	defer wg.Done()
		hostgroups.Sync(host, cfg)
		//}(host)
	}
	//wg.Wait()
}

func puppetClassUpdate(cfg *models.Config) {
	//var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		//wg.Add(1)
		//go func(host string) {
		//	defer wg.Done()
		puppetclass.Update(host, cfg)
		//}(host)
	}
	//wg.Wait()
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
			hostgroups.Sync(host, cfg)

			// Match smart classes to puppet class ==
			//==========================================================================================================
			puppetclass.Update(host, cfg)

		}(host)
	}
	wg.Wait()
}
