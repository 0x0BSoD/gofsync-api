package environment

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
)

// ======================================================
// CHECKS and GETS
// ======================================================
func ID(host string, env string, ctx *user.GlobalCTX) int {

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

func ForemanID(host string, env string, ctx *user.GlobalCTX) int {

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
func DbGetRepo(host string, ctx *user.GlobalCTX) string {

	var r string

	stmt, err := ctx.Config.Database.DB.Prepare("select repo from environments where host=?")
	if err != nil {
		logger.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host).Scan(&r)
	if err != nil {
		return ""
	}
	return r
}

func DbGet(host string, env string, ctx *user.GlobalCTX) Environment {

	var state string
	var repo string

	stmt, err := ctx.Config.Database.DB.Prepare("select state, repo from environments where host=? and env=?")
	if err != nil {
		logger.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, env).Scan(&state, &repo)
	if err != nil {
		return Environment{}
	}
	return Environment{
		Repo:  repo,
		State: state,
	}
}
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
func DbInsert(host, env, state string, foremanId int, codeInfo SvnInfo, ctx *user.GlobalCTX) {

	meta := "{}"
	if (SvnInfo{}) != codeInfo {
		tmp, _ := json.Marshal(codeInfo)
		meta = string(tmp)
	}

	eId := ID(host, env, ctx)
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
// Update
// ======================================================
func DbSetRepo(repo, host string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("UPDATE environments SET  `repo` = ? WHERE (`host` = ?)")
	if err != nil {
		logger.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(repo, host)
	if err != nil {
		logger.Warning.Printf("%q, updateEnvironments", err)
	}
}

func DbSetUpdated(state, host, name string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("UPDATE environments SET  `state` = ? WHERE (`host` = ?) AND (`env` = ?)")
	if err != nil {
		logger.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(state, host, name)
	if err != nil {
		logger.Warning.Printf("%q, updateEnvironments", err)
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
