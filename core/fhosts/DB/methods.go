package DB

import (
	"database/sql"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Return DB ID for puppet master host parameter
// Here user *sql.DB because ctx  is not ready yet
func (Get) ID(host string, db *sql.DB) int {

	// VARS
	var ID int

	// =========
	stmt, err := db.Prepare("select id from hosts where host=?")
	if err != nil {
		utils.Warning.Println(err)
		ID = -1
	}

	defer utils.DeferCloseStmt(stmt)
	err = stmt.QueryRow(host).Scan(&ID)
	if err != nil {
		utils.Warning.Println(err)
		ID = -1
	}

	return ID
}

// Return all puppet master hosts with environments
func (Get) All(ctx *user.GlobalCTX) []ForemanHost {

	// VARS
	var result []ForemanHost

	// ==========
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

// Return Environment for puppet master host
func (Get) Environment(host string, ctx *user.GlobalCTX) string {

	// VARS
	var hostEnv string

	// =========
	stmt, err := ctx.Config.Database.DB.Prepare("select env from hosts where host=?")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host).Scan(&hostEnv)
	if err != nil {
		utils.Warning.Println(err)
		return ""
	}

	return hostEnv
}

// =====================================================================================================================
// INSERT
// =====================================================================================================================

// Insert puppet master host
func (Insert) Add(host string, db *sql.DB) int {

	// VARS
	var gDB Get
	ID := gDB.ID(host, db)

	// =====
	if ID == -1 {
		stmt, err := db.Prepare("insert into hosts (host) values(?)")
		if err != nil {
			utils.Warning.Println(err)
			return -1
		}
		defer utils.DeferCloseStmt(stmt)
		response, err := stmt.Exec(host)
		if err != nil {
			utils.Warning.Println(err)
			return -1
		}
		ID64, _ := response.LastInsertId()
		ID = int(ID64)
	}

	return ID
}
