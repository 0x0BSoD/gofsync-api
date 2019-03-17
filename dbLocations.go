package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// ======================================================
// CHECKS
// ======================================================
func checkLoc(host string, loc string) int {
	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id from locations where host=? and loc=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(host, loc).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================

// ======================================================
// INSERT
// ======================================================
func insertToLocations(host string, loc string) {
	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	eId := checkLoc(host, loc)
	if eId == -1 {

		stmt, err := tx.Prepare("insert into locations(host, loc) values(?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(host, loc)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
	}
}
