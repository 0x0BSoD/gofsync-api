package locations

import (
	"git.ringcentral.com/archops/goFsync/core/user"
	logger "git.ringcentral.com/archops/goFsync/utils"
)

// ======================================================
// CHECKS
// ======================================================
func DbID(host, loc string, ctx *user.GlobalCTX) int {

	var id int

	stmt, err := ctx.Config.Database.DB.Prepare("select id from locations where host=? and loc=?")
	if err != nil {
		logger.Warning.Printf("%q, checkLoc", err)
	}
	logger.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, loc).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func DbAll(host string, ctx *user.GlobalCTX) ([]string, string) {

	var res []string
	var env string
	stmt, err := ctx.Config.Database.DB.Prepare("select l.loc, h.env from locations as l, hosts as h where l.host=? and l.host=h.host")
	if err != nil {
		logger.Warning.Println(err)
	}
	logger.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		logger.Warning.Printf("%q, getAllLocNames", err)
	}

	for rows.Next() {
		var loc string
		err = rows.Scan(&loc, &env)
		if err != nil {
			logger.Warning.Printf("%q, getAllLocNames", err)
		}
		res = append(res, loc)
	}
	return res, env
}

func DbAllForemanID(host string, ctx *user.GlobalCTX) []int {

	var foremanIds []int

	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id from locations where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	logger.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		logger.Warning.Printf("%q, getAllLocations", err)
	}

	for rows.Next() {
		var foremanId int
		err = rows.Scan(&foremanId)
		if err != nil {
			logger.Warning.Printf("%q, getAllLocations", err)
		}
		foremanIds = append(foremanIds, foremanId)
	}

	return foremanIds
}

// ======================================================
// INSERT
// ======================================================
func DbInsert(host, loc string, foremanId int, ctx *user.GlobalCTX) {

	eId := DbID(host, loc, ctx)
	if eId == -1 {

		stmt, err := ctx.Config.Database.DB.Prepare("insert into locations(host, loc, foreman_id) values(?, ?, ?)")
		if err != nil {
			logger.Warning.Printf("%q, insertToLocations", err)
		}
		logger.DeferCloseStmt(stmt)

		_, err = stmt.Exec(host, loc, foremanId)
		if err != nil {
			logger.Warning.Printf("%q, insertToLocations", err)
		}
	}
}

// ======================================================
// DELETE
// ======================================================
func DbDelete(host, loc string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM locations WHERE (`host` = ? and `loc`=?);")
	if err != nil {
		logger.Warning.Println(err)
	}
	logger.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, loc)
	if err != nil {
		logger.Warning.Printf("%q, deleteLocation", err)
	}
}
