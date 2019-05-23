package locations

import (
	"git.ringcentral.com/archops/goFsync/models"
	logger "git.ringcentral.com/archops/goFsync/utils"
)

// ======================================================
// CHECKS
// ======================================================
func DbID(host string, loc string, cfg *models.Config) int {

	var id int

	stmt, err := cfg.Database.DB.Prepare("select id from locations where host=? and loc=?")
	if err != nil {
		logger.Warning.Printf("%q, checkLoc", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(host, loc).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func DbAll(host string, cfg *models.Config) []string {

	var res []string

	stmt, err := cfg.Database.DB.Prepare("select loc from locations where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(host)
	if err != nil {
		logger.Warning.Printf("%q, getAllLocNames", err)
	}

	for rows.Next() {
		var loc string
		err = rows.Scan(&loc)
		if err != nil {
			logger.Warning.Printf("%q, getAllLocNames", err)
		}
		res = append(res, loc)
	}
	return res
}

func DbAllForemanID(host string, cfg *models.Config) []int {

	var foremanIds []int

	stmt, err := cfg.Database.DB.Prepare("select foreman_id from locations where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

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
func DbInsert(host string, loc string, foremanId int, cfg *models.Config) {

	eId := DbID(host, loc, cfg)
	if eId == -1 {

		stmt, err := cfg.Database.DB.Prepare("insert into locations(host, loc, foreman_id) values(?, ?, ?)")
		if err != nil {
			logger.Warning.Printf("%q, insertToLocations", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(host, loc, foremanId)
		if err != nil {
			logger.Warning.Printf("%q, insertToLocations", err)
		}
	}
}

// ======================================================
// DELETE
// ======================================================
func DbDelete(host string, loc string, cfg *models.Config) {
	stmt, err := cfg.Database.DB.Prepare("DELETE FROM locations WHERE (`host` = ? and `loc`=?);")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(host, loc)
	if err != nil {
		logger.Warning.Printf("%q, deleteLocation", err)
	}
}
