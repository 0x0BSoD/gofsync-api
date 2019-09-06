package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"strconv"
)

// ======================================================
// CHECKS
// ======================================================
func CheckSCByForemanId(host string, foremanId int, ctx *user.GlobalCTX) int {

	var id int

	stmt, err := ctx.Config.Database.DB.Prepare("select id from smart_classes where host=? and foreman_id=?")
	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, foremanId).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

func CheckOvrByForemanId(scId int, foremanId int, ctx *user.GlobalCTX) int {

	var id int

	stmt, err := ctx.Config.Database.DB.Prepare("select id from override_values where sc_id=? and foreman_id=?")
	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(scId, foremanId).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func GetSC(host string, puppetClass string, parameter string, ctx *user.GlobalCTX) SCGetResAdv {

	var id int
	var foremanId int
	var ovrCount int

	stmt, err := ctx.Config.Database.DB.Prepare("select id, override_values_count, foreman_id from smart_classes where parameter=? and puppetclass=? and host=?")

	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(parameter, puppetClass, host).Scan(&id, &ovrCount, &foremanId)
	if err != nil {
		return SCGetResAdv{}
	}

	return SCGetResAdv{
		ID:                  id,
		ForemanID:           foremanId,
		Name:                parameter,
		OverrideValuesCount: ovrCount,
	}
}

func GetSCData(scID int, ctx *user.GlobalCTX) SCGetResAdv {

	stmt, err := ctx.Config.Database.DB.Prepare("select id, parameter, override_values_count, foreman_id, parameter_type, puppetclass, dump from smart_classes where id=?")
	if err != nil {
		logger.Warning.Printf("%q, getSCData", err)
	}
	defer utils.DeferCloseStmt(stmt)

	var (
		id        int
		foremanId int
		paramName string
		ovrCount  int
		_type     string
		pc        string
		dump      string
	)

	err = stmt.QueryRow(scID).Scan(&id, &paramName, &ovrCount, &foremanId, &_type, &pc, &dump)
	if err != nil {
		return SCGetResAdv{}
	}

	return SCGetResAdv{
		ID:                  id,
		ForemanID:           foremanId,
		Name:                paramName,
		OverrideValuesCount: ovrCount,
		ValueType:           _type,
		PuppetClass:         pc,
		Dump:                dump,
	}
}

func GetOvrData(scId int, name string, parameter string, ctx *user.GlobalCTX) (SCOParams, error) {

	matchStr := fmt.Sprintf("hostgroup=SWE/%s", name)
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id, `match`, value, sc_id from override_values where sc_id=? and `match` = ?")
	if err != nil {
		logger.Warning.Printf("%q, getOvrData", err)
	}
	defer utils.DeferCloseStmt(stmt)

	var foremanId int
	var match string
	var scID int
	var val string

	err = stmt.QueryRow(scId, matchStr).Scan(&foremanId, &match, &val, &scID)
	if err != nil {
		return SCOParams{}, err
	}

	return SCOParams{
		ForemanID:    foremanId,
		SmartClassId: scID,
		Parameter:    parameter,
		Match:        match,
		Value:        val,
	}, nil
}

func GetOverridesHG(hgName string, ctx *user.GlobalCTX) []OvrParams {

	var results []OvrParams

	qStr := fmt.Sprintf("hostgroup=SWE/%s", hgName)
	stmt, err := ctx.Config.Database.DB.Prepare("select `match`, value, sc_id from override_values where `match`= ?")
	if err != nil {
		logger.Warning.Printf("%q, getOverridesHG", err)
	}
	defer utils.DeferCloseStmt(stmt)
	rows, err := stmt.Query(qStr)
	if err != nil {
		logger.Warning.Printf("%q, getOverridesHG", err)
	}
	for rows.Next() {
		var smartClassId int
		var value string
		var match string
		err = rows.Scan(&match, &value, &smartClassId)
		if err != nil {
			logger.Warning.Printf("%q, getOverridesHG", err)
		}
		scData := GetSCData(smartClassId, ctx)
		results = append(results, OvrParams{
			SmartClassName: scData.Name,
			Value:          value,
		})
	}

	return results
}

func GetOverridesLoc(host, locName string, ctx *user.GlobalCTX) []OverrideParameters {

	var results []OverrideParameters
	qStr := fmt.Sprintf("location=%s", locName)
	stmt, err := ctx.Config.Database.DB.Prepare("select  ov.`match`, ov.value, ov.sc_id, ov.foreman_id as ovr_foreman_id, sc.foreman_id  as sc_foreman_id, sc.parameter,sc.parameter_type, sc.puppetclass from override_values as ov, smart_classes as sc where ov.`match`= ? and sc.id = ov.sc_id and sc.host = ?")
	if err != nil {
		logger.Warning.Printf("%q, getOverridesLoc", err)
	}
	defer utils.DeferCloseStmt(stmt)
	rows, err := stmt.Query(qStr, host)
	if err != nil {
		logger.Warning.Printf("%q, getOverridesLoc", err)
	}

	resTmp := make(map[string][]OvrParams)

	for rows.Next() {
		var ovrFId int
		var scFId int
		var param string
		var _type string
		var pc string
		var smartClassId int
		var value string
		var match string
		err = rows.Scan(&match, &value, &smartClassId,
			&ovrFId, &scFId, &param, &_type, &pc)
		if err != nil {
			logger.Warning.Printf("%q, getOverridesLoc", err)
		}

		var dumpObj SCParameterDef
		scData := GetSCData(smartClassId, ctx)
		_ = json.Unmarshal([]byte(scData.Dump), &dumpObj)
		resTmp[pc] = append(resTmp[pc], OvrParams{
			SmartClassName: scData.Name,
			Value:          value,
			OvrForemanId:   ovrFId,
			PuppetClass:    pc,
			SCForemanId:    scFId,
			Type:           _type,
			DefaultValue:   dumpObj.DefaultValue,
		})
	}

	for pc, data := range resTmp {
		var tmp []OverrideParameter
		for _, i := range data {
			tmp = append(tmp, OverrideParameter{
				OverrideForemanId:  i.OvrForemanId,
				ParameterForemanId: i.SCForemanId,
				Name:               i.SmartClassName,
				Value:              i.Value,
				Type:               i.Type,
				DefaultValue:       i.DefaultValue,
			})
		}
		results = append(results, OverrideParameters{
			PuppetClass: pc,
			Parameters:  tmp,
		})
	}

	return results
}

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
