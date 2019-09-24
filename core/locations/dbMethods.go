package locations

import (
	"database/sql"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

func ErrRow(rows *sql.Rows) {
	if err := rows.Err(); err != nil {
		utils.Error.Fatal(err)
	}
}
func ErrQuery(err error) bool {
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Trace.Println("empty result")
		} else {
			utils.Error.Fatal(err)
		}
		return true
	}
	return false
}

// ======================================================
// CHECKS
// ======================================================
func ID(host, loc string, ctx *user.GlobalCTX) int {

	stmt, err := ctx.Config.Database.DB.Prepare("select id from locations where host=? and loc=?")
	if err != nil {
		utils.Warning.Printf("%q, checkLoc", err)
	}
	defer utils.DeferCloseStmt(stmt)

	// VARS
	var id int
	err = stmt.QueryRow(host, loc).Scan(&id)
	if ErrQuery(err) {
		return -1
	}

	return id
}

// ======================================================
// GET
// ======================================================
func DbAll(host string, ctx *user.GlobalCTX) ([]string, string) {

	stmt, err := ctx.Config.Database.DB.Prepare("select l.loc, h.env from locations as l, hosts as h where l.host=? and l.host=h.host")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	ErrQuery(err)
	cols, _ := rows.Columns()

	// VARS
	res := make([]string, 0, len(cols))
	var env string
	var loc string

	for rows.Next() {
		err = rows.Scan(&loc, &env)
		if err != nil {
			utils.Warning.Printf("%q, getAllLocNames", err)
		}
		res = append(res, loc)
	}
	ErrRow(rows)

	return res, env
}

func DbAllForemanID(host string, ctx *user.GlobalCTX) []int {

	var foremanIds []int

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

// ======================================================
// INSERT
// ======================================================
func DbInsert(host, loc string, foremanId int, ctx *user.GlobalCTX) {

	ID := ID(host, loc, ctx)
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
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE locations SET  `foreman_id` =? WHERE (`id` = ?)")
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(foremanId, ID)
		if err != nil {
			utils.Warning.Printf("%q, updateEnvironments", err)
		}
	}
}

// ======================================================
// DELETE
// ======================================================
func DbDelete(host, loc string, ctx *user.GlobalCTX) {
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
