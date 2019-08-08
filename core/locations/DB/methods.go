package DB

import (
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

func (Get) ID(host, loc string, ctx *user.GlobalCTX) int {

	// VARS
	var ID int

	// ======
	stmt, err := ctx.Config.Database.DB.Prepare("select id from locations where host=? and loc=?")
	if err != nil {
		utils.Warning.Printf("%q, checkLoc", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, loc).Scan(&ID)
	if err != nil {
		ID = -1
	}

	return ID
}

func (Get) ForemanIDs(host string, ctx *user.GlobalCTX) []int {

	// VARS
	var foremanIds []int

	// =====
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id from locations where host=?")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		utils.Warning.Printf("%q, getAllLocations", err)
	}

	for rows.Next() {
		var foremanId int
		err = rows.Scan(&foremanId)
		if err != nil {
			utils.Warning.Printf("%q, getAllLocations", err)
		}
		foremanIds = append(foremanIds, foremanId)
	}

	return foremanIds
}

func (Get) All(host string, ctx *user.GlobalCTX) ([]string, string) {

	// VARS
	var res []string
	var env string

	stmt, err := ctx.Config.Database.DB.Prepare("select l.loc, h.env from locations as l, hosts as h where l.host=? and l.host=h.host")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		utils.Warning.Printf("%q, getAllLocNames", err)
	}

	for rows.Next() {
		var loc string
		err = rows.Scan(&loc, &env)
		if err != nil {
			utils.Warning.Printf("%q, getAllLocNames", err)
		}
		res = append(res, loc)
	}

	return res, env
}

// =====================================================================================================================
// INSERT
// =====================================================================================================================

func (Insert) Add(host, loc string, foremanId int, ctx *user.GlobalCTX) {

	// VARS
	var gDB Get
	ID := gDB.ID(host, loc, ctx)

	// =====
	if ID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into locations(host, loc, foreman_id) values(?, ?, ?)")
		if err != nil {
			utils.Warning.Printf("%q, insertToLocations", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(host, loc, foremanId)
		if err != nil {
			utils.Warning.Printf("%q, insertToLocations", err)
		}
	}
}

// =====================================================================================================================
// DELETE
// =====================================================================================================================

func (Delete) ByName(host, loc string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM locations WHERE (`host` = ? and `loc`=?);")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, loc)
	if err != nil {
		utils.Warning.Printf("%q, deleteLocation", err)
	}
}
