package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
)

// ======================================================
// CHECKS
// ======================================================
func checkPC(subclass string, host string) int64 {

	stmt, err := globConf.DB.Prepare("select id from puppet_classes where host=? and subclass=?")
	if err != nil {
		log.Fatal(err)
	}
	var id int64
	err = stmt.QueryRow(host, subclass).Scan(&id)
	if err != nil {
		return -1
	}
	stmt.Close()
	return id
}
func checkPCHostId(host string, pcId int) int {

	q := fmt.Sprintf("select id from pc_host_ids where pc_id=%d and '%s' = -1", pcId, host)
	var id int
	err := globConf.DB.QueryRow(q).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func getByNamePC(subclass string, host string) PC {

	stmt, err := globConf.DB.Prepare("select id, class, subclass, sc_ids, env_ids, hg_ids, foreman_id from puppet_classes where subclass=? and host=?")
	if err != nil {
		log.Fatal(err)
	}

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
		}
	}

	stmt.Close()

	return r
}
func getPC(pId int) PC {

	stmt, err := globConf.DB.Prepare("select class, subclass, sc_ids, env_ids, hg_ids from puppet_classes where id=?")
	if err != nil {
		log.Fatal(err)
	}

	var class string
	var subclass string
	var sCIDs string
	var envIDs string
	var hGIDs string

	err = stmt.QueryRow(pId).Scan(&class, &subclass, &sCIDs, &envIDs, &hGIDs)

	stmt.Close()

	return PC{
		Class:    class,
		Subclass: subclass,
		SCIDs:    sCIDs,
	}
}

// PuppetclassesNI for getting from base
type PuppetclassesNI struct {
	Class     string
	SubClass  string
	ForemanID int
}

func getAllPCBase(host string) []PuppetclassesNI {

	stmt, err := globConf.DB.Prepare("select foreman_id, class, subclass from puppet_classes where host=?")
	if err != nil {
		log.Fatal(err)
	}

	var r []PuppetclassesNI

	rows, err := stmt.Query(host)
	if err != nil {
		return []PuppetclassesNI{}
	}
	for rows.Next() {
		var foremanId int
		var class string
		var subClass string
		err = rows.Scan(&foremanId, &class, &subClass)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, PuppetclassesNI{class, subClass, foremanId})
	}

	rows.Close()
	stmt.Close()

	return r
}

// ======================================================
// INSERT
// ======================================================
func insertPC(host string, class string, subclass string, foremanId int) int64 {

	existID := checkPC(subclass, host)
	if existID == -1 {
		stmt, err := globConf.DB.Prepare("insert into puppet_classes(host, class, subclass, foreman_id, sc_ids, env_ids, hg_ids) values(?,?,?,?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}

		res, err := stmt.Exec(host, class, subclass, foremanId, "NULL", "NULL", "NULL")
		if err != nil {
			log.Fatal(err)
		}
		stmt.Close()

		lastID, _ := res.LastInsertId()
		return lastID
	} else {
		return existID
	}
}

func updatePC(host string, subClass string, data PCSCParameters) {

	var strScList []string
	var strEnvList []string
	var strHGList []string

	for _, i := range data.SmartClassParameters {
		scID := checkSC(data.Name, i.Name, host)
		if scID != -1 {
			strScList = append(strScList, strconv.Itoa(int(scID)))
		}
	}

	// TODO: Will be see, maybe its not needed
	//for _, i := range data.Environments {
	//	scID := checkEnv(host, i.Name)
	//	strEnvList = append(strEnvList, strconv.Itoa(int(scID)))
	//}
	//
	//for _, i := range data.HostGroups {
	//	scID := checkHGID(i.Name, host)
	//	strHGList = append(strHGList, strconv.Itoa(int(scID)))
	//}

	stmt, err := globConf.DB.Prepare("update puppet_classes set sc_ids=?, env_ids=?, hg_ids=? where host=? and subclass=?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(
		strings.Join(strScList, ","),
		strings.Join(strEnvList, ","),
		strings.Join(strHGList, ","),
		host,
		subClass)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Close()
}
