package main

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"log"
	"strings"
	"time"
)

// ======================================================
// CHECKS
// ======================================================
// Check HG by name
func checkHG(name string, host string, cfg *models.Config) int {

	stmt, err := cfg.Database.DB.Prepare("select id from hg where name=? and host=?")
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
func checkHGID(name string, host string, cfg *models.Config) int {

	stmt, err := cfg.Database.DB.Prepare("select foreman_id from hg where name=? and host=?")
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
func checkParams(hgId int64, name string, cfg *models.Config) int {

	stmt, err := cfg.Database.DB.Prepare("select id from hg_parameters where hg_id=? and name=?")
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
func getHGAllList(cfg *models.Config) []models.HGListElem {

	stmt, err := cfg.Database.DB.Prepare("select id, name from hg")
	if err != nil {
		log.Fatal(err)
	}

	var list []models.HGListElem
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
		if !utils.StringInSlice(name, chkList) {
			chkList = append(chkList, name)
			list = append(list, models.HGListElem{
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
func getHGList(host string, cfg *models.Config) []models.HGListElem {

	stmt, err := cfg.Database.DB.Prepare("select id, name from hg where host=?")
	if err != nil {
		log.Fatal(err)
	}

	var list []models.HGListElem

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
		list = append(list, models.HGListElem{
			ID:   id,
			Name: name,
		})
	}

	rows.Close()
	stmt.Close()

	return list
}

func getHGParams(hgId int, cfg *models.Config) []models.HGParam {

	stmt, err := cfg.Database.DB.Prepare("select name, value from hg_parameters where hg_id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var list []models.HGParam

	rows, err := stmt.Query(hgId)
	if err != nil {
		return []models.HGParam{}
	}

	for rows.Next() {
		var name string
		var value string
		err = rows.Scan(&name, &value)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, models.HGParam{
			Name:  name,
			Value: value,
		})
	}

	return list
}

func getHG(id int, cfg *models.Config) models.HGElem {

	// VARS
	var d models.HostGroup
	var name string
	var pClassesStr string
	var dump string
	var foremanId int
	pClasses := make(map[string][]models.PuppetClassesWeb)

	// Hg Data
	stmt, err := cfg.Database.DB.Prepare("select foreman_id, name, pcList, dump from hg where id=?")
	if err != nil {
		logger.Error.Println("HostGroup getting..", err)
	}

	err = stmt.QueryRow(id).Scan(&foremanId, &name, &pClassesStr, &dump)
	if err != nil {
		logger.Error.Println("HostGroup getting..", err)
	}

	// HG Parameters
	params := getHGParams(id, cfg)

	err = json.Unmarshal([]byte(dump), &d)
	if err != nil {
		log.Fatalf("Error on Parsing HG: %s", err)
	}

	// PuppetClasses and Parameters
	for _, cl := range utils.Integers(pClassesStr) {
		res := getPC(cl, cfg)

		var SCList []string
		var OvrList []models.SCOParams
		scList := utils.Integers(res.SCIDs)
		for _, SCID := range scList {
			data := getSCData(SCID, cfg)
			if data.Name != "" {
				SCList = append(SCList, data.Name)
			}
			if data.OverrideValuesCount > 0 {
				ovrData, err := getOvrData(SCID, name, data.Name, cfg)
				if err != nil {
					logger.Warning.Println("Host group dont have a overrides, ", SCID, name, data.Name)
				} else {
					OvrList = append(OvrList, ovrData)
				}
			}
		}

		pClasses[res.Class] = append(pClasses[res.Class], models.PuppetClassesWeb{
			Subclass:     res.Subclass,
			SmartClasses: SCList,
			Overrides:    OvrList,
		})
	}

	stmt.Close()

	return models.HGElem{
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
func insertHG(name string, host string, data string, foremanId int, cfg *models.Config) int64 {
	hgExist := checkHG(name, host, cfg)
	if hgExist == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into hg(name, host, dump, created_at, updated_at, foreman_id, pcList, locList) values(?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			logger.Error.Println(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(name, host, data, time.Now(), time.Now(), foremanId, "NULL", "NULL")
		if err != nil {
			logger.Error.Println(err)
		}

		lastID, _ := res.LastInsertId()
		return lastID
	} else {
		stmt, err := cfg.Database.DB.Prepare("UPDATE `goFsync`.`hg` SET `foreman_id` = ?, `updated_at` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Error.Println(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(foremanId, time.Now(), hgExist)
		if err != nil {
			logger.Error.Println(err)
		}

		return int64(hgExist)
	}
}

func insertHGP(sweId int64, name string, pVal string, priority int, cfg *models.Config) {

	oldId := checkParams(sweId, name, cfg)
	if oldId == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into hg_parameters(hg_id, name, `value`, priority) values(?, ?, ?, ?)")
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

func updateLocInHG(hgId int64, lIdList []int64, cfg *models.Config) {

	db := *cfg.Database.DB

	var strLocList []string

	for _, i := range lIdList {
		if i != 0 {
			strLocList = append(strLocList, utils.String(i))
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
func updatePCinHG(hgId int64, pcList []int64, cfg *models.Config) {

	var strPcList []string

	for _, i := range pcList {
		if i != 0 {
			strPcList = append(strPcList, utils.String(i))
		}
	}
	pcListStr := strings.Join(strPcList, ",")
	stmt, err := cfg.Database.DB.Prepare("update hg set pcList=? where id=?")
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
func deleteHGbyId(hgId int, cfg *models.Config) {
	stmt, err := cfg.Database.DB.Prepare("DELETE FROM hg WHERE id=?")
	if err != nil {
		logger.Error.Println(err)
	}

	_, err = stmt.Exec(hgId)
	if err != nil {
		logger.Error.Println(err)
	}

	stmt.Close()
}
