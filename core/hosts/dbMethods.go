package hosts

import (
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// Return all puppet master hosts with environments
func PuppetHosts(ctx *user.GlobalCTX) []ForemanHost {
	var result []ForemanHost
	stmt, err := ctx.Config.Database.DB.Prepare("select host, env from hosts")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query()
	if err != nil {
		utils.Warning.Println(err)
	}
	for rows.Next() {
		var name string
		var env string
		err = rows.Scan(&name, &env)
		if err != nil {
			utils.Error.Println(err)
		}
		if utils.StringInSlice(name, ctx.Config.Hosts) {
			result = append(result, ForemanHost{
				Name: name,
				Env:  env,
			})
		}
	}
	return result
}
