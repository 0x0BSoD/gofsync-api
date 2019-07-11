package environment

import (
	"encoding/json"
	"fmt"
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
func DbAll(ctx *user.GlobalCTX) map[string][]Environment {

	list := make(map[string][]Environment)

	rows, err := ctx.Config.Database.DB.Query("select id, host, env, state, repo from environments")
	if err != nil {
		logger.Warning.Printf("%q, getEnvList", err)
	}

	for rows.Next() {
		var ID int
		var env string
		var host string
		var state string
		var repo string
		err = rows.Scan(&ID, &host, &env, &state, &repo)
		if err != nil {
			logger.Error.Printf("%q, getEnvList", err)
		}
		list[host] = append(list[host], Environment{
			ID:    ID,
			Name:  env,
			State: state,
			Repo:  repo,
		})
	}

	err = rows.Close()
	if err != nil {
		logger.Error.Printf("%q, getEnvList", err)
	}

	return list
}
func DbByHost(host string, ctx *user.GlobalCTX) []string {

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
func DbInsert(host string, env string, foremanId int, codeInfo utils.SvnInfo, ctx *user.GlobalCTX) {

	fmt.Println(codeInfo)

	meta := "{}"
	state := "absent"

	if (utils.SvnInfo{}) != codeInfo {
		tmp, _ := json.Marshal(codeInfo)
		meta = string(tmp)
		if codeInfo.LastRev == codeInfo.Revision {
			state = "ok"
		} else {
			state = "outdated"
		}
	}

	eId := DbID(host, env, ctx)
	if eId == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into environments(host, env, meta, state, foreman_id) values(?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Printf("%q, insertToEnvironments", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(host, env, meta, state, foremanId)
		if err != nil {
			logger.Warning.Printf("%q, insertToEnvironments", err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE environments SET  `meta` = ?, `state` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(meta, state, eId)
		if err != nil {
			logger.Warning.Printf("%q, updateEnvironments", err)
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
