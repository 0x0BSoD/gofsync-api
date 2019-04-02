package main

import (
	"git.ringcentral.com/alexander.simonov/goFsync/logger"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// ======================================================
// CHECKS
// ======================================================
func checkLoc(host string, loc string) int {

	var id int

	stmt, err := globConf.DB.Prepare("select id from goFsync.locations where host=? and loc=?")
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
func getAllLocations(host string) []int {

	stmt, err := globConf.DB.Prepare("select foreman_id from goFsync.locations where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(host)
	if err != nil {
		logger.Warning.Printf("%q, getAllLocations", err)
	}
	var foremanIds []int
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
func insertToLocations(host string, loc string, foremanId int) {

	eId := checkLoc(host, loc)
	if eId == -1 {

		stmt, err := globConf.DB.Prepare("insert into locations(host, loc, foreman_id) values(?, ?, ?)")
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
