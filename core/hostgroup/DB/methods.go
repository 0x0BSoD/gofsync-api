package DB

import (
	"git.ringcentral.com/archops/goFsync/core/puppetclass/DB"
	DB2 "git.ringcentral.com/archops/goFsync/core/smartclass/DB"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Return DB ID for host group
func (Get) ID(name, host string, ctx *user.GlobalCTX) int {

	// VARS
	var ID int

	// ====
	stmt, err := ctx.Config.Database.DB.Prepare("select id from hg where name=? and host=?")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(name, host).Scan(&ID)
	if err != nil {
		return -1
	}

	return ID
}

// Return Foreman ID for host group
func (Get) ForemanID(name, host string, ctx *user.GlobalCTX) int {

	// VARS
	var ID int

	// ==========
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id from hg where name=? and host=?")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(name, host).Scan(&ID)
	if err != nil {
		ID = -1
	}

	return ID
}

// Return DB ID for host group parameter
func (Get) ParameterID(hgId int, name string, ctx *user.GlobalCTX) int {

	// VARS
	var ID int

	// ====
	stmt, err := ctx.Config.Database.DB.Prepare("select id from hg_parameters where hg_id=? and name=?")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)
	err = stmt.QueryRow(hgId, name).Scan(&ID)
	if err != nil {
		return -1
	}

	return ID
}

// Return Foreman ID for puppet master host
func (Get) ForemanIDs(host string, ctx *user.GlobalCTX) []int {

	// VARS
	var result []int

	// =======
	stmt, err := ctx.Config.Database.DB.Prepare("SELECT foreman_id FROM hg WHERE host=?;")
	if err != nil {
		utils.Warning.Printf("%q, GetForemanIDs", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
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

// Return all host groups
func (Get) List(ctx *user.GlobalCTX) []string {

	stmt, err := ctx.Config.Database.DB.Prepare("select id, name from hg")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var result []string

	rows, err := stmt.Query()
	if err != nil {
		return []string{}
	}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			utils.Error.Println(err)
		}
		if !utils.StringInSlice(name, result) {
			result = append(result, name)
		}

	}
	sort.Strings(result)

	return result
}

// Return all host groups for foreman host
func (Get) ListByHost(host string, ctx *user.GlobalCTX) []HostGroupJSON {

	stmt, err := ctx.Config.Database.DB.Prepare("select id, foreman_id, name, status, updated_at from hg where host=?")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var result []HostGroupJSON

	rows, err := stmt.Query(host)
	if err != nil {
		return []HostGroupJSON{}
	}
	for rows.Next() {

		var (
			ID        int
			foremanID int
			name      string
			status    string
			updated   string
		)

		err = rows.Scan(&ID, &foremanID, &name, &status, &updated)
		if err != nil {
			utils.Error.Println(err)
		}
		result = append(result, HostGroupJSON{
			ID:        ID,
			ForemanID: foremanID,
			Name:      name,
			Status:    status,
			Updated:   updated,
		})

	}

	return result
}

// Return Host Group data by name
func (Get) ByName(host, name string, ctx *user.GlobalCTX) (HostGroupJSON, error) {
	// VARS
	var pcGDB DB.Get
	var scGDB DB2.Get
	var (
		id        int
		foremanId int
		dump      string
		pcList    string
		status    string
		updated   string
	)
	var puppetClasses = make(map[string][]HostGroupPuppetCLass)

	// ====
	stmt, err := ctx.Config.Database.DB.Prepare("SELECT id, foreman_id, dump, pcList, status, updated_at FROM hg WHERE host=? AND name=?;")
	if err != nil {
		utils.Warning.Printf("%q, errror while getting host group", err)
	}

	err = stmt.QueryRow(host, name).Scan(&id, &foremanId, &dump, &pcList, &status, &updated)
	if err != nil {
		return HostGroupJSON{}, err
	}
	utils.DeferCloseStmt(stmt)

	// Puppet Classes ==
	pcIDs := utils.Integers(pcList)
	for _, pcID := range pcIDs {
		puppetClass := pcGDB.ByID(pcID, ctx)
		var smartClasses []DB2.SmartClass
		for _, scID := range puppetClass.SmartClassIDs {
			if scID != 0 {
				smartClass, err := scGDB.ByHostGroup(scID, name, ctx)
				if err != nil {
					utils.Warning.Println(err)
				}
				smartClasses = append(smartClasses, smartClass)
			} else {
				smartClasses = nil
				break
			}
		}
		puppetClasses[puppetClass.Class] = append(puppetClasses[puppetClass.Class], HostGroupPuppetCLass{
			ID:           puppetClass.ID,
			ForemanID:    puppetClass.ForemanID,
			Class:        puppetClass.Class,
			Subclass:     puppetClass.Subclass,
			SmartClasses: smartClasses,
		})
	}

	return HostGroupJSON{
		ID:            id,
		ForemanID:     foremanId,
		Name:          name,
		Status:        status,
		PuppetClasses: puppetClasses,
		Updated:       updated,
	}, nil
}

// Return Host Group data by ID
func (Get) ByID(ID int, ctx *user.GlobalCTX) (HostGroupJSON, error) {
	// VARS
	var pcGDB DB.Get
	var scGDB DB2.Get
	var (
		id        int
		foremanId int
		name      string
		dump      string
		pcList    string
		status    string
		updated   string
	)
	var puppetClasses = make(map[string][]HostGroupPuppetCLass)

	// ====
	stmt, err := ctx.Config.Database.DB.Prepare("SELECT id, name, foreman_id, dump, pcList, status, updated_at FROM hg WHERE id=?;")
	if err != nil {
		utils.Warning.Printf("%q, errror while getting host group", err)
	}

	err = stmt.QueryRow(ID).Scan(&id, &name, &foremanId, &dump, &pcList, &status, &updated)
	if err != nil {
		return HostGroupJSON{}, err
	}
	utils.DeferCloseStmt(stmt)

	// Puppet Classes ==
	pcIDs := utils.Integers(pcList)
	for _, pcID := range pcIDs {
		puppetClass := pcGDB.ByID(pcID, ctx)
		var smartClasses []DB2.SmartClass
		for _, scID := range puppetClass.SmartClassIDs {
			if scID != 0 {
				smartClass, err := scGDB.ByHostGroup(scID, name, ctx)
				if err != nil {
					utils.Warning.Println(err)
				}
				smartClasses = append(smartClasses, smartClass)
			} else {
				smartClasses = nil
				break
			}
		}
		puppetClasses[puppetClass.Class] = append(puppetClasses[puppetClass.Class], HostGroupPuppetCLass{
			ID:           puppetClass.ID,
			ForemanID:    puppetClass.ForemanID,
			Class:        puppetClass.Class,
			Subclass:     puppetClass.Subclass,
			SmartClasses: smartClasses,
		})
	}

	return HostGroupJSON{
		ID:            id,
		ForemanID:     foremanId,
		Name:          name,
		Status:        status,
		PuppetClasses: puppetClasses,
		Updated:       updated,
	}, nil
}
