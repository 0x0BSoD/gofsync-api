package DB

import (
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
		ForemanId:           foremanId,
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
		ForemanId:           foremanID,
		Name:                parameter,
		OverrideValuesCount: ovrCount,
	}, nil
}
