package hostgroups

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/core/puppetclass"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"time"
)

// ======================================================
// CHECKS
// ======================================================
// Check HG by name
func CheckHG(name string, host string, cfg *models.Config) int {

	stmt, err := cfg.Database.DB.Prepare("select id from hg where name=? and host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}
func CheckHGID(name string, host string, cfg *models.Config) int {

	stmt, err := cfg.Database.DB.Prepare("select foreman_id from hg where name=? and host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}
func CheckParams(hgId int64, name string, cfg *models.Config) int {

	stmt, err := cfg.Database.DB.Prepare("select id from hg_parameters where hg_id=? and name=?")
	if err != nil {
		logger.Warning.Println(err)
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
func GetHGAllList(cfg *models.Config) []models.HGListElem {

	stmt, err := cfg.Database.DB.Prepare("select id, name from hg")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	var list []models.HGListElem
	var chkList []string

	rows, err := stmt.Query()
	if err != nil {
		return list
	}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			logger.Error.Println(err)
		}
		if !utils.StringInSlice(name, chkList) {
			chkList = append(chkList, name)
			list = append(list, models.HGListElem{
				ID:   id,
				Name: name,
			})
		}

	}

	return list
}

// For Web Server =======================================
func GetHGList(host string, cfg *models.Config) []models.HGListElem {

	stmt, err := cfg.Database.DB.Prepare("select id, name from hg where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	var list []models.HGListElem

	rows, err := stmt.Query(host)
	if err != nil {
		return list
	}

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			logger.Error.Println(err)
		}
		list = append(list, models.HGListElem{
			ID:   id,
			Name: name,
		})
	}

	return list
}

func GetHGParams(hgId int, cfg *models.Config) []models.HGParam {

	stmt, err := cfg.Database.DB.Prepare("select name, value from hg_parameters where hg_id=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	var list []models.HGParam

	rows, err := stmt.Query(hgId)
	if err != nil {
		return list
	}

	for rows.Next() {
		var name string
		var value string
		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Error.Println(err)
		}
		list = append(list, models.HGParam{
			Name:  name,
			Value: value,
		})
	}

	return list
}

func GetHG(id int, cfg *models.Config) models.HGElem {

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
		logger.Warning.Println("HostGroup getting..", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&foremanId, &name, &pClassesStr, &dump)
	if err != nil {
		return models.HGElem{}
	}

	// HG Parameters
	params := GetHGParams(id, cfg)

	err = json.Unmarshal([]byte(dump), &d)
	if err != nil {
		logger.Warning.Printf("Error on Parsing HG: %s", err)
	}

	// PuppetClasses and Parameters
	for _, cl := range utils.Integers(pClassesStr) {
		res := puppetclass.GetPC(cl, cfg)

		var SCList []string
		var OvrList []models.SCOParams
		scList := utils.Integers(res.SCIDs)
		for _, SCID := range scList {
			data := smartclass.GetSCData(SCID, cfg)
			if data.Name != "" {
				SCList = append(SCList, data.Name)
			}
			if data.OverrideValuesCount > 0 {
				ovrData, err := smartclass.GetOvrData(SCID, name, data.Name, cfg)
				if err != nil {
					logger.Trace.Println("Host group dont have a overrides, ", SCID, name, data.Name)
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
func InsertHG(name string, host string, data string, foremanId int, cfg *models.Config) int64 {
	hgExist := CheckHG(name, host, cfg)
	if hgExist == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into hg(name, host, dump, created_at, updated_at, foreman_id, pcList, locList) values(?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(name, host, data, time.Now(), time.Now(), foremanId, "NULL", "NULL")
		if err != nil {
			return int64(-1)
		}

		lastID, _ := res.LastInsertId()
		return lastID
	} else {
		stmt, err := cfg.Database.DB.Prepare("UPDATE `goFsync`.`hg` SET `foreman_id` = ?, `updated_at` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(foremanId, time.Now(), hgExist)
		if err != nil {
			return int64(-1)
		}

		return int64(hgExist)
	}
}

func InsertHGP(sweId int64, name string, pVal string, priority int, cfg *models.Config) {

	oldId := CheckParams(sweId, name, cfg)
	if oldId == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into hg_parameters(hg_id, name, `value`, priority) values(?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(sweId, name, pVal, priority)
		if err != nil {
			logger.Warning.Println(err)
		}
	}
}

// ======================================================
// DELETE
// ======================================================
func DeleteHGbyId(hgId int, cfg *models.Config) {
	stmt, err := cfg.Database.DB.Prepare("DELETE FROM hg WHERE id=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(hgId)
	if err != nil {
		logger.Warning.Println(err)
	}
}
