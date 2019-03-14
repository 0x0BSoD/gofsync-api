package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"strings"
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
func getPC(pId int) PC {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select class, subclass, sc_ids, env_ids, hg_ids from puppet_classes where id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var r PC

	rows, err := stmt.Query(pId)
	if err != nil {
		return PC{}
	}
	for rows.Next() {
		var class string
		var subclass string
		var sCIDs string
		var envIDs string
		var hGIDs string
		err = rows.Scan(&class, &subclass, &sCIDs, &envIDs, &hGIDs)
		if err != nil {
			log.Fatal(err)
		}
		r = PC{
			Class:    class,
			Subclass: subclass,
			SCIDs:    sCIDs,
			//EnvIDs: envIDs,
			//HGIDs: hGIDs,
		}
	}
	return r
}
func getAllPCBase(host string) []string {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select subclass from puppet_classes where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var r []string

	rows, err := stmt.Query(host)
	if err != nil {
		return []string{}
	}
	for rows.Next() {
		var subclass string
		err = rows.Scan(&subclass)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, subclass)
	}
	return r
}

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

func updatePC(host string, ss string, data PCSCParameters) {

	var strScList []string
	var strEnvList []string
	var strHGList []string

	db := getDBConn()
	defer db.Close()

	for _, i := range data.SmartClassParameters {
		scID := checkSC(i.Name, host)
		strScList = append(strScList, strconv.Itoa(int(scID)))
	}

	for _, i := range data.Environments {
		scID := checkEnv(host, i.Name)
		strEnvList = append(strEnvList, strconv.Itoa(int(scID)))
	}

	for _, i := range data.HostGroups {
		scID := checkHGID(i.Name, host)
		strHGList = append(strHGList, strconv.Itoa(int(scID)))
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("update puppet_classes set sc_ids=?, env_ids=?, hg_ids=? where host=? and subclass=?")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		strings.Join(strScList, ","),
		strings.Join(strEnvList, ","),
		strings.Join(strHGList, ","),
		host,
		ss)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}
