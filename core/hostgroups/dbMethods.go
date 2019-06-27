package hostgroups

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"time"
)

// ======================================================
// CHECKS
// ======================================================
// Check HG by name
func CheckHG(name string, host string, ss *models.Session) int {

	stmt, err := ss.Config.Database.DB.Prepare("select id from hg where name=? and host=?")
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
func StateID(hgName string, ss *models.Session) int {

	stmt, err := ss.Config.Database.DB.Prepare("SELECT id FROM goFsync.hg_state where host_group=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(hgName).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}
func CheckHGID(name string, host string, ss *models.Session) int {

	stmt, err := ss.Config.Database.DB.Prepare("select foreman_id from hg where name=? and host=?")
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
func CheckParams(hgId int, name string, ss *models.Session) int {

	stmt, err := ss.Config.Database.DB.Prepare("select id from hg_parameters where hg_id=? and name=?")
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

func CheckHost(host string, cfg *models.Config) int {
	stmt, err := cfg.Database.DB.Prepare("select id from hosts where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(host).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func HostEnv(host string, ss *models.Session) string {

	stmt, err := ss.Config.Database.DB.Prepare("select env from hosts where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	var hostEnv string
	err = stmt.QueryRow(host).Scan(&hostEnv)
	if err != nil {
		logger.Warning.Println(err)
		return ""
	}

	return hostEnv
}
func AllHosts(ss *models.Session) []models.ForemanHost {
	var result []models.ForemanHost
	stmt, err := ss.Config.Database.DB.Prepare("select host, env from hosts")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Warning.Println(err)
	}
	for rows.Next() {
		var name string
		var env string
		err = rows.Scan(&name, &env)
		if err != nil {
			logger.Error.Println(err)
		}
		if logger.StringInSlice(name, ss.Config.Hosts) {
			result = append(result, models.ForemanHost{
				Name: name,
				Env:  env,
			})
		}
	}
	return result
}
func GetHGAllList(ss *models.Session) []models.HGListElem {

	stmt, err := ss.Config.Database.DB.Prepare("select id, name from hg")
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
func GetHGList(host string, ss *models.Session) []models.HGListElem {

	stmt, err := ss.Config.Database.DB.Prepare("select id, foreman_id, name, status from hg where host=?")
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
		var foremanId int
		var name string
		var status string
		err = rows.Scan(&id, &foremanId, &name, &status)
		if err != nil {
			logger.Error.Println(err)
		}
		list = append(list, models.HGListElem{
			ID:        id,
			ForemanID: foremanId,
			Name:      name,
			Status:    status,
		})
	}

	return list
}

func GetHGParams(hgId int, ss *models.Session) []models.HGParam {

	stmt, err := ss.Config.Database.DB.Prepare("select foreman_id, name, value from hg_parameters where hg_id=?")
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
		var foremanId int
		err = rows.Scan(&foremanId, &name, &value)
		if err != nil {
			logger.Error.Println(err)
		}
		list = append(list, models.HGParam{
			ForemanID: foremanId,
			Name:      name,
			Value:     value,
		})
	}

	return list
}

func GetHG(id int, ss *models.Session) models.HGElem {

	// VARS
	var d models.HostGroup
	var name string
	var status string
	var pClassesStr string
	var dump string
	var foremanId int
	var updatedStr string
	pClasses := make(map[string][]models.PuppetClassesWeb)

	// Hg Data
	stmt, err := ss.Config.Database.DB.Prepare("select foreman_id, name, pcList, status, dump, updated_at from hg where id=?")
	if err != nil {
		logger.Warning.Println("HostGroup getting..", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&foremanId, &name, &pClassesStr, &status, &dump, &updatedStr)
	if err != nil {
		return models.HGElem{}
	}

	// HG Parameters
	params := GetHGParams(id, ss)

	err = json.Unmarshal([]byte(dump), &d)
	if err != nil {
		logger.Warning.Printf("Error on Parsing HG: %s", err)
	}

	// PuppetClasses and Parameters
	for _, cl := range utils.Integers(pClassesStr) {
		res := puppetclass.DbByID(cl, ss)

		var SCList []models.SmartClass
		var OvrList []models.SCOParams
		scList := utils.Integers(res.SCIDs)
		for _, SCID := range scList {
			data := smartclass.GetSCData(SCID, ss)
			if data.Name != "" {
				SCList = append(SCList, models.SmartClass{
					Id:        data.ID,
					ForemanId: data.ForemanId,
					Name:      data.Name,
				})
			}
			if data.OverrideValuesCount > 0 {
				ovrData, err := smartclass.GetOvrData(SCID, name, data.Name, ss)
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
		Status:        status,
		Params:        params,
		Environment:   d.EnvironmentName,
		ParentId:      d.Ancestry,
		PuppetClasses: pClasses,
		Updated:       updatedStr,
	}

}

func GetForemanIDs(host string, ss *models.Session) []int {
	var result []int

	stmt, err := ss.Config.Database.DB.Prepare("SELECT foreman_id FROM hg WHERE host=?;")
	if err != nil {
		logger.Warning.Printf("%q, GetForemanIDs", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(host)
	if err != nil {
		logger.Warning.Printf("%q, GetForemanIDs", err)
	}
	for rows.Next() {
		var _id int
		err = rows.Scan(&_id)
		if err != nil {
			logger.Warning.Printf("%q, GetForemanIDs", err)
		}

		result = append(result, _id)
	}
	return result
}

// ======================================================
// INSERT
// ======================================================
func Insert(name string, host string, data string, sweStatus string, foremanId int, ss *models.Session) int {
	hgExist := CheckHG(name, host, ss)
	if hgExist == -1 {
		stmt, err := ss.Config.Database.DB.Prepare("insert into hg(name, host, dump, created_at, updated_at, foreman_id, pcList, status) values(?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(name, host, data, time.Now(), time.Now(), foremanId, "NULL", sweStatus)
		if err != nil {
			return -1
		}

		lastID, _ := res.LastInsertId()
		return int(lastID)
	} else {
		stmt, err := ss.Config.Database.DB.Prepare("UPDATE hg SET  `status` = ?, `foreman_id` = ?, `updated_at` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(sweStatus, foremanId, time.Now(), hgExist)
		if err != nil {
			return -1
		}

		return hgExist
	}
}

func InsertParameters(sweId int, p models.HostGroupP, ss *models.Session) {

	oldId := CheckParams(sweId, p.Name, ss)
	if oldId == -1 {
		stmt, err := ss.Config.Database.DB.Prepare("insert into hg_parameters(hg_id, foreman_id, name, `value`, priority) values(?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(sweId, p.ID, p.Name, p.Value, p.Priority)
		if err != nil {
			logger.Warning.Println(err)
		}
	} else {
		stmt, err := ss.Config.Database.DB.Prepare("UPDATE `goFsync`.`hg_parameters` SET `foreman_id` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(p.ID, oldId)
		if err != nil {
			logger.Warning.Println(err)
		}
	}
}

func InsertHost(host string, cfg *models.Config) {
	if id := CheckHost(host, cfg); id == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into hosts (host) values(?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(host)
		if err != nil {
			logger.Warning.Println(err)
		}
	}
}

func insertState(hgName, host, state string, ss *models.Session) {
	ID := StateID(hgName, ss)
	if ID == -1 {
		q := fmt.Sprintf("insert into hg_state (host_group, `%s`) values(?, ?)", host)
		stmt, err := ss.Config.Database.DB.Prepare(q)
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(hgName, state)
		if err != nil {
			logger.Warning.Println(err)
		}
	} else {
		q := fmt.Sprintf("UPDATE `goFsync`.`hg_state` SET `%s` = ? WHERE (`id` = ?)", host)
		stmt, err := ss.Config.Database.DB.Prepare(q)
		if err != nil {
			logger.Warning.Println(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(state, ID)
		if err != nil {
			logger.Warning.Println(err)
		}
	}

}

// ======================================================
// UPDATE
// ======================================================

// ======================================================
// DELETE
// ======================================================
func DeleteHGbyId(hgId int, ss *models.Session) {
	stmt, err := ss.Config.Database.DB.Prepare("DELETE FROM hg WHERE id=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(hgId)
	if err != nil {
		logger.Warning.Println(err)
	}
}
