package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

// ======================================================
// CHECKS
// ======================================================
// Check HG by name
func checkHG(name string, host string) int {

	stmt, err := globConf.DB.Prepare("select id from hg where name=? and host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}
func checkHGID(name string, host string) int {

	stmt, err := globConf.DB.Prepare("select foreman_id from hg where name=? and host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}
func checkParams(hgId int64, name string) int {

	stmt, err := globConf.DB.Prepare("select id from hg_parameters where hg_id=? and name=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(hgId, name).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

// ======================================================
// GET
// ======================================================
func getHGAllList() []HGListElem {

	stmt, err := globConf.DB.Prepare("select id, name from hg")
	if err != nil {
		log.Fatal(err)
	}

	var list []HGListElem
	var chkList []string

	rows, err := stmt.Query()
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
		if !stringInSlice(name, chkList) {
			chkList = append(chkList, name)
			list = append(list, HGListElem{
				ID:   id,
				Name: name,
			})
		}

	}

	rows.Close()
	stmt.Close()

	return list
}

// For Web Server =======================================
func getHGList(host string) []HGListElem {

	stmt, err := globConf.DB.Prepare("select id, name from hg where host=?")
	if err != nil {
		log.Fatal(err)
	}

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

	rows.Close()
	stmt.Close()

	return list
}

func getHGParams(hgId int) []HGParam {

	stmt, err := globConf.DB.Prepare("select name, value from hg_parameters where hg_id=?")
	if err != nil {
		log.Fatal(err)
	}

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

	rows.Close()
	stmt.Close()

	return list
}

func getHG(id int) HGElem {

	// VARS
	var d SWE
	var name string
	var pClassesStr string
	var dump string
	var foremanId int
	pClasses := make(map[string][]PuppetClassesWeb)

	// Hg Data
	stmt, err := globConf.DB.Prepare("select foreman_id, name, pcList, dump from hg where id=?")
	if err != nil {
		log.Println("HostGroup getting..", err)
	}

	err = stmt.QueryRow(id).Scan(&foremanId, &name, &pClassesStr, &dump)
	if err != nil {
		log.Println("HostGroup getting..", err)
	}

	// HG Parameters
	params := getHGParams(id)

	err = json.Unmarshal([]byte(dump), &d)
	if err != nil {
		log.Fatalf("Error on Parsing HG: %s", err)
	}

	// PuppetClasses and Parameters
	for _, cl := range Integers(pClassesStr) {
		res := getPC(cl)

		var SCList []string
		var OvrList []SCOParams
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
			Subclass:     res.Subclass,
			SmartClasses: SCList,
			Overrides:    OvrList,
		})
	}

	stmt.Close()

	return HGElem{
		ID:            id,
		ForemanID:     foremanId,
		Name:          name,
		Params:        params,
		Environment:   d.EnvironmentName,
		ParentId:      d.Ancestry,
		PuppetClasses: pClasses,
	}

}

// ======================================================
// INSERT
// ======================================================
func insertHG(name string, host string, data string, foremanId int) int64 {
	hgExist := checkHG(name, host)
	if hgExist == -1 {

		stmt, err := globConf.DB.Prepare("insert into hg(name, host, dump, created_at, updated_at, foreman_id, pcList, locList) values(?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err := stmt.Exec(name, host, data, time.Now(), time.Now(), foremanId, "NULL", "NULL")
		if err != nil {
			log.Fatal(err)
		}

		stmt.Close()

		lastID, _ := res.LastInsertId()
		return lastID
	}
	return int64(hgExist)
}

func insertHGP(sweId int64, name string, pVal string, priority int) {

	oldId := checkParams(sweId, host)
	if oldId == -1 {
		stmt, err := globConf.DB.Prepare("insert into hg_parameters(hg_id, name, `value`, priority) values(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(sweId, name, pVal, priority)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func updateLocInHG(hgId int64, lIdList []int64) {

	db := *globConf.DB

	var strLocList []string

	for _, i := range lIdList {
		if i != 0 {
			strLocList = append(strLocList, String(i))
		}
	}
	LocList := strings.Join(strLocList, ",")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("update hg set locList=? where id=?")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(LocList, hgId)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}
func updatePCinHG(hgId int64, pcList []int64) {

	var strPcList []string

	for _, i := range pcList {
		if i != 0 {
			strPcList = append(strPcList, String(i))
		}
	}
	pcListStr := strings.Join(strPcList, ",")
	stmt, err := globConf.DB.Prepare("update hg set pcList=? where id=?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(pcListStr, hgId)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Close()
}

// ======================================================
// DELETE
// ======================================================
func deleteHGbyId(hgId int) {
	stmt, err := globConf.DB.Prepare("DELETE FROM hg WHERE id=?")
	if err != nil {
		log.Println(err)
	}

	_, err = stmt.Exec(hgId)
	if err != nil {
		log.Println(err)
	}

	stmt.Close()
}
