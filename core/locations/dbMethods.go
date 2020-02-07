package locations

import (
	"database/sql"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// ======================================================
// STATEMENTS
// ======================================================
var (
	selectID            = "select id from locations where host_id=? and name=?"
	selectForemanID     = "select foreman_id from locations where host_id=? and name=?"
	selectAllForemanIDs = "select foreman_id from locations where host_id=?"
	selectAll           = "select l.name, h.env from locations as l inner join hosts as h on l.host_id=h.id where host_id=?" // TODO: host.name required?

	insert = "insert into locations(host_id, name, foreman_id) values(?, ?, ?)"

	updateForemanID = "update locations set  `foreman_id`=? where (`id`=?)"

	deleteLoc = "DELETE FROM locations WHERE (`host_id` = ? and `name`=?);"
)

func ErrRow(rows *sql.Rows) {
	if err := rows.Err(); err != nil {
		utils.Error.Fatal(err)
	}
}
func ErrQuery(err error) bool {
	if err != nil {
		if err == sql.ErrNoRows {
			//utils.Trace.Trace("empty result")
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
func ID(hostID int, name string, ctx *user.GlobalCTX) int {
	stmt, err := ctx.Config.Database.DB.Prepare(selectID)
	if err != nil {
		utils.Warning.Printf("%q, checkLoc", err)
	}
	defer utils.DeferCloseStmt(stmt)

	// VARS
	var id int
	err = stmt.QueryRow(hostID, name).Scan(&id)
	if ErrQuery(err) {
		return -1
	}

	return id
}

func ForemanID(hostID int, name string, ctx *user.GlobalCTX) int {
	stmt, err := ctx.Config.Database.DB.Prepare(selectForemanID)
	if err != nil {
		utils.Warning.Printf("%q, checkLoc", err)
	}
	defer utils.DeferCloseStmt(stmt)

	// VARS
	var id int
	err = stmt.QueryRow(hostID, name).Scan(&id)
	if ErrQuery(err) {
		return -1
	}

	return id
}

// ======================================================
// GET
// ======================================================
func DbAll(hostID int, ctx *user.GlobalCTX) ([]string, string) {

	stmt, err := ctx.Config.Database.DB.Prepare(selectAll)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(hostID)
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

func DbAllForemanID(hostID int, ctx *user.GlobalCTX) []int {

	var foremanIds []int

	stmt, err := ctx.Config.Database.DB.Prepare(selectAllForemanIDs)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(hostID)
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
func DbInsert(hostID int, foremanID int, name string, ctx *user.GlobalCTX) {

	ID := ID(hostID, name, ctx)
	if ID == -1 {

		stmt, err := ctx.Config.Database.DB.Prepare(insert)
		if err != nil {
			utils.Warning.Printf("%q, insertToLocations", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(hostID, name, foremanID)
		if err != nil {
			utils.Warning.Printf("%q, insertToLocations", err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare(updateForemanID)
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(foremanID, ID)
		if err != nil {
			utils.Warning.Printf("%q, updateEnvironments", err)
		}
	}
}

// ======================================================
// DELETE
// ======================================================
func DbDelete(hostID int, name string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare(deleteLoc)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(hostID, name)
	if err != nil {
		utils.Warning.Printf("%q, deleteLocation", err)
	}
}
