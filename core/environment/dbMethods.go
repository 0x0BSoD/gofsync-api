package environment

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// ======================================================
// CHECKS and GETS
// ======================================================
func ID(hostID int, name string, ctx *user.GlobalCTX) int {
	var id int

	err := ctx.Config.DBGetOne("select id from environments where host_id=? and name=?", &id, hostID, name)
	if err != nil {
		return -1
	}

	return id
}

func ForemanID(hostID int, name string, ctx *user.GlobalCTX) int {
	var id int

	err := ctx.Config.DBGetOne("select foreman_id  from environments where host_id=? and name=?", &id, hostID, name)
	if err != nil {
		return -1
	}

	return id
}

func ForemanIDs(hostID int, ctx *user.GlobalCTX) []int {
	var result []int

	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id  from environments where host_id=?")
	if err != nil {
		utils.Warning.Printf("%q, GetForemanIDs", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(hostID)
	if err != nil {
		utils.Warning.Printf("%q, GetForemanIDs", err)
	}

	for rows.Next() {
		var _id int
		err = rows.Scan(&_id)
		if err != nil {
			utils.Warning.Printf("%q, GetForemanIDs", err)
		}

		result = append(result, _id)
	}

	return result
}

// ======================================================
// GET
// ======================================================
func DbGetRepo(hostID int, ctx *user.GlobalCTX) string {
	var r string

	err := ctx.Config.DBGetOne("select repo from environments where host_id=?", &r, hostID)
	if err != nil {
		return ""
	}

	return r
}

func DbGet(hostID int, env string, ctx *user.GlobalCTX) Environment {
	var state string
	var repo string

	stmt, err := ctx.Config.Database.DB.Prepare("select state, repo from environments where host_id=? and name=?")
	if err != nil {
		utils.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(hostID, env).Scan(&state, &repo)
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

	rows, err := ctx.Config.Database.DB.Query("select e.id, h.name as host, e.name, state, repo, h.env from environments as e inner join hosts h on e.host_id = h.id")
	if err != nil {
		utils.Warning.Printf("%q, getEnvList", err)
	}

	for rows.Next() {
		var ID int
		var env string
		var host string
		var state string
		var repo string
		var hostEnv string
		err = rows.Scan(&ID, &host, &env, &state, &repo, &hostEnv)
		if err != nil {
			utils.Error.Printf("%q, getEnvList", err)
		}
		if _, ok := ctx.Config.Hosts[host]; ok {
			list[host] = append(list[host], Environment{
				ID:    ID,
				Name:  env,
				State: state,
				Repo:  repo,
				Env:   hostEnv,
			})
		}

	}

	err = rows.Close()
	if err != nil {
		utils.Error.Printf("%q, getEnvList", err)
	}

	return list
}

func DbGetByHost(hostID int, ctx *user.GlobalCTX) []string {
	var list []string

	stmt, err := ctx.Config.Database.DB.Prepare("select name from environments where host_id=?")
	if err != nil {
		utils.Warning.Printf("%q, getEnvList", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(hostID)
	if err != nil {
		return list
	}

	for rows.Next() {
		var env string
		err = rows.Scan(&env)
		if err != nil {
			utils.Error.Printf("%q, getEnvList", err)
		}
		list = append(list, env)
	}

	return list
}

// ======================================================
// INSERT
// ======================================================
func DbInsert(hostID int, env, repo, state string, foremanId int, codeInfo SvnDirInfo, ctx *user.GlobalCTX) {
	meta := "{}"
	if (SvnDirInfo{}) != codeInfo {
		tmp, _ := json.Marshal(codeInfo)
		meta = string(tmp)
	}
	eId := ID(hostID, env, ctx)

	if eId == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into environments(host_id, name, repo, meta, state, foreman_id) values(?, ?, ?, ?, ?, ?)")
		if err != nil {
			utils.Warning.Printf("%q, insertToEnvironments", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(hostID, env, repo, meta, state, foremanId)
		if err != nil {
			utils.Warning.Printf("%q, insertToEnvironments", err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("update environments set  `foreman_id` =?, `meta` = ?, `state` = ? where (`id` = ?)")
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(foremanId, meta, state, eId)
		if err != nil {
			utils.Warning.Printf("%q, updateEnvironments", err)
		}
	}
}

// ======================================================
// Update
// ======================================================
func DbSetRepo(hostID int, repo string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("update environments set  `repo` = ? where (`host_id` = ?)")
	if err != nil {
		utils.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(repo, hostID)
	if err != nil {
		utils.Warning.Printf("%q, updateEnvironments", err)
	}
}

func DbSetUpdated(hostID int, name, state string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("update environments set  `state` = ? where (`host_id` = ?) and (`name` = ?)")
	if err != nil {
		utils.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(state, hostID, name)
	if err != nil {
		utils.Warning.Printf("%q, updateEnvironments", err)
	}
}

// ======================================================
// DELETE
// ======================================================
func DbDelete(hostID int, env string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("delete FROM environments where (`host_id` = ? and `name`=?);")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	res, err := stmt.Exec(hostID, env)
	if err != nil {
		utils.Warning.Printf("%q, DeleteEnvironment", err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Println(affect)
}
