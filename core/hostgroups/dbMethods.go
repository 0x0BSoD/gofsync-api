package hostgroups

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/hosts"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"time"
)

// =====================================================================================================================
// IDS
// =====================================================================================================================

// Return DB ID for host group
func ID(name, host string, ctx *user.GlobalCTX) int {
	stmt, err := ctx.Config.Database.DB.Prepare("select id from hg where name=? and host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

// Return Foreman ID for host group
func FID(name, host string, ctx *user.GlobalCTX) int {
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id from hg where name=? and host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

// Return Foreman ID for puppet master host
func FIDs(host string, ctx *user.GlobalCTX) []int {
	var result []int

	stmt, err := ctx.Config.Database.DB.Prepare("SELECT foreman_id FROM hg WHERE host=?;")
	if err != nil {
		logger.Warning.Printf("%q, GetForemanIDs", err)
	}
	defer utils.DeferCloseStmt(stmt)

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

// Return DB ID for host group parameter
func PID(hgId int, name string, ctx *user.GlobalCTX) int {
	stmt, err := ctx.Config.Database.DB.Prepare("select id from hg_parameters where hg_id=? and name=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)
	var id int
	err = stmt.QueryRow(hgId, name).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

// Return DB ID for puppet master host parameter
func HID(host string, cfg *models.Config) int {
	stmt, err := cfg.Database.DB.Prepare("select id from hosts where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)
	var id int
	err = stmt.QueryRow(host).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// =====================================================================================================================
// GET
// =====================================================================================================================

// Return all puppet master hosts with environments
func PuppetHosts(ctx *user.GlobalCTX) []hosts.ForemanHost {
	var result []hosts.ForemanHost
	stmt, err := ctx.Config.Database.DB.Prepare("select host, env from hosts")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

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
		if logger.StringInSlice(name, ctx.Config.Hosts) {
			result = append(result, hosts.ForemanHost{
				Name: name,
				Env:  env,
			})
		}
	}
	return result
}

// Return Environment for puppet master host
func PuppetHostEnv(host string, ctx *user.GlobalCTX) string {
	stmt, err := ctx.Config.Database.DB.Prepare("select env from hosts where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var hostEnv string
	err = stmt.QueryRow(host).Scan(&hostEnv)
	if err != nil {
		logger.Warning.Println(err)
		return ""
	}

	return hostEnv
}

// Return all host groups
func All(ctx *user.GlobalCTX) []HGListElem {
	stmt, err := ctx.Config.Database.DB.Prepare("select id, name from hg")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var list []HGListElem
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
			list = append(list, HGListElem{
				ID:   id,
				Name: name,
			})
		}

	}

	return list
}

// Return all host groups for puppet master host
func OnHost(host string, ctx *user.GlobalCTX) []HGListElem {
	stmt, err := ctx.Config.Database.DB.Prepare("select id, foreman_id, name, status from hg where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var list []HGListElem

	rows, err := stmt.Query(host)
	if err != nil {
		return list
	}

	for rows.Next() {
		var (
			id, foremanId int
			name, status  string
		)
		err = rows.Scan(&id, &foremanId, &name, &status)
		if err != nil {
			logger.Error.Println(err)
		}
		list = append(list, HGListElem{
			ID:        id,
			ForemanID: foremanId,
			Name:      name,
			Status:    status,
		})
	}

	return list
}

// Return host group parameters by hg id
func HGParams(hgId int, ctx *user.GlobalCTX) []HGParam {
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id, name, value from hg_parameters where hg_id=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var list []HGParam

	rows, err := stmt.Query(hgId)
	if err != nil {
		return list
	}

	for rows.Next() {
		var (
			foremanId   int
			name, value string
		)
		err = rows.Scan(&foremanId, &name, &value)
		if err != nil {
			logger.Error.Println(err)
		}
		list = append(list, HGParam{
			ForemanID: foremanId,
			Name:      name,
			Value:     value,
		})
	}

	return list
}

// Get host group by DB ID
func Get(id int, ctx *user.GlobalCTX) HGElem {
	var (
		foremanId                                   int
		name, status, pClassesStr, dump, updatedStr string
		d                                           HostGroupForeman
	)
	pClasses := make(map[string][]puppetclass.PuppetClassesWeb)

	// Hg Data
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id, name, pcList, status, dump, updated_at from hg where id=?")
	if err != nil {
		logger.Warning.Println("HostGroup getting..", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(id).Scan(&foremanId, &name, &pClassesStr, &status, &dump, &updatedStr)
	if err != nil {
		return HGElem{}
	}

	err = json.Unmarshal([]byte(dump), &d)
	if err != nil {
		logger.Warning.Printf("Error on Parsing HG: %s", err)
	}

	// PuppetClasses and Parameters
	for _, cl := range utils.Integers(pClassesStr) {
		res := puppetclass.DbByID(cl, ctx)

		var SCList []smartclass.SmartClass
		var OvrList []smartclass.SCOParams

		scList := utils.Integers(res.SCIDs)
		for _, SCID := range scList {
			data := smartclass.GetSCData(SCID, ctx)

			if len(data.Name) > 0 {
				SCList = append(SCList, smartclass.SmartClass{
					Id:        data.ID,
					ForemanId: data.ForemanID,
					Name:      data.Name,
				})
			}
			if data.OverrideValuesCount > 0 {
				ovrData, err := smartclass.GetOvrData(SCID, name, data.Name, ctx)
				if err != nil {
					logger.Trace.Println("Host group dont have a overrides, ", SCID, name, data.Name)
				} else {
					OvrList = append(OvrList, ovrData)
				}
			}
		}

		pClasses[res.Class] = append(pClasses[res.Class], puppetclass.PuppetClassesWeb{
			Subclass:     res.Subclass,
			SmartClasses: SCList,
			Overrides:    OvrList,
		})
	}

	return HGElem{
		ID:            id,
		ForemanID:     foremanId,
		Name:          name,
		Status:        status,
		Params:        HGParams(id, ctx),
		Environment:   d.EnvironmentName,
		ParentId:      d.Ancestry,
		PuppetClasses: pClasses,
		Updated:       updatedStr,
	}

}

// =====================================================================================================================
// INSERT
// =====================================================================================================================

// Insert/Update host group
func Insert(name, host, data, sweStatus string, foremanId int, ctx *user.GlobalCTX) int {
	hgID := ID(name, host, ctx)
	if hgID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into hg(name, host, dump, created_at, updated_at, foreman_id, pcList, status) values(?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		res, err := stmt.Exec(name, host, data, time.Now(), time.Now(), foremanId, "NULL", sweStatus)
		if err != nil {
			return -1
		}

		lastID, _ := res.LastInsertId()
		return int(lastID)
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE hg SET  `status` = ?, `foreman_id` = ?, `updated_at` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(sweStatus, foremanId, time.Now(), hgID)
		if err != nil {
			return -1
		}

		return hgID
	}
}

// Insert/Update host group parameters
func InsertParameters(sweId int, p HostGroupP, ctx *user.GlobalCTX) {
	PID := PID(sweId, p.Name, ctx)
	if PID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into hg_parameters(hg_id, foreman_id, name, `value`, priority) values(?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(sweId, p.ID, p.Name, p.Value, p.Priority)
		if err != nil {
			logger.Warning.Println(err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE hg_parameters SET `foreman_id` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(p.ID, PID)
		if err != nil {
			logger.Warning.Println(err)
		}
	}
}

// Insert puppet master host
func InsertHost(host string, cfg *models.Config) {
	if id := HID(host, cfg); id == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into hosts (host) values(?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)
		_, err = stmt.Exec(host)
		if err != nil {
			logger.Warning.Println(err)
		}
	}
}

//func insertState(hgName, host, state string, ctx *user.GlobalCTX) {
//	ID := StateID(hgName, ctx)
//	if ID == -1 {
//		q := fmt.Sprintf("insert into hg_state (host_group, `%s`) values(?, ?)", host)
//		stmt, err := ctx.Config.Database.DB.Prepare(q)
//		if err != nil {
//			logger.Warning.Println(err)
//		}
//		defer utils.DeferCloseStmt(stmt)
//		_, err = stmt.Exec(hgName, state)
//		if err != nil {
//			logger.Warning.Println(err)
//		}
//	} else {
//		q := fmt.Sprintf("UPDATE `goFsync`.`hg_state` SET `%s` = ? WHERE (`id` = ?)", host)
//		stmt, err := ctx.Config.Database.DB.Prepare(q)
//		if err != nil {
//			logger.Warning.Println(err)
//		}
//		defer utils.DeferCloseStmt(stmt)
//		_, err = stmt.Exec(state, ID)
//		if err != nil {
//			logger.Warning.Println(err)
//		}
//	}
//
//}

// =====================================================================================================================
// UPDATE
// =====================================================================================================================

// =====================================================================================================================
// DELETE
// =====================================================================================================================
//func DeleteHGbyId(hgId int, ctx *user.GlobalCTX) {
//	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM hg WHERE id=?")
//	if err != nil {
//		logger.Warning.Println(err)
//	}
//	defer utils.DeferCloseStmt(stmt)
//
//	_, err = stmt.Exec(hgId)
//	if err != nil {
//		logger.Warning.Println(err)
//	}
//}

func Delete(foremanID int, host string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM hg WHERE foreman_id=? AND host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(foremanID, host)
	if err != nil {
		logger.Warning.Println(err)
	}
}
