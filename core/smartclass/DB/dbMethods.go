package DB

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"strconv"
)

//func CheckOvr(scId int, match string, ctx *user.GlobalCTX) int {
//
//	var id int
//	//fmt.Printf("select id from override_values where sc_id=%d and `match`=%s\n", scId, match)
//	stmt, err := ctx.Config.Database.DB.Prepare("select id from override_values where sc_id=? and `match`=?")
//	if err != nil {
//		logger.Warning.Printf("%q, checkSC", err)
//	}
//	defer utils.DeferCloseStmt(stmt)
//
//	err = stmt.QueryRow(scId, match).Scan(&id)
//	if err != nil {
//		return -1
//	}
//	return id
//}

// ======================================================
// GET
// ======================================================

func GetForemanIDs(host string, ctx *user.GlobalCTX) []int {
	var result []int

	stmt, err := ctx.Config.Database.DB.Prepare("SELECT foreman_id FROM smart_classes WHERE host=?;")
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

func GetForemanIDsBySCid(scId int, ctx *user.GlobalCTX) []int {
	var result []int

	stmt, err := ctx.Config.Database.DB.Prepare("SELECT foreman_id FROM override_values WHERE sc_id=?;")
	if err != nil {
		logger.Warning.Printf("%q, GetOverrodesForemanIDs", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(scId)
	if err != nil {
		logger.Warning.Printf("%q, GetOverrodesForemanIDs", err)
	}
	for rows.Next() {
		var _id int
		err = rows.Scan(&_id)
		if err != nil {
			logger.Warning.Printf("%q, GetOverrodesForemanIDs", err)
		}

		result = append(result, _id)
	}
	return result
}

// ======================================================
// INSERT
// ======================================================
func InsertSC(host string, data SCParameter, ctx *user.GlobalCTX) {

	var dbId int

	existID := CheckSCByForemanId(host, data.ID, ctx)
	if existID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into smart_classes(host, puppetclass, parameter, parameter_type, foreman_id, override_values_count, dump) values(?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Printf("%q, insertSC", err)
		}
		defer utils.DeferCloseStmt(stmt)

		sJson, _ := json.Marshal(data)
		res, err := stmt.Exec(host, data.PuppetClass.Name, data.Parameter, data.ParameterType, data.ID, data.OverrideValuesCount, sJson)
		if err != nil {
			logger.Warning.Printf("%q, insertSC", err)
		}

		lastId, _ := res.LastInsertId()
		dbId = int(lastId)
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE smart_classes SET `override_values_count` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Printf("%q, updateSC", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(data.OverrideValuesCount, existID)
		if err != nil {
			logger.Warning.Printf("%q, updateSC", err)
		}
		dbId = existID
	}

	if data.OverrideValuesCount > 0 {

		//fmt.Println(utils.PrintJsonStep(models.Step{
		//	Actions: fmt.Sprintf("Storing Smart classes Overrides %s", data.Parameter),
		//	Host:    host,
		//}))

		beforeUpdateOvr := GetForemanIDsBySCid(dbId, ctx)
		var afterUpdateOvr []int
		for _, ovr := range data.OverrideValues {
			afterUpdateOvr = append(afterUpdateOvr, ovr.ID)
			InsertSCOverride(dbId, ovr, data.ParameterType, ctx)
		}

		for _, j := range beforeUpdateOvr {

			//fmt.Println(utils.PrintJsonStep(models.Step{
			//	Actions: fmt.Sprintf("Checking Overrides ... %s", data.Parameter),
			//	Host:    host,
			//}))

			if !utils.Search(afterUpdateOvr, j) {

				//fmt.Println(utils.PrintJsonStep(models.Step{
				//	Actions: fmt.Sprintf("Deleting Overrides ... %s", data.Parameter),
				//	Host:    host,
				//}))

				DeleteOverride(dbId, j, ctx)
			}
		}

	}

}

// Insert Smart Class override
func InsertSCOverride(scId int, data OverrideValue, pType string, ctx *user.GlobalCTX) {

	var strData string

	// Value assertion
	// =================================================================================================================
	if data.Value != nil {
		switch data.Value.(type) {
		case string:
			strData = data.Value.(string)
		case []interface{}:
			var tmpResInt []string
			for _, i := range data.Value.([]interface{}) {
				tmpResInt = append(tmpResInt, i.(string))
			}
			strIng, _ := json.Marshal(tmpResInt)
			strData = string(strIng)
		case bool:
			strData = string(strconv.FormatBool(data.Value.(bool)))
		case int:
			strData = strconv.FormatFloat(data.Value.(float64), 'f', 6, 64)
		case float64:
			strData = strconv.FormatFloat(data.Value.(float64), 'f', 6, 64)
		default:
			logger.Warning.Printf("Type not known try save as string, Type: %s, Val: %s, Match: %s", pType, data.Value, data.Match)
			strData = data.Value.(string)
		}
	}
	// =================================================================================================================

	existId := CheckOvrByForemanId(scId, data.ID, ctx)
	if existId == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into override_values(foreman_id, `match`, value, sc_id, use_puppet_default) values(?,?,?,?,?)")
		if err != nil {
			logger.Warning.Printf("%q, insertSCOverride", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(data.ID, data.Match, strData, scId, data.UsePuppetDefault)
		if err != nil {
			logger.Warning.Printf("%q, insertSCOverride", err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE override_values SET `value` = ?, foreman_id=? WHERE id= ?")
		if err != nil {
			logger.Warning.Printf("%q, Prepare updateSCOverride data: %q, %d", err, strData, existId)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(strData, data.ID, existId)
		if err != nil {
			logger.Warning.Printf("%q, Exec updateSCOverride data: %q, %d", err, strData, existId)
		}
	}
}

// ======================================================
// DELETE
// ======================================================
func DeleteSmartClass(host string, foremanId int, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM smart_classes WHERE host=? and foreman_id=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, foremanId)
	if err != nil {
		logger.Warning.Printf("%q, DeleteSmartClass	", err)
	}
}

func DeleteOverride(scId int, foremanId int, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM override_values WHERE sc_id=? AND foreman_id=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(scId, foremanId)
	if err != nil {
		logger.Warning.Printf("%q, DeleteSmartClass	", err)
	}
}
