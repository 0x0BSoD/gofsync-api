package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// ======================================================
// CHECKS
// ======================================================
func checkEnv(host string, env string) int {
	db := getDBConn()
	defer db.Close()
	stmt, err := db.Prepare("select id from environments where host=? and env=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(host, env).Scan(&id)
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
func insertToEnvironments(host string, env string) {
	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	eId := checkEnv(host, env)
	if eId == -1 {

		stmt, err := tx.Prepare("insert into environments(host, env) values(?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(host, env)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
	}
}
