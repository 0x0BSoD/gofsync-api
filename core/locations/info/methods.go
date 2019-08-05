package info

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
)

func Sync(host string, ctx *user.GlobalCTX) {

	// Step LOG to stdout ======================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Locations Info",
		Host:    host,
	}))
	// =========================================

	// from foreman
	locationsResult := ApiReportsDaily(host, ctx)
	Update(host, locationsResult, ctx)

	// Socket Broadcast ---
	data := models.Step{
		Host:    host,
		Actions: "Dashboard updated",
	}
	msg, _ := json.Marshal(data)
	ctx.Broadcast(msg)
	// ---

}
