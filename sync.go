package main

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/locations/info"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/jasonlvhit/gocron"
	"sync"
)

// =================================================================
// RUN
// =================================================================
func locSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			locations.Sync(host, ctx)
		}(host)
	}
	wg.Wait()
}

func envSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			environment.Sync(host, ctx)
		}(host)
	}
	wg.Wait()
}

func puppetClassSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			puppetclass.Sync(host, ctx)
		}(host)
	}
	wg.Wait()
}

func smartClassSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			smartclass.Sync(host, ctx)
		}(host)
	}
	wg.Wait()
}

func hostGroupsSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			hostgroups.Sync(host, ctx)
		}(host)
	}
	wg.Wait()
}

func puppetClassUpdate(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			puppetclass.UpdateSCID(host, ctx)
		}(host)
	}
	wg.Wait()
}

func DashboardUpdate(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			info.Sync(host, ctx)
		}(host)
	}
	wg.Wait()
	_, time := gocron.NextRun()
	fmt.Println("Next Run dashboard update: ", time)
}

func fullSync(ctx *user.GlobalCTX) {
	var wg sync.WaitGroup
	for _, host := range globConf.Hosts {

		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			// Socket Broadcast ---
			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "started",
					},
				})
			}

			// ---

			// Locations ===
			//==========================================================================================================
			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "Locations",
						Counter: struct {
							Current int `json:"current"`
							Total   int `json:"total"`
						}{1, 7},
					},
				})
			}
			locations.Sync(host, ctx)

			// Environments ===
			//==========================================================================================================
			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "Environments",
						Counter: struct {
							Current int `json:"current"`
							Total   int `json:"total"`
						}{2, 7},
					},
				})
			}
			environment.Sync(host, ctx)

			// Puppet classes ===
			//==========================================================================================================
			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "Puppet classes",
						Counter: struct {
							Current int `json:"current"`
							Total   int `json:"total"`
						}{3, 7},
					},
				})
			}
			puppetclass.Sync(host, ctx)

			// Smart classes ===
			//==========================================================================================================
			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "Smart classes",
						Counter: struct {
							Current int `json:"current"`
							Total   int `json:"total"`
						}{4, 7},
					},
				})
			}
			smartclass.Sync(host, ctx)

			// Host groups ===
			//==========================================================================================================
			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "Host groups",
						Counter: struct {
							Current int `json:"current"`
							Total   int `json:"total"`
						}{5, 7},
					},
				})
			}
			hostgroups.Sync(host, ctx)

			// Match smart classes to puppet class ==
			//==========================================================================================================
			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "Matching smart classes to puppet class",
						Counter: struct {
							Current int `json:"current"`
							Total   int `json:"total"`
						}{6, 7},
					},
				})
			}
			puppetclass.UpdateSCID(host, ctx)

			// Save to json files
			//==========================================================================================================
			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "Saving to json files",
						Counter: struct {
							Current int `json:"current"`
							Total   int `json:"total"`
						}{7, 7},
					},
				})
			}
			hostgroups.SaveHGToJson(ctx)

			if webServer {
				ctx.Session.SendMsg(models.WSMessage{
					Broadcast: true,
					Operation: "hostUpdate",
					Data: models.Step{
						Host:   host,
						Status: ctx.Session.UserName,
						State:  "done",
					},
				})
			}
		}(host)
	}
	wg.Wait()

	//_, time := gocron.NextRun()
	//fmt.Println("Next Run fullSync: ", time)
}

func startScheduler(ctx *user.GlobalCTX) {
	localCTX := ctx
	//gocron.Every(2).Hours().Do(fullSync, localCTX)
	gocron.Every(5).Minutes().Do(DashboardUpdate, localCTX)
	_, time := gocron.NextRun()
	fmt.Println("Next Run: ", time)
	<-gocron.Start()
}
