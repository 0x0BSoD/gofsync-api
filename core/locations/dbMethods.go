package locations

import (
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

// ======================================================
// CHECKS
// ======================================================
func CheckLoc(host string, loc string, cfg *models.Config) int {

	var id int

	stmt, err := cfg.Database.DB.Prepare("select id from goFsync.locations where host=? and loc=?")
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
func GetAllLocNames(host string, cfg *models.Config) []string {

	var res []string

	stmt, err := cfg.Database.DB.Prepare("select loc from goFsync.locations where host=?")
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

func GetAllLocations(host string, cfg *models.Config) []int {

	var foremanIds []int

	stmt, err := cfg.Database.DB.Prepare("select foreman_id from goFsync.locations where host=?")
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
func InsertToLocations(host string, loc string, foremanId int, cfg *models.Config) {

	eId := CheckLoc(host, loc, cfg)
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
func DeleteLocation(host string, loc string, cfg *models.Config) {
	stmt, err := cfg.Database.DB.Prepare("DELETE FROM `goFsync`.`locations` WHERE (`host` = ? and `loc`=?);")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(host, loc)
	if err != nil {
		logger.Warning.Printf("%q, deleteLocation", err)
	}
}
