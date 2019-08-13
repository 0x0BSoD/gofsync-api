package hostgroup

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/hostgroup/API"
	"git.ringcentral.com/archops/goFsync/core/hostgroup/DB"
	API2 "git.ringcentral.com/archops/goFsync/core/puppetclass/API"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
)

func Sync(host string, ctx *user.GlobalCTX) {

	//==========================================================================================================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Filling HostGroups",
		Host:    host,
	}))

	// VARS
	var gDB DB.Get
	var iDB DB.Insert
	var dDB DB.Delete
	var gAPI API.Get
	var iAPIPuppetClass API2.Insert

	// Socket Broadcast ---
	if ctx.Session.PumpStarted {
		data := models.Step{
			Host:    host,
			Actions: "Getting HostGroups",
			State:   "",
		}
		msg, _ := json.Marshal(data)
		ctx.Session.SendMsg(msg)
	}
	// ---

	beforeUpdate := gDB.ForemanIDs(host, ctx)
	var afterUpdate []int

	results := gAPI.All(host, ctx)

	// RT SWEs =================================================================================================
	SWEs := RackTablesData(gDB.HostEnvironment(host, ctx), ctx)

	for idx, i := range results {

		// Socket Broadcast ---
		if ctx.Session.PumpStarted {
			data := models.Step{
				Host:    host,
				Actions: "Saving HostGroups",
				State:   fmt.Sprintf("HostGroup: %s %d/%d", i.Name, idx+1, len(results)),
			}
			msg, _ := json.Marshal(data)
			ctx.Session.SendMsg(msg)
		}
		// ---

		sJson, _ := json.Marshal(i)

		sweStatus := GetFromRT(i.Name, SWEs)
		fmt.Printf("Get: %s\tStatus:%s\n", i.Name, sweStatus)

		// Add Host Group to base
		lastId := iDB.Add(i.Name, host, string(sJson), sweStatus, i.ForemanID, ctx)
		afterUpdate = append(afterUpdate, i.ForemanID)

		// If success
		if lastId != -1 {
			// Get Puppet classes and smart classes and store in DB
			iAPIPuppetClass.Add(host, i.ForemanID, lastId, ctx)

			// Get host group parameters
			params := gAPI.Parameters(host, lastId, i.ForemanID, ctx)
			for _, p := range params {
				iDB.Parameter(lastId, p, ctx)
			}
		}
	}

	for _, i := range beforeUpdate {
		if !utils.Search(afterUpdate, i) {
			dDB.ByID(i, host, ctx)
		}
	}
}

func RackTablesData(env string, ctx *user.GlobalCTX) map[string]string {
	// RT SWEs =================================================================================================
	var swes []RackTablesSWE
	result := make(map[string]string)

	var rtHost string
	if env == "stage" {
		rtHost = ctx.Config.RackTables.Stage
	} else if env == "prod" {
		rtHost = ctx.Config.RackTables.Production
	}

	body, err := utils.RackTablesAPI("GET", rtHost, "rchwswelookups/search?q=name~.*&fields=name,swestatus&format=json", "", ctx)
	if err != nil {
		utils.Error.Println(err)
	}
	err = json.Unmarshal(body.Body, &swes)
	if err != nil {
		utils.Warning.Printf("%q:\n %s\n", err, body.Body)
	}
	// RT SWEs =================================================================================================
	for _, swe := range swes {
		result[swe.Name] = swe.SweStatus
	}
	return result
}

func GetFromRT(name string, SWEs map[string]string) string {
	if val, ok := SWEs[name]; ok {
		return val
	}
	return "nope"
}
