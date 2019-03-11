package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

// ======================================================
// CHECKS
// ======================================================
func checkLoc(host string) string {
	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select list from locations where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var list string
	err = stmt.QueryRow(host).Scan(&list)
	if err != nil {
		return ""
	}
	return list
}

// ======================================================
// GET
// ======================================================

// ======================================================
// INSERT
// ======================================================
func insertToLocations(host string, list string) bool {
	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	oldList := checkLoc(host)
	if len(oldList) == 0 {

		stmt, err := tx.Prepare("insert into locations(host, list, created_at, updated_at) values(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(host, list, time.Now(), time.Now())
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
	} else {
		if oldList != list {
			stmt, err := tx.Prepare("update locations set list=?, updated_at=? where host=?")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(list, time.Now(), host)
			if err != nil {
				log.Fatal(err)
			}

			tx.Commit()
		}
	}

	return true
}
