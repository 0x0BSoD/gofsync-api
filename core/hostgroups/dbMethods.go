package hostgroups

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"time"
)

// ======================================================
// STATEMENTS
// ======================================================
var (
	selectID         = "select id         from hg where name=? and host_id=?"
	selectForemanID  = "select foreman_id from hg where name=? and host_id=?"
	selectForemanIDs = "select foreman_id from hg where host_id=?;"
	selectHGName     = "select name       from hg where foreman_id=? and host_id=?"
	selectAll        = "select id, name   from hg"
	selectAllByHost  = "select id, foreman_id, name, status from hg where host_id=?"
	selectByID       = "select foreman_id, name, pcList, status, dump, updated_at from hg where id=?"

	insert   = "insert into hg(name, host_id, dump, created_at, updated_at, foreman_id, pcList, status) values(?, ?, ?, ?, ?, ?, ?, ?)"
	update   = "update hg set  `dump` = ?, `status` = ?, `foreman_id` = ?, `updated_at` = ? where (`id` = ?)"
	deleteHG = "delete from hg where foreman_id=? and host_id=?"

	selectParameterID    = "select id from hg_parameters where hg_id=? and name=?"
	selectParametersByHG = "select foreman_id, name, value from hg_parameters where hg_id=?"

	insertParameter = "insert into hg_parameters(hg_id, foreman_id, name, `value`, priority) values(?, ?, ?, ?, ?)"
	updateParameter = "update hg_parameters set `foreman_id` = ? where (`id` = ?)"
)

// =====================================================================================================================
// IDS
// =====================================================================================================================

// Return DB ID for host group
func ID(hostID int, name string, ctx *user.GlobalCTX) int {
	stmt, err := ctx.Config.Database.DB.Prepare(selectID)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var id int
	err = stmt.QueryRow(name, hostID).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

// Return Foreman ID for host group
func ForemanID(hostID int, hostGroupName string, ctx *user.GlobalCTX) int {
	stmt, err := ctx.Config.Database.DB.Prepare(selectForemanID)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var id int
	err = stmt.QueryRow(hostGroupName, hostID).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

// Return Foreman ID for puppet master host
func ForemanIDs(hostID int, ctx *user.GlobalCTX) []int {
	var result []int

	stmt, err := ctx.Config.Database.DB.Prepare(selectForemanIDs)
	if err != nil {
		utils.Warning.Printf("%q, GetForemanIDs", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(hostID)
	if err != nil {
		utils.Warning.Printf("%q, GetForemanIDs", err)
	}

	for rows.Next() {
		var _id int
		err = rows.Scan(&_id)
		if err != nil {
			utils.Warning.Printf("%q, GetForemanIDs", err)
		}

		result = append(result, _id)
	}

	return result
}

// Return DB ID for host group parameter
func ParameterID(hgID int, name string, ctx *user.GlobalCTX) int {
	stmt, err := ctx.Config.Database.DB.Prepare(selectParameterID)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)
	var id int
	err = stmt.QueryRow(hgID, name).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

// =====================================================================================================================
// GET
// =====================================================================================================================

// Return Host Group name by ID
func Name(hostID, foremanID int, ctx *user.GlobalCTX) string {

	// VARS
	var name string

	// ===========
	stmt, err := ctx.Config.Database.DB.Prepare(selectHGName)
	if err != nil {
		utils.Warning.Println(err)
		return ""
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(foremanID, hostID).Scan(&name)
	if err != nil {
		utils.Warning.Println(err)
		return ""
	}

	return name
}

// Return all host groups
func All(ctx *user.GlobalCTX) []HGListElem {
	stmt, err := ctx.Config.Database.DB.Prepare(selectAll)
	if err != nil {
		utils.Warning.Println(err)
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
			utils.Error.Println(err)
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
func OnHost(hostID int, ctx *user.GlobalCTX) []HGListElem {
	stmt, err := ctx.Config.Database.DB.Prepare(selectAllByHost)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var list []HGListElem

	rows, err := stmt.Query(hostID)
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
			utils.Error.Println(err)
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
	stmt, err := ctx.Config.Database.DB.Prepare(selectParametersByHG)
	if err != nil {
		utils.Warning.Println(err)
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
			utils.Error.Println(err)
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
func Get(ID int, ctx *user.GlobalCTX) HGElem {

	var (
		foremanId                                   int
		name, status, pClassesStr, dump, updatedStr string
		d                                           HostGroupForeman
	)
	pClasses := make(map[string][]puppetclass.PuppetClassesWeb)

	// Hg Data
	stmt, err := ctx.Config.Database.DB.Prepare(selectByID)
	if err != nil {
		utils.Warning.Println("HostGroup getting..", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(ID).Scan(&foremanId, &name, &pClassesStr, &status, &dump, &updatedStr)
	if err != nil {
		return HGElem{}
	}

	err = json.Unmarshal([]byte(dump), &d)
	if err != nil {
		utils.Warning.Printf("Error on Parsing HG: %s", err)
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
					utils.Trace.Println("Host group dont have a overrides, ", SCID, name, data.Name)
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
		ID:            ID,
		ForemanID:     foremanId,
		Name:          name,
		Status:        status,
		Params:        HGParams(ID, ctx),
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
func Insert(hostID, foremanID int, name, data, sweStatus string, ctx *user.GlobalCTX) int {
	hgID := ID(hostID, name, ctx)
	if hgID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare(insert)
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		res, err := stmt.Exec(name, hostID, data, time.Now(), time.Now(), foremanID, "NULL", sweStatus)
		if err != nil {
			return -1
		}

		lastID, _ := res.LastInsertId()
		return int(lastID)
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare(update)
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(data, sweStatus, hostID, time.Now(), hgID)
		if err != nil {
			return -1
		}

		return hgID
	}
}

// Insert/Update host group parameters
func InsertParameters(sweID int, p HostGroupP, ctx *user.GlobalCTX) {
	PID := ParameterID(sweID, p.Name, ctx)
	if PID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare(insertParameter)
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(sweID, p.ID, p.Name, p.Value, p.Priority)
		if err != nil {
			utils.Warning.Println(err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare(updateParameter)
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(p.ID, PID)
		if err != nil {
			utils.Warning.Println(err)
		}
	}
}

// =====================================================================================================================
// DELETE
// =====================================================================================================================

func Delete(hostID, foremanID int, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare(deleteHG)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(foremanID, hostID)
	if err != nil {
		utils.Warning.Println(err)
	}
}
