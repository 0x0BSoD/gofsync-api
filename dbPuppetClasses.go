package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// ======================================================
// CHECKS
// ======================================================
func checkPC(subclass string, host string) int64 {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id from puppet_classes where host=? and subclass=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(host, subclass).Scan(&id)
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
func insertPC(host string, class string, subclass string) int64 {

	db := getDBConn()
	defer db.Close()
	existID := checkPC(subclass, host)
	if existID == -1 {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into puppet_classes(host, class, subclass) values(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(host, class, subclass)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()

		lastID, _ := res.LastInsertId()
		return lastID
	} else {
		return existID
	}
}
