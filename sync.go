package main

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/foremans"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"github.com/jasonlvhit/gocron"
	"sync"
)

// =================================================================
// RUN
// =================================================================
func locSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for hostname := range globConf.Hosts {

		wg.Add(1)
		go func(hostname string) {
			defer wg.Done()
			locations.Sync(hostname, ctx)
		}(hostname)
	}
	wg.Wait()
}

func envSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for hostname := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			err := environment.Sync(host, ctx)
			if err != nil {
				panic(err)
			}
		}(hostname)
	}
	wg.Wait()
}

func puppetClassSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for hostname := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			puppetclass.Sync(host, ctx)
		}(hostname)
	}
	wg.Wait()
}

func smartClassSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for hostname := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			smartclass.Sync(host, ctx)
		}(hostname)
	}
	wg.Wait()
}

func hostGroupsSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for hostname := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			hostgroups.Sync(host, ctx)
		}(hostname)
	}
	wg.Wait()
}

func puppetClassUpdate(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for hostname := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			puppetclass.UpdateSCID(host, ctx)
		}(hostname)
	}
	wg.Wait()
}

func DashboardUpdate(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for hostname := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			foremans.Sync(host, ctx)
		}(hostname)
	}
	wg.Wait()
	_, time := gocron.NextRun()
	fmt.Println("Next Run dashboard update: ", time)
}

func fullSync(ctx *user.GlobalCTX) {
	if !ctx.SyncWIP {
		ctx.SyncWIP = true
		var wg sync.WaitGroup
		for hostname := range globConf.Hosts {

			wg.Add(1)
			go func(host string) {
				defer wg.Done()

				// Locations ===
				//==========================================================================================================
				locations.Sync(host, ctx)

				// Environments ===
				//=========================================================================================================
				err := environment.Sync(host, ctx)
				if err != nil {
					panic(err)
				}

				// Puppet classes ===
				//==========================================================================================================
				puppetclass.Sync(host, ctx)

				// Smart classes ===
				//==========================================================================================================
				smartclass.Sync(host, ctx)

				// Host groups ===
				//==========================================================================================================
				hostgroups.Sync(host, ctx)

				// Match smart classes to puppet class ==
				//==========================================================================================================
				puppetclass.UpdateSCID(host, ctx)

				// Save to json files
				//==========================================================================================================
				hostgroups.SaveHGToJson(ctx)
			}(hostname)
		}
		wg.Wait()
		ctx.SyncWIP = false
	}
}

func startScheduler(ctx *user.GlobalCTX) {
	localCTX := ctx
	gocron.Every(1).Day().At("21:30").Do(fullSync, localCTX)
	gocron.Every(5).Minutes().Do(DashboardUpdate, localCTX)
	_, time := gocron.NextRun()
	fmt.Println("Next Run: ", time)
	<-gocron.Start()
}
