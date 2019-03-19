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
func getAllLocations(host string) []int {
	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select foreman_id from locations where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

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

	return foremanIds
}

// ======================================================
// INSERT
// ======================================================
func insertToLocations(host string, loc string, foremanId int) {
	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	eId := checkLoc(host, loc)
	if eId == -1 {

		stmt, err := tx.Prepare("insert into locations(host, loc, foreman_id) values(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(host, loc, foremanId)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
	}
}
