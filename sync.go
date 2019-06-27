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
func locSync(ss *models.Session) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			locations.Sync(host, ss)
		}(host)
	}
	wg.Wait()
}

func envSync(ss *models.Session) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			environment.Sync(host, ss)
		}(host)
	}
	wg.Wait()
}

func puppetClassSync(ss *models.Session) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			puppetclass.Sync(host, ss)
		}(host)
	}
	wg.Wait()
}

func smartClassSync(ss *models.Session) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			smartclass.Sync(host, ss)
		}(host)
	}
	wg.Wait()
}

func hostGroupsSync(ss *models.Session) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			hostgroups.Sync(host, ss)
		}(host)
	}
	wg.Wait()
}

func puppetClassUpdate(ss *models.Session) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			puppetclass.UpdateSCID(host, ss)
		}(host)
	}
	wg.Wait()
}

func fullSync(ss *models.Session) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			// Locations ===
			//==========================================================================================================
			locations.Sync(host, ss)

			// Environments ===
			//==========================================================================================================
			environment.Sync(host, ss)

			// Puppet classes ===
			//==========================================================================================================
			puppetclass.Sync(host, ss)

			// Smart classes ===
			//==========================================================================================================
			smartclass.Sync(host, ss)

			// Host groups ===
			//==========================================================================================================
			hostgroups.Sync(host, ss)

			// Match smart classes to puppet class ==
			//==========================================================================================================
			puppetclass.UpdateSCID(host, ss)

		}(host)
	}
	wg.Wait()
}
