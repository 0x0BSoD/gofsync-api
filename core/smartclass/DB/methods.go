package DB

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/smartclass/API"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
	"strconv"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Get Smart Class ID from DB
func (Get) ID(host string, pc string, parameter string, ctx *user.GlobalCTX) int {

	// VARS
	var ID int

	// =======
	stmt, err := ctx.Config.Database.DB.Prepare("select id from smart_classes where host=? and parameter=? and puppetclass=?")
	if err != nil {
		utils.Warning.Printf("%q, error while getting smart calss ID", err)
		ID = -1
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, parameter, pc).Scan(&ID)
	if err != nil {
		ID = -1
	}
	return ID
}

// Return ID for override by Smart Class ID and Foreman ID
func (Get) OverrideID(scID, foremanID int, ctx *user.GlobalCTX) int {

	// VARS
	var ID int

	// =========
	stmt, err := ctx.Config.Database.DB.Prepare("SELECT id FROM override_values where sc_id=? and foreman_id=?")
	if err != nil {
		utils.Warning.Printf("%q, error while getting override ID", err)
		ID = -1
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(scID, foremanID).Scan(&ID)
	if err != nil {
		utils.Warning.Printf("%q, error while getting override ID", err)
		ID = -1
	}

	return ID
}

// Get Overrides Foreman IDs by Smart Class ID from DB
func (Get) OverridesFIDsBySmartClassID(scId int, ctx *user.GlobalCTX) []int {

	// VARS
	var IDs []int

	// =========
	stmt, err := ctx.Config.Database.DB.Prepare("SELECT foreman_id FROM override_values WHERE sc_id=?;")
	if err != nil {
		utils.Warning.Printf("%q, error while getting overrides IDs", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(scId)
	if err != nil {
		utils.Warning.Printf("%q, error while getting overrides IDs", err)
	}

	for rows.Next() {
		var ID int
		err = rows.Scan(&ID)
		if err != nil {
			utils.Warning.Printf("%q, error while getting overrides IDs", err)
		}
		IDs = append(IDs, ID)
	}

	return IDs
}

// Get all Smart Classes Foreman IDs from DB
func (Get) ForemanIDs(host string, ctx *user.GlobalCTX) []int {

	// VARS
	var IDs []int

	// =================
	stmt, err := ctx.Config.Database.DB.Prepare("SELECT foreman_id FROM smart_classes WHERE host=?;")
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Classes Foreman IDs", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Classes Foreman IDs", err)
	}

	for rows.Next() {
		var ID int
		err = rows.Scan(&ID)
		if err != nil {
			utils.Warning.Printf("%q, error while getting Smart Classes Foreman IDs", err)
		}
		IDs = append(IDs, ID)
	}

	return IDs
}

// Get Smart Class ID by Foreman ID and host name from DB
func (Get) IDByForemanID(host string, foremanID int, ctx *user.GlobalCTX) int {

	// VARS
	var id int

	// ======
	stmt, err := ctx.Config.Database.DB.Prepare("select id from smart_classes where host=? and foreman_id=?")
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Class Foreman ID", err)
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

	// VARS
	var (
		id        int
		foremanId int
		paramName string
		ovrCount  int
		_type     string
		pc        string
		dump      string
		gDB       Get
		dumpObj   SmartClass
	)

	// =======
	stmt, err := ctx.Config.Database.DB.Prepare("select id, parameter, override_values_count, foreman_id, parameter_type, puppetclass, dump from smart_classes where id=?")
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Class", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(scID).Scan(&id, &paramName, &ovrCount, &foremanId, &_type, &pc, &dump)
	if err != nil {
		return SmartClass{}, err
	}

	overrides, _ := gDB.Overrides(scID, ctx)

	_ = json.Unmarshal([]byte(dump), &dumpObj)

	return SmartClass{
		ID:                  id,
		ForemanID:           foremanId,
		Name:                paramName,
		OverrideValuesCount: ovrCount,
		ValueType:           _type,
		Override:            overrides,
		DefaultVal:          dumpObj.DefaultVal,
		PuppetClass:         pc,
		Dump:                dump,
	}, nil
}

// Get Smart Class by ID
func (Get) ByHostGroup(host, hostGroup string, ctx *user.GlobalCTX) ([]SmartClass, error) {

	// VARS
	var (
		gDB     Get
		dumpObj SmartClass
		result  []SmartClass
	)

	// =======
	stmt, err := ctx.Config.Database.DB.Prepare("select id, parameter, override_values_count, foreman_id, parameter_type, puppetclass, dump from smart_classes where host=?")
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Class", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		return []SmartClass{}, err
	}

	for rows.Next() {
		var (
			id        int
			foremanId int
			paramName string
			ovrCount  int
			_type     string
			pc        string
			dump      string
		)

		err := rows.Scan(&id, &paramName, &ovrCount, &foremanId, &_type, &pc, &dump)
		if err != nil {
			utils.Warning.Printf("%q, error while getting Smart Class Overrides for host group", err)
		}

		match := fmt.Sprintf("hostgroup=SWE/%s", hostGroup)
		var override Override
		ovrResult := []Override{override}
		if ovrCount > 0 {
			ovrResult[0], err = gDB.OverrideByMatch(id, match, ctx)
			if err != nil {
				utils.Warning.Printf("%q, error while getting Smart Class Overrides for host group", err)
			}
		} else {
			ovrResult = nil
		}

		_ = json.Unmarshal([]byte(dump), &dumpObj)
		result = append(result, SmartClass{
			ID:                  id,
			ForemanID:           foremanId,
			Name:                paramName,
			OverrideValuesCount: ovrCount,
			ValueType:           _type,
			Override:            ovrResult,
			DefaultVal:          dumpObj.DefaultVal,
			PuppetClass:         pc,
			Dump:                dump,
		})
	}

	return result, nil
}

// Get Smart Class by parameter and puppet class name
func (Get) ByParameter(host string, puppetClass string, parameter string, ctx *user.GlobalCTX) (SmartClass, error) {

	var (
		id        int
		foremanID int
		ovrCount  int
	)

	stmt, err := ctx.Config.Database.DB.Prepare("select id, override_values_count, foreman_id from smart_classes where parameter=? and puppetclass=? and host=?")

	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Class", err)
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

// Return the Smart Class parameters overrides by Smart Class ID
func (Get) OverrideByMatch(scID int, matchParameter string, ctx *user.GlobalCTX) (Override, error) {

	// VARS
	//matchStr := fmt.Sprintf("hostgroup=SWE/%s", matchParameter)

	// ======
	fmt.Printf("select foreman_id, `value` from override_values where sc_id=%d and `match` = %s\n", scID, matchParameter)
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id, `value` from override_values where sc_id=? and `match` = ?")
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Class Overrides", err)
	}
	defer utils.DeferCloseStmt(stmt)

	var foremanId int
	var val string

	err = stmt.QueryRow(scID, matchParameter).Scan(&foremanId, &val)
	if err != nil {
		return Override{}, err
	}

	return Override{
		ForemanID: foremanId,
		Match:     matchParameter,
		Value:     val,
	}, nil
}

// Return the Smart Class parameters overrides by Smart Class ID
func (Get) Overrides(scID int, ctx *user.GlobalCTX) ([]Override, error) {

	// VARS
	var result []Override

	// ======
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id, `value`, `match` from override_values where sc_id=?")
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Class Overrides", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(scID)
	if err != nil {
		return []Override{}, err
	}

	for rows.Next() {
		var foremanId int
		var match string
		var val string

		err = rows.Scan(&foremanId, &val, &match)
		if err != nil {
			utils.Warning.Printf("%q, error while getting Smart Class Overrides", err)
		}
		result = append(result, Override{
			ForemanID: foremanId,
			Match:     match,
			Value:     val,
		})
	}

	return result, nil
}

// Return the Smart Class parameters overrides by match
func (Get) OverridesByMatch(host, matchParameter string, ctx *user.GlobalCTX) map[string][]Override {

	// VARS
	var gDB Get
	results := make(map[string][]Override)

	// ========
	stmt, err := ctx.Config.Database.DB.Prepare("select  ov.id, ov.`match`, ov.value, ov.sc_id, ov.foreman_id as ovr_foreman_id, sc.foreman_id  as sc_foreman_id, sc.parameter,sc.parameter_type, sc.puppetclass from override_values as ov, smart_classes as sc where ov.`match`= ? and sc.id = ov.sc_id and sc.host = ?")
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Class Overrides for Location", err)
	}
	defer utils.DeferCloseStmt(stmt)
	rows, err := stmt.Query(matchParameter, host)
	if err != nil {
		utils.Warning.Printf("%q, error while getting Smart Class Overrides for Location", err)
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
			utils.Warning.Printf("%q, error while getting Smart Class Overrides for Location", err)
		}

		var dumpObj API.Parameter
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

// =====================================================================================================================
// INSERT
// =====================================================================================================================

// Insert new smart class to DB
func (Insert) Add(host string, SmartClass API.Parameter, ctx *user.GlobalCTX) {

	// VARS
	var gDB Get
	var iDB Insert
	var dDB Delete

	// ==== SMART CLASS =====
	ID := gDB.IDByForemanID(host, SmartClass.ForemanID, ctx)
	if ID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into smart_classes(host, puppetclass, parameter, parameter_type, foreman_id, override_values_count, dump) values(?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			utils.Warning.Printf("%q, error while adding new Smart Class", err)
		}
		defer utils.DeferCloseStmt(stmt)

		sJson, _ := json.Marshal(SmartClass)
		res, err := stmt.Exec(host, SmartClass.PuppetClass.Name, SmartClass.Parameter, SmartClass.ParameterType, SmartClass.ForemanID, SmartClass.OverrideValuesCount, sJson)
		if err != nil {
			utils.Warning.Printf("%q, error while adding new Smart Class", err)
		}

		lastId, _ := res.LastInsertId()
		ID = int(lastId)
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE smart_classes SET `override_values_count` = ? WHERE (`id` = ?)")
		if err != nil {
			utils.Warning.Printf("%q, error while updating Smart Class", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(SmartClass.OverrideValuesCount, ID)
		if err != nil {
			utils.Warning.Printf("%q, error while updating Smart Class", err)
		}
	}

	// ==== OVERRIDES =====
	if SmartClass.OverrideValuesCount > 0 {

		beforeUpdateIDs := gDB.OverridesFIDsBySmartClassID(ID, ctx)
		sort.Ints(beforeUpdateIDs)
		var afterUpdateIDs []int

		for _, ovr := range SmartClass.OverrideValues {
			afterUpdateIDs = append(afterUpdateIDs, ovr.ForemanID)
			iDB.AddOverride(ID, ovr, SmartClass.ParameterType, ctx)
		}

		sort.Ints(afterUpdateIDs)
		for _, ForemanID := range beforeUpdateIDs {
			if !utils.Search(afterUpdateIDs, ForemanID) {
				err := dDB.Override(ID, ForemanID, ctx)
				if err != nil {
					utils.Warning.Printf("%q, error while deleting Smart Class Override", err)
				}
			}
		}

	}

}

// Insert Smart Class override
func (Insert) AddOverride(scId int, override API.OverrideValue, pType string, ctx *user.GlobalCTX) {

	// VARS
	var gDB Get
	var strData string

	// Value assertion
	// =================================================================================================================
	if override.Value != nil {
		switch override.Value.(type) {
		case string:
			strData = override.Value.(string)
		case []interface{}:
			var tmpResInt []string
			for _, i := range override.Value.([]interface{}) {
				tmpResInt = append(tmpResInt, i.(string))
			}
			strIng, _ := json.Marshal(tmpResInt)
			strData = string(strIng)
		case bool:
			strData = strconv.FormatBool(override.Value.(bool))
		case int:
			strData = strconv.FormatFloat(override.Value.(float64), 'f', 6, 64)
		case float64:
			strData = strconv.FormatFloat(override.Value.(float64), 'f', 6, 64)
		default:
			utils.Warning.Printf("Type not known try save as string, Type: %s, Val: %s, Match: %s", pType, override.Value, override.Match)
			strData = override.Value.(string)
		}
	}
	// =================================================================================================================

	ID := gDB.OverrideID(scId, override.ForemanID, ctx)
	if ID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into override_values(foreman_id, `match`, value, sc_id, use_puppet_default) values(?,?,?,?,?)")
		if err != nil {
			utils.Warning.Printf("%q, insertSCOverride", err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(override.ForemanID, override.Match, strData, scId, override.UsePuppetDefault)
		if err != nil {
			utils.Warning.Printf("%q, insertSCOverride", err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE override_values SET `value` = ?, foreman_id=? WHERE id= ?")
		if err != nil {
			utils.Warning.Printf("%q, Prepare updateSCOverride override: %q, %d", err, strData, ID)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(strData, override.ForemanID, ID)
		if err != nil {
			utils.Warning.Printf("%q, updateSCOverride override: %q, %d", err, strData, ID)
		}
	}
}

// =====================================================================================================================
// DELETE
// =====================================================================================================================

func (Delete) SmartClass(host string, foremanID int, ctx *user.GlobalCTX) error {

	// =======
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM smart_classes WHERE host=? and foreman_id=?")
	if err != nil {
		utils.Warning.Println(err)
		return err
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, foremanID)
	if err != nil {
		utils.Warning.Printf("%q, DeleteSmartClass	", err)
		return err
	}

	return nil
}

func (Delete) Override(scID int, foremanID int, ctx *user.GlobalCTX) error {

	// =======
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM override_values WHERE sc_id=? AND foreman_id=?")
	if err != nil {
		utils.Warning.Println(err)
		return err
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(scID, foremanID)
	if err != nil {
		utils.Warning.Printf("%q, DeleteSmartClass	", err)
		return err
	}

	return nil
}
