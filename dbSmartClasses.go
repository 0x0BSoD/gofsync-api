package main

import (
	"encoding/json"
	"log"
)

// ======================================================
// CHECKS
// ======================================================
func checkSC(parameter string, host string) int64 {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id from smart_classes where host=? and parameter=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(host, parameter).Scan(&id)
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
func insertSC(host string, data SCParameter) {
	db := getDBConn()
	defer db.Close()

	existID := checkSC(data.Parameter, host)

	if existID == -1 {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into smart_classes(host, parameter, override_values_count, dump) values(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		sJson, _ := json.Marshal(data)

		_, err = stmt.Exec(host, data.Parameter, data.OverrideValuesCount, sJson)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
	}
}
