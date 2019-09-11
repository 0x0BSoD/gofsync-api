package global

import (
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
)

func Sync(host string, ctx *user.GlobalCTX) {

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:   host,
			Status: ctx.Session.UserName,
			State:  "started",
		},
	})
	// ---

	// Locations ===
	//==========================================================================================================
	locations.Sync(host, ctx)

	// Environments ===
	//==========================================================================================================
	environment.Sync(host, ctx)

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

	// Socket Broadcast ---
	ctx.Session.SendMsg(models.WSMessage{
		Broadcast: true,
		Operation: "hostUpdate",
		Data: models.Step{
			Host:   host,
			Status: ctx.Session.UserName,
			State:  "done",
		},
	})
	// ---

}
