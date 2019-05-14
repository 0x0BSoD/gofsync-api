package puppetclass

import (
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"strconv"
	"strings"
)

// ======================================================
// CHECKS
// ======================================================
func CheckPC(subclass string, host string, cfg *models.Config) int64 {

	var id int64

	stmt, err := cfg.Database.DB.Prepare("select id from puppet_classes where host=? and subclass=?")
	if err != nil {
		logger.Warning.Printf("%q, checkPC", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(host, subclass).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func GetAllPCDB(host string, cfg *models.Config) []models.PCintId {

	var res []models.PCintId

	stmt, err := cfg.Database.DB.Prepare("SELECT id, foreman_id, class, subclass, sc_ids from goFsync.puppet_classes where host=?;")
	if err != nil {
		logger.Warning.Printf("%q, getByNamePC", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(host)
	if err != nil {
		return []models.PCintId{}
	}

	for rows.Next() {
		var foremanId int
		var class string
		var subclass string
		var scIds string
		var _id int
		err := rows.Scan(&_id, &foremanId, &class, &subclass, &scIds)
		if err != nil {
			logger.Warning.Println("No result while getting puppet classes")
		}

		if scIds != "" {
			intScIds := logger.Integers(scIds)
			res = append(res, models.PCintId{
				ID:        _id,
				ForemanId: foremanId,
				Class:     class,
				Subclass:  subclass,
				SCIDs:     intScIds,
			})
		} else {
			res = append(res, models.PCintId{
				ForemanId: foremanId,
				Class:     class,
				Subclass:  subclass,
			})
		}
	}

	return res
}

func GetByNamePC(subclass string, host string, cfg *models.Config) models.PC {

	var class string
	var sCIDs string
	var envIDs string
	var hGIDs string
	var foremanId int
	var id int

	stmt, err := cfg.Database.DB.Prepare("select id, class, sc_ids, env_ids, hg_ids, foreman_id from puppet_classes where subclass=? and host=?")
	if err != nil {
		logger.Warning.Printf("%q, getByNamePC", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(subclass, host).Scan(&id, &class, &sCIDs, &envIDs, &hGIDs, &foremanId)
	if err != nil {
		return models.PC{}
	}

	return models.PC{
		ID:        id,
		ForemanId: foremanId,
		Class:     class,
		Subclass:  subclass,
		SCIDs:     sCIDs,
	}
}
func GetPC(pId int, cfg *models.Config) models.PC {

	var class string
	var subclass string
	var sCIDs string
	var envIDs string
	var hGIDs string

	stmt, err := cfg.Database.DB.Prepare("select class, subclass, sc_ids, env_ids, hg_ids from puppet_classes where id=?")
	if err != nil {
		logger.Warning.Printf("%q, getPC", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(pId).Scan(&class, &subclass, &sCIDs, &envIDs, &hGIDs)

	return models.PC{
		Class:    class,
		Subclass: subclass,
		SCIDs:    sCIDs,
	}
}

func GetAllPCBase(host string, cfg *models.Config) []models.PuppetclassesNI {

	var r []models.PuppetclassesNI

	stmt, err := cfg.Database.DB.Prepare("select foreman_id, class, subclass from puppet_classes where host=?")
	if err != nil {
		logger.Warning.Printf("%q, getAllPCBase", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(host)
	if err != nil {
		return []models.PuppetclassesNI{}
	}
	for rows.Next() {
		var foremanId int
		var class string
		var subClass string
		err = rows.Scan(&foremanId, &class, &subClass)
		if err != nil {
			logger.Warning.Printf("%q, getAllPCBase", err)
		}
		r = append(r, models.PuppetclassesNI{class, subClass, foremanId})
	}

	return r
}

// ======================================================
// INSERT
// ======================================================
func InsertPC(host string, class string, subclass string, foremanId int, cfg *models.Config) int64 {

	existID := CheckPC(subclass, host, cfg)
	if existID == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into puppet_classes(host, class, subclass, foreman_id, sc_ids, env_ids) values(?,?,?,?,?,?)")
		if err != nil {
			logger.Warning.Printf("%q, insertPC", err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(host, class, subclass, foremanId, "NULL", "NULL")
		if err != nil {
			logger.Warning.Printf("%q, checkPC", err)
		}

		lastID, _ := res.LastInsertId()
		return lastID
	} else {
		return existID
	}
}

func UpdatePC(host string, subClass string, data models.PCSCParameters, cfg *models.Config) {

	var strScList []string
	var strEnvList []string
	var strHGList []string

	for _, i := range data.SmartClassParameters {
		scID := smartclass.CheckSC(data.Name, i.Name, host, cfg)
		if scID != -1 {
			strScList = append(strScList, strconv.Itoa(int(scID)))
		}
	}

	stmt, err := cfg.Database.DB.Prepare("update puppet_classes set sc_ids=?, env_ids=?, hg_ids=? where host=? and subclass=?")
	if err != nil {
		logger.Warning.Printf("%q, updatePC", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		strings.Join(strScList, ","),
		strings.Join(strEnvList, ","),
		strings.Join(strHGList, ","),
		host,
		subClass)
	if err != nil {
		logger.Warning.Printf("%q, updatePC", err)
	}
}

func UpdatePCinHG(hgId int64, pcList []int64, cfg *models.Config) {

	var strPcList []string

	for _, i := range pcList {
		if i != 0 {
			strPcList = append(strPcList, utils.String(i))
		}
	}
	pcListStr := strings.Join(strPcList, ",")
	stmt, err := cfg.Database.DB.Prepare("update hg set pcList=? where id=?")
	if err != nil {
		logger.Error.Println(err)
	}

	_, err = stmt.Exec(pcListStr, hgId)
	if err != nil {
		logger.Error.Println(err)
	}

	stmt.Close()
}

// ======================================================
// DELETE
// ======================================================
func DeletePuppetClass(host string, subClass string, cfg *models.Config) {
	stmt, err := cfg.Database.DB.Prepare("DELETE FROM `goFsync`.`puppet_classes` WHERE (`host` = ? and `subclass`=?);")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(host, subClass)
	if err != nil {
		logger.Warning.Printf("%q, DeletePuppetClass", err)
	}
}
