package environment

import (
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
)

// ======================================================
// CHECKS and GETS
// ======================================================
func DbID(host string, env string, ctx *user.GlobalCTX) int {

	var id int

	stmt, err := ctx.Config.Database.DB.Prepare("select id from environments where host=? and env=?")
	if err != nil {
		logger.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, env).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}
func DbForemanID(host string, env string, ctx *user.GlobalCTX) int {

	var id int

	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id from environments where host=? and env=?")
	if err != nil {
		logger.Warning.Printf("%q, checkPostEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, env).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func DbAll(host string, ctx *user.GlobalCTX) []string {

	var list []string

	stmt, err := ctx.Config.Database.DB.Prepare("select env from environments where host=?")
	if err != nil {
		logger.Warning.Printf("%q, getEnvList", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		return list
	}

	for rows.Next() {
		var env string
		err = rows.Scan(&env)
		if err != nil {
			logger.Error.Printf("%q, getEnvList", err)
		}
		list = append(list, env)
	}

	return list
}

// ======================================================
// INSERT
// ======================================================
func DbInsert(host string, env string, foremanId int, ctx *user.GlobalCTX) {

	eId := DbID(host, env, ctx)
	if eId == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into environments(host, env, foreman_id) values(?, ?, ?)")
		if err != nil {
			logger.Warning.Printf("%q, insertToEnvironments", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(host, env, foremanId)
		if err != nil {
			logger.Warning.Printf("%q, insertToEnvironments", err)
		}
	}
}

// ======================================================
// DELETE
// ======================================================
func DbDelete(host string, env string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM environments WHERE (`host` = ? and `env`=?);")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, env)
	if err != nil {
		logger.Warning.Printf("%q, DeleteEnvironment", err)
	}
}
