package hostgroups

import (
	"encoding/json"
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
func CheckParams(hgId int, name string, cfg *models.Config) int {

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
func HostEnv(host string, cfg *models.Config) string {

	stmt, err := cfg.Database.DB.Prepare("select env from hosts where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer stmt.Close()

	var hostEnv string
	err = stmt.QueryRow(host).Scan(&hostEnv)
	if err != nil {
		logger.Warning.Println(err)
	}

	return hostEnv
}
func AllHosts(cfg *models.Config) []models.Host {
	var result []models.Host
	stmt, err := cfg.Database.DB.Prepare("select host, env from hosts")
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
		if logger.StringInSlice(name, cfg.Hosts) {
			result = append(result, models.Host{
				Name: name,
				Env:  env,
			})
		}
	}
	return result
}
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

	stmt, err := cfg.Database.DB.Prepare("select id, name, status from hg where host=?")
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
		var status string
		err = rows.Scan(&id, &name, &status)
		if err != nil {
			logger.Error.Println(err)
		}
		list = append(list, models.HGListElem{
			ID:     id,
			Name:   name,
			Status: status,
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
	var status string
	var pClassesStr string
	var dump string
	var foremanId int
	var updatedStr string
	pClasses := make(map[string][]models.PuppetClassesWeb)

	// Hg Data
	stmt, err := cfg.Database.DB.Prepare("select foreman_id, name, pcList, status, dump, updated_at from hg where id=?")
	if err != nil {
		logger.Warning.Println("HostGroup getting..", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&foremanId, &name, &pClassesStr, &status, &dump, &updatedStr)
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
		res := puppetclass.DbByID(cl, cfg)

		var SCList []models.SmartClass
		var OvrList []models.SCOParams
		scList := utils.Integers(res.SCIDs)
		for _, SCID := range scList {
			data := smartclass.GetSCData(SCID, cfg)
			if data.Name != "" {
				SCList = append(SCList, models.SmartClass{
					Id:        data.ID,
					ForemanId: data.ForemanId,
					Name:      data.Name,
				})
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
		Status:        status,
		Params:        params,
		Environment:   d.EnvironmentName,
		ParentId:      d.Ancestry,
		PuppetClasses: pClasses,
		Updated:       updatedStr,
	}

}

func GetForemanIDs(host string, cfg *models.Config) []int {
	var result []int

	stmt, err := cfg.Database.DB.Prepare("SELECT foreman_id FROM hg WHERE host=?;")
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
func Insert(name string, host string, data string, sweStatus string, foremanId int, cfg *models.Config) int {
	hgExist := CheckHG(name, host, cfg)
	if hgExist == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into hg(name, host, dump, created_at, updated_at, foreman_id, pcList, status) values(?, ?, ?, ?, ?, ?, ?, ?)")
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
		stmt, err := cfg.Database.DB.Prepare("UPDATE hg SET  `status` = ?, `foreman_id` = ?, `updated_at` = ? WHERE (`id` = ?)")
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

func InsertParameters(sweId int, name string, pVal string, priority int, cfg *models.Config) {

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

// ======================================================
// UPDATE
// ======================================================

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
