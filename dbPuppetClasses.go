package main

import (
	"fmt"
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
func checkPCHostId(host string, pcId int) int {

	db := getDBConn()
	defer db.Close()

	q := fmt.Sprintf("select id from pc_host_ids where pc_id=%d and '%s' = -1", pcId, host)
	var id int
	err := db.QueryRow(q).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

// ======================================================
// GET
// ======================================================
func getByNamePC(subclass string, host string) PC {

	db := getDBConn()
	defer db.Close()
	stmt, err := db.Prepare("select id, class, subclass, sc_ids, env_ids, hg_ids, foreman_id from puppet_classes where subclass=? and host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var r PC

	rows, err := stmt.Query(subclass, host)
	if err != nil {
		return PC{}
	}
	for rows.Next() {
		var class string
		var subclass string
		var sCIDs string
		var envIDs string
		var hGIDs string
		var foremanId int
		var id int
		err = rows.Scan(&id, &class, &subclass, &sCIDs, &envIDs, &hGIDs, &foremanId)
		if err != nil {
			log.Fatal(err)
		}
		r = PC{
			ID:        id,
			ForemanId: foremanId,
			Class:     class,
			Subclass:  subclass,
			SCIDs:     sCIDs,
			//EnvIDs: envIDs,
			//HGIDs: hGIDs,
		}
	}
	return r
}
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
	defer rows.Close()

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
func insertPC(host string, class string, subclass string, foremanId int) int64 {

	db := getDBConn()
	defer db.Close()

	existID := checkPC(subclass, host)
	if existID == -1 {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := tx.Prepare("insert into puppet_classes(host, class, subclass, foreman_id) values(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(host, class, subclass, foremanId)
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
func insertPCHostID(host string, pcId int, id int) {
	lastId := checkPCHostId(host, pcId)
	if lastId == -1 {
		db := getDBConn()
		defer db.Close()

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		q := fmt.Sprintf("update pc_host_ids set '%s'=? where pc_id=?", host)
		stmt, err := tx.Prepare(q)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(id, pcId)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
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
