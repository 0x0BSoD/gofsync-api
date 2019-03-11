package main

import (
	"log"
	"strings"
	"time"
)

// ======================================================
// CHECKS
// ======================================================
func checkHG(name string, host string) bool {
	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id from hg where name=? and host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return false
	}
	return true
}

// ======================================================
// GET
// ======================================================
//func getAllHGdb(host string) {
//	db := getDBConn()
//	defer db.Close()
//
//	stmt, err := db.Prepare("select id, dump from hg where host=?")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer stmt.Close()
//
//	var list string
//	err = stmt.QueryRow(host).Scan(&list)
//	if err != nil {
//		return ""
//	}
//}

// ======================================================
// INSERT
// ======================================================
func insertHG(name string, host string, data string, pcList []int64) int64 {

	db := getDBConn()
	defer db.Close()

	if !checkHG(name, host) {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into hg(name, host, dump, pcList, created_at, updated_at) values(?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		var strPcList []string
		for _, i := range pcList {
			if i != 0 {
				strPcList = append(strPcList, String(i))
			}
		}

		res, err := stmt.Exec(name, host, data, strings.Join(strPcList, ","), time.Now(), time.Now())
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()

		lastID, _ := res.LastInsertId()
		return lastID
	}
	return -1
}

func insertHGP(sweId int64, name string, pVal string, priority int) {

	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into hg_parameters(hg_id, name, 'value', priority) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(sweId, name, pVal, priority)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

}
