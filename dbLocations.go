package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// ======================================================
// CHECKS
// ======================================================
func checkLoc(host string, loc string) int {

	stmt, err := globConf.DB.Prepare("select id from goFsync.locations where host=? and loc=?")
	if err != nil {
		log.Fatalf("SQL: %q, \n checkLoc", err)
	}

	var id int
	err = stmt.QueryRow(host, loc).Scan(&id)
	if err != nil {
		stmt.Close()
		return -1
	}
	stmt.Close()
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

	rows, err := stmt.Query(host)
	if err != nil {
		log.Fatal(err)
	}
	var foremanIds []int
	for rows.Next() {
		//var location string
		var foremanId int
		err = rows.Scan(&foremanId)
		if err != nil {
			log.Fatal(err)
		}
		foremanIds = append(foremanIds, foremanId)
	}

	stmt.Close()

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
			log.Fatalf("SQL: %q, \n insertToLocations", err)
		}

		_, err = stmt.Exec(host, loc, foremanId)
		if err != nil {
			log.Fatalf("SQL: %q, \n insertToLocations", err)
		}
		stmt.Close()
	}
}
