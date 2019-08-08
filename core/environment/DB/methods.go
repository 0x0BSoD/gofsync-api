package DB

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/environment/API"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Return environment ID by host and env name
func (Get) ID(host, env string, ctx *user.GlobalCTX) int {

	// VARS
	var ID int

	// =============
	stmt, err := ctx.Config.Database.DB.Prepare("select id from environments where host=? and env=?")
	if err != nil {
		utils.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, env).Scan(&ID)
	if err != nil {
		ID = -1
	}

	return ID
}

// Return environment Foreman ID by host and env name
func (Get) ForemanID(host, env string, ctx *user.GlobalCTX) int {

	// VARS
	var foremanID int

	// =============
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id from environments where host=? and env=?")
	if err != nil {
		utils.Warning.Printf("%q, checkPostEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, env).Scan(&foremanID)
	if err != nil {
		foremanID = -1
	}

	return foremanID
}

// Return environment by name
func (Get) ByName(host string, env string, ctx *user.GlobalCTX) API.Environment {

	// VARS
	var state string
	var repo string

	// ============
	stmt, err := ctx.Config.Database.DB.Prepare("select state, repo from environments where host=? and env=?")
	if err != nil {
		utils.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, env).Scan(&state, &repo)
	if err != nil {
		return API.Environment{}
	}

	return API.Environment{
		Repo:  repo,
		State: state,
	}
}

// Return list of environments by host
func (Get) ByHost(host string, ctx *user.GlobalCTX) []string {

	// VARS
	var list []string

	// ============
	stmt, err := ctx.Config.Database.DB.Prepare("select env from environments where host=?")
	if err != nil {
		utils.Warning.Printf("%q, getEnvList", err)
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
			utils.Error.Printf("%q, getEnvList", err)
		}
		list = append(list, env)
	}

	return list
}

// Return all environments with hosts
func (Get) All(ctx *user.GlobalCTX) map[string][]API.Environment {

	list := make(map[string][]API.Environment)

	rows, err := ctx.Config.Database.DB.Query("select id, host, env, state, repo from environments")
	if err != nil {
		utils.Warning.Printf("%q, getEnvList", err)
	}

	for rows.Next() {

		var (
			ID    int
			env   string
			host  string
			state string
			repo  string
		)

		err = rows.Scan(&ID, &host, &env, &state, &repo)
		if err != nil {
			utils.Error.Printf("%q, getEnvList", err)
		}
		list[host] = append(list[host], API.Environment{
			ForemanID: ID,
			Name:      env,
			State:     state,
			Repo:      repo,
		})
	}

	err = rows.Close()
	if err != nil {
		utils.Error.Printf("%q, getEnvList", err)
	}

	return list
}

// Return repo for environment
func (Get) Repo(host string, ctx *user.GlobalCTX) string {

	var repo string

	stmt, err := ctx.Config.Database.DB.Prepare("select repo from environments where host=?")
	if err != nil {
		utils.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host).Scan(&repo)
	if err != nil {
		return ""
	}

	return repo
}

// =====================================================================================================================
// INSERT
// =====================================================================================================================

// Add or updated new env
func (Insert) Add(host, env, state string, foremanId int, codeInfo SvnInfo, ctx *user.GlobalCTX) {

	// VARS
	var gDB Get
	meta := "{}"
	if (SvnInfo{}) != codeInfo {
		tmp, _ := json.Marshal(codeInfo)
		meta = string(tmp)
	}

	// ===========
	ID := gDB.ID(host, env, ctx)
	if ID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into environments(host, env, meta, state, foreman_id) values(?, ?, ?, ?, ?)")
		if err != nil {
			utils.Warning.Printf("%q, insertToEnvironments", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(host, env, meta, state, foremanId)
		if err != nil {
			utils.Warning.Printf("%q, insertToEnvironments", err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE environments SET  `meta` = ?, `state` = ? WHERE (`id` = ?)")
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(meta, state, ID)
		if err != nil {
			utils.Warning.Printf("%q, updateEnvironments", err)
		}
	}
}

// =====================================================================================================================
// UPDATE
// =====================================================================================================================

// Set repo url for env
func (Update) SetRepo(repo, host string, ctx *user.GlobalCTX) {

	// ======
	stmt, err := ctx.Config.Database.DB.Prepare("UPDATE environments SET  `repo` = ? WHERE (`host` = ?)")
	if err != nil {
		utils.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(repo, host)
	if err != nil {
		utils.Warning.Printf("%q, updateEnvironments", err)
	}
}

// Set code state for env
func (Update) SetState(state, host, name string, ctx *user.GlobalCTX) {

	// ======
	stmt, err := ctx.Config.Database.DB.Prepare("UPDATE environments SET  `state` = ? WHERE (`host` = ?) AND (`env` = ?)")
	if err != nil {
		utils.Warning.Printf("%q, checkEnv", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(state, host, name)
	if err != nil {
		utils.Warning.Printf("%q, updateEnvironments", err)
	}
}

// =====================================================================================================================
// DELETE
// =====================================================================================================================

// Remove env by name
func (Delete) ByName(host string, env string, ctx *user.GlobalCTX) {

	// ======
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM environments WHERE (`host` = ? and `env`=?);")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, env)
	if err != nil {
		utils.Warning.Printf("%q, DeleteEnvironment", err)
	}
}
