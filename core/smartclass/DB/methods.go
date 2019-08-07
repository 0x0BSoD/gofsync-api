package DB

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Get Smart Class ID from DB
func (Get) ID(host string, pc string, parameter string, ctx *user.GlobalCTX) int {

	// VARS
	var id int

	// =======
	stmt, err := ctx.Config.Database.DB.Prepare("select id from smart_classes where host=? and parameter=? and puppetclass=?")
	if err != nil {
		utils.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, parameter, pc).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// Get Smart Class ID by Foreman ID and host name from DB
func (Get) IDByForemanID(host string, foremanID int, ctx *user.GlobalCTX) int {

	// VARS
	var id int

	// ======
	stmt, err := ctx.Config.Database.DB.Prepare("select id from smart_classes where host=? and foreman_id=?")
	if err != nil {
		utils.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, foremanID).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// Get Smart Class by ID
func (Get) ByID(scID int, ctx *user.GlobalCTX) (SmartClass, error) {

	stmt, err := ctx.Config.Database.DB.Prepare("select id, parameter, override_values_count, foreman_id, parameter_type, puppetclass, dump from smart_classes where id=?")
	if err != nil {
		utils.Warning.Printf("%q, getSCData", err)
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
		return SmartClass{}, err
	}

	return SmartClass{
		ID:                  id,
		ForemanID:           foremanId,
		Name:                paramName,
		OverrideValuesCount: ovrCount,
		ValueType:           _type,
		PuppetClass:         pc,
		Dump:                dump,
	}, nil
}

// Get Smart Class by parameter and puppet class name
func (Get) GetSC(host string, puppetClass string, parameter string, ctx *user.GlobalCTX) (SmartClass, error) {

	var (
		id        int
		foremanID int
		ovrCount  int
	)

	stmt, err := ctx.Config.Database.DB.Prepare("select id, override_values_count, foreman_id from smart_classes where parameter=? and puppetclass=? and host=?")

	if err != nil {
		utils.Warning.Printf("%q, checkSC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(parameter, puppetClass, host).Scan(&id, &ovrCount, &foremanID)
	if err != nil {
		return SmartClass{}, err
	}

	return SmartClass{
		ID:                  id,
		ForemanID:           foremanID,
		Name:                parameter,
		OverrideValuesCount: ovrCount,
	}, nil
}

// Return the Smart Class parameters overrides by Smart Class ID and 'match'
func (Get) Override(scID int, name string, parameter string, ctx *user.GlobalCTX) (Override, error) {

	// VARS
	matchStr := fmt.Sprintf("hostgroup=SWE/%s", name)

	// ======
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id, `value` from override_values where sc_id=? and `match` = ?")
	if err != nil {
		utils.Warning.Printf("%q, getOvrData", err)
	}
	defer utils.DeferCloseStmt(stmt)

	var foremanId int
	var val string

	err = stmt.QueryRow(scID, matchStr).Scan(&foremanId, &val)
	if err != nil {
		return Override{}, err
	}

	return Override{
		Match: matchStr,
		Value: val,
	}, nil
}

// Return the Smart Class parameters overrides by match
func (Get) OverridesByMatch(host, matchParameter string, ctx *user.GlobalCTX) map[string][]Override {

	// VARS
	var gDB Get
	results := make(map[string][]Override)
	//qStr := fmt.Sprintf("location=%s", matchParameter)

	// ========
	stmt, err := ctx.Config.Database.DB.Prepare("select  ov.id, ov.`match`, ov.value, ov.sc_id, ov.foreman_id as ovr_foreman_id, sc.foreman_id  as sc_foreman_id, sc.parameter,sc.parameter_type, sc.puppetclass from override_values as ov, smart_classes as sc where ov.`match`= ? and sc.id = ov.sc_id and sc.host = ?")
	if err != nil {
		utils.Warning.Printf("%q, getOverridesLoc", err)
	}
	defer utils.DeferCloseStmt(stmt)
	rows, err := stmt.Query(matchParameter, host)
	if err != nil {
		utils.Warning.Printf("%q, getOverridesLoc", err)
	}

	for rows.Next() {
		var (
			ovrId        int
			ovrFId       int
			scFId        int
			smartClassId int
			param        string
			_type        string
			pc           string
			value        string
			match        string
		)

		err = rows.Scan(&ovrId, &match, &value, &smartClassId,
			&ovrFId, &scFId, &param, &_type, &pc)
		if err != nil {
			utils.Warning.Printf("%q, getOverridesLoc", err)
		}

		var dumpObj APISmartClass
		scData, _ := gDB.ByID(smartClassId, ctx)
		_ = json.Unmarshal([]byte(scData.Dump), &dumpObj)

		results[pc] = append(results[pc], Override{
			SmartClass: &scData,
			ID:         ovrId,
			ForemanID:  ovrFId,
			Value:      value,
		})
	}

	return results
}
