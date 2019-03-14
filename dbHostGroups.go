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
func checkHGID(name string, host string) int {
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
		return id
	}
	return id
}

// ======================================================
// GET
// ======================================================
func getHGDump(host string) []string {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select dump from hg where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var HGs []string

	rows, err := stmt.Query(host)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var dump string
		err = rows.Scan(&dump)
		if err != nil {
			log.Fatal(err)
		}
		HGs = append(HGs, dump)
	}
	return HGs
}

// For Web Server =======================================
func getHGList(host string) []HGListElem {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id, name from hg where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var list []HGListElem

	rows, err := stmt.Query(host)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, HGListElem{
			ID:   id,
			Name: name,
		})
	}
	return list
}

func getHGParams(hgId int) []HGParam {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select name, value from hg_parameters where hg_id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var list []HGParam

	rows, err := stmt.Query(hgId)
	if err != nil {
		return []HGParam{}
	}
	for rows.Next() {
		var name string
		var value string
		err = rows.Scan(&name, &value)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, HGParam{
			Name:  name,
			Value: value,
		})
	}
	return list
}

func getHG(host string, id string) []HGElem {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id, name, pcList from hg where host=? and id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var list []HGElem
	pClasses := make(map[string][]PuppetClassesWeb)

	rows, err := stmt.Query(host, id)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int
		var name string
		var pClassesStr string
		err = rows.Scan(&id, &name, &pClassesStr)
		if err != nil {
			log.Fatal(err)
		}
		params := getHGParams(id)

		for _, cl := range Integers(pClassesStr) {
			res := getPC(cl)
			var SCList []string
			var OvrList []SCOParams
			//envList := Integers(res.EnvIDs)
			//hgList := Integers(res.HGIDs)

			scList := Integers(res.SCIDs)
			for _, SCID := range scList {
				data := getSCData(SCID)
				if data.Name != "" {
					SCList = append(SCList, data.Name)
				}
				if data.OverrideValuesCount > 0 {
					ovrData := getOvrData(SCID, name, data.Name)
					for _, p := range ovrData {
						OvrList = append(OvrList, p)
					}
				}
			}

			pClasses[res.Class] = append(pClasses[res.Class], PuppetClassesWeb{
				Subclass: res.Subclass,
				//EnvIds:envList,
				//HostGroupsIds:hgList,
				SmartClasses: SCList,
				Overrides:    OvrList,
			})
		}

		list = append(list, HGElem{
			ID:            id,
			Name:          name,
			Params:        params,
			PuppetClasses: pClasses,
		})
	}
	return list
}

// ======================================================
// INSERT
// ======================================================
func insertHG(name string, host string, data string) int64 {

	db := getDBConn()
	defer db.Close()

	if !checkHG(name, host) {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into hg(name, host, dump, created_at, updated_at) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		res, err := stmt.Exec(name, host, data, time.Now(), time.Now())
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()

		lastID, _ := res.LastInsertId()
		return lastID
	}
	return -1
}

func updatePCinHG(hgId int64, pcList []int64) {

	var strPcList []string
	db := getDBConn()
	defer db.Close()

	for _, i := range pcList {
		if i != 0 {
			strPcList = append(strPcList, String(i))
		}
	}
	pcListStr := strings.Join(strPcList, ",")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("update hg set pcList=? where id=?")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(pcListStr, hgId)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
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
