package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
)

// ======================================================
// STATEMENTS
// ======================================================
var (
	selectScID           = "select id from smart_classes where host_id=? and foreman_id=?"
	selectOvrID          = "select id from override_values where sc_id=? and foreman_id=?"
	selectSmartClass     = "select id, override_values_count, foreman_id from smart_classes where parameter=? and puppetclass=? and host_id=?"
	selectSmartClassMore = "select id, parameter, override_values_count, foreman_id, parameter_type, puppetclass, dump from smart_classes where id=?"
	selectOvr            = "select foreman_id, `match`, value from override_values where sc_id=? and `match` = ?"
	selectOvrByMatch     = "select `match`, value, sc_id from override_values where `match`= ?"
	selectOvrForLocation = `select ov.'match', 
								   ov.value, 
								   ov.sc_id, 
								   ov.foreman_id as ovr_foreman_id, 
								   sc.foreman_id  as sc_foreman_id, 
								   sc.parameter,
								   sc.parameter_type, 
								   sc.puppetclass from override_values as ov
						    inner join 
								       smart_classes as sc on sc.id = ov.sc_id 
							where ov.'match'= ? and sc.host_id = ?`
	selectAllForemanIDs = "select foreman_id from smart_classes where host_id=?;"

	insertSC  = "insert into smart_classes(host_id, puppetclass, parameter, parameter_type, foreman_id, override, override_values_count, dump) values(?, ?, ?, ?, ?, ?, ?, ?)"
	insertOvr = "insert into override_values(foreman_id, `match`, value, sc_id, use_puppet_default) values(?,?,?,?,?)"

	updateSC  = "update smart_classes set `override_values_count` = ? where (`id` = ?)"
	updateOvr = "update override_values set `value` = ?, foreman_id=? where id= ?"

	deleteSC  = "delete from smart_classes where host_id=? and foreman_id=?"
	deleteOvr = "delete from override_values where sc_id=? and foreman_id=?"
)

// ======================================================
// CHECKS
// ======================================================
func ScID(hostID, foremanID int, ctx *user.GlobalCTX) int {

	var id int

	stmt, err := ctx.Config.Database.DB.Prepare(selectScID)
	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(hostID, foremanID).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

func OvrID(scID, foremanID int, ctx *user.GlobalCTX) int {

	var id int

	stmt, err := ctx.Config.Database.DB.Prepare(selectOvrID)
	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(scID, foremanID).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func GetSC(hostID int, puppetClass, parameter string, ctx *user.GlobalCTX) SCGetResAdv {

	var id int
	var foremanId int
	var ovrCount int

	stmt, err := ctx.Config.Database.DB.Prepare(selectSmartClass)

	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(parameter, puppetClass, hostID).Scan(&id, &ovrCount, &foremanId)
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

	stmt, err := ctx.Config.Database.DB.Prepare(selectSmartClassMore)
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

func GetOvrData(scID int, name, parameter string, ctx *user.GlobalCTX) (SCOParams, error) {

	matchStr := fmt.Sprintf("hostgroup=SWE/%s", name)
	stmt, err := ctx.Config.Database.DB.Prepare(selectOvr)
	if err != nil {
		logger.Warning.Printf("%q, getOvrData", err)
	}
	defer utils.DeferCloseStmt(stmt)

	var foremanId int
	var match string
	var val string

	err = stmt.QueryRow(scID, matchStr).Scan(&foremanId, &match, &val)
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
	stmt, err := ctx.Config.Database.DB.Prepare(selectOvrByMatch)
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

	qStr := fmt.Sprintf("location=%s", locName)
	stmt, err := ctx.Config.Database.DB.Prepare(selectOvrForLocation)
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

	var results = make([]OverrideParameters, 0, len(resTmp))
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

func GetForemanIDs(hostID int, ctx *user.GlobalCTX) []int {

	var result []int

	stmt, err := ctx.Config.Database.DB.Prepare(selectAllForemanIDs)
	if err != nil {
		logger.Warning.Printf("%q, GetForemanIDs", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(hostID)
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
func InsertSC(hostID int, data SCParameter, ctx *user.GlobalCTX) {

	var dbId int

	existID := ScID(hostID, data.ID, ctx)
	if existID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare(insertSC)
		if err != nil {
			logger.Warning.Printf("%q, insertSC", err)
		}
		defer utils.DeferCloseStmt(stmt)

		sJson, _ := json.Marshal(data)
		res, err := stmt.Exec(hostID,
			data.PuppetClass.Name,
			data.Parameter,
			data.ParameterType,
			data.ID,
			data.Override,
			data.OverrideValuesCount,
			sJson)
		if err != nil {
			logger.Warning.Printf("%q, insertSC", err)
		}

		lastId, _ := res.LastInsertId()
		dbId = int(lastId)
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare(updateSC)
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

		fmt.Println(utils.PrintJsonStep(models.Step{
			Actions: "Storing overrides " + data.Parameter,
			Host:    string(hostID),
		}))

		//beforeUpdateOvr := GetForemanIDsBySCid(dbId, ctx)
		aLen := len(data.OverrideValues)
		afterUpdateOvr := make([]int, 0, aLen)

		for _, ovr := range data.OverrideValues {
			afterUpdateOvr = append(afterUpdateOvr, ovr.ID)
			InsertSCOverride(dbId, ovr, data.ParameterType, ctx)
		}

		// TODO: inspect out of index error
		//bLen := len(beforeUpdateOvr)
		//if aLen != bLen {
		//	sort.Ints(afterUpdateOvr)
		//	sort.Ints(beforeUpdateOvr)
		//	for _, j := range beforeUpdateOvr {
		//		if !utils.Search(afterUpdateOvr, j) {
		//			DeleteOverride(dbId, j, ctx)
		//		}
		//	}
		//}

	}

}

// Insert Smart Class override
func InsertSCOverride(scID int, data OverrideValue, pType string, ctx *user.GlobalCTX) {
	strData := utils.AllToStr(data.Value, pType)
	existId := OvrID(scID, data.ID, ctx)
	if existId == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare(insertOvr)
		if err != nil {
			logger.Warning.Printf("%q, insertSCOverride", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(data.ID, data.Match, strData, scID, data.UsePuppetDefault)
		if err != nil {
			logger.Warning.Printf("%q, insertSCOverride", err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare(updateOvr)
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
func DeleteSmartClass(hostID, foremanID int, ctx *user.GlobalCTX) {

	stmt, err := ctx.Config.Database.DB.Prepare(deleteSC)
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(hostID, foremanID)
	if err != nil {
		logger.Warning.Printf("%q, DeleteSmartClass	", err)
	}
}

func DeleteOverride(scId int, foremanId int, ctx *user.GlobalCTX) {

	stmt, err := ctx.Config.Database.DB.Prepare(deleteOvr)
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(scId, foremanId)
	if err != nil {
		logger.Warning.Printf("%q, DeleteSmartClass	", err)
	}
}
