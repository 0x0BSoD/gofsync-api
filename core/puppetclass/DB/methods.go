package DB

import (
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"strings"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

// Get ID by subclass and host return -1 if not found
func (Get) ID(subclass string, host string, ctx *user.GlobalCTX) int {

	// VARS
	var id int

	// =====
	stmt, err := ctx.Config.Database.DB.Prepare("select id from puppet_classes where host=? and subclass=?")
	if err != nil {
		utils.Warning.Printf("%q, getting puppet class id error", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, subclass).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// Return all puppet classes by host
func (Get) All(host string, ctx *user.GlobalCTX) []PuppetClass {

	// VARS
	var res []PuppetClass

	// =====
	stmt, err := ctx.Config.Database.DB.Prepare("SELECT id, foreman_id, class, subclass, sc_ids from goFsync.puppet_classes where host=?;")
	if err != nil {
		utils.Warning.Printf("%q, getByNamePC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		return []PuppetClass{}
	}

	for rows.Next() {

		var (
			foremanId int
			class     string
			subclass  string
			scIds     string
			_id       int
		)

		err := rows.Scan(&_id, &foremanId, &class, &subclass, &scIds)
		if err != nil {
			utils.Warning.Println("No result while getting puppet classes")
		}

		if scIds != "" {
			intScIds := utils.Integers(scIds)
			res = append(res, PuppetClass{
				ID:            _id,
				ForemanID:     foremanId,
				Class:         class,
				Subclass:      subclass,
				SmartClassIDs: intScIds,
			})
		} else {
			res = append(res, PuppetClass{
				ForemanID: foremanId,
				Class:     class,
				Subclass:  subclass,
			})
		}
	}

	return res
}

// Return Puppet class by name
func (Get) ByName(subclass string, host string, ctx *user.GlobalCTX) PuppetClass {

	// VARS
	var (
		class     string
		scIds     string
		envIDs    string
		foremanId int
		id        int
	)

	// =====
	stmt, err := ctx.Config.Database.DB.Prepare("select id, class, sc_ids, env_ids, foreman_id from puppet_classes where subclass=? and host=?")
	if err != nil {
		utils.Warning.Printf("%q, getByNamePC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(subclass, host).Scan(&id, &class, &scIds, &envIDs, &foremanId)
	if err != nil {
		return PuppetClass{}
	}

	intScIds := utils.Integers(scIds)

	return PuppetClass{
		ID:            id,
		ForemanID:     foremanId,
		Class:         class,
		Subclass:      subclass,
		SmartClassIDs: intScIds,
	}
}

// Return Puppet class by ID
func (Get) ByID(pId int, ctx *user.GlobalCTX) PuppetClass {

	// VARS
	var (
		class     string
		subclass  string
		scIds     string
		envIDs    string
		foremanID int
	)

	// =====
	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id, class, subclass, sc_ids, env_ids from puppet_classes where id=?")
	if err != nil {
		utils.Warning.Printf("%q, getPC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(pId).Scan(&foremanID, &class, &subclass, &scIds, &envIDs)

	intScIds := utils.Integers(scIds)
	intEnvIds := utils.Integers(envIDs)

	return PuppetClass{
		ID:             pId,
		ForemanID:      foremanID,
		Class:          class,
		Subclass:       subclass,
		SmartClassIDs:  intScIds,
		EnvironmentIDs: intEnvIds,
	}
}

// =====================================================================================================================
// INSERT
// =====================================================================================================================

// Inserting new puppet class to base, return ID
func (Insert) Insert(host string, class string, subclass string, foremanId int, ctx *user.GlobalCTX) int {

	// VARS
	var PcDb Get

	// =========
	ID := PcDb.ID(subclass, host, ctx)
	if ID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into puppet_classes(host, class, subclass, foreman_id, sc_ids, env_ids) values(?,?,?,?,?,?)")
		if err != nil {
			utils.Warning.Printf("%q, insertPC", err)
		}
		defer utils.DeferCloseStmt(stmt)

		res, err := stmt.Exec(host, class, subclass, foremanId, "NULL", "NULL")
		if err != nil {
			utils.Warning.Printf("%q, error while inserting new puppet class", err)
		}

		lastID, _ := res.LastInsertId()
		return int(lastID)
	} else {
		return ID
	}
}

// =====================================================================================================================
// UPDATE
// =====================================================================================================================

// Update puppet class ids in host group ny id
func (Update) HostGroupIDs(hgId int, pcList []int, ctx *user.GlobalCTX) {

	// VARS
	var strPcList []string

	// =======
	for _, i := range pcList {
		if i != 0 {
			strPcList = append(strPcList, utils.String(i))
		}
	}
	pcListStr := strings.Join(strPcList, ",")

	stmt, err := ctx.Config.Database.DB.Prepare("update hg set pcList=? where id=?")
	if err != nil {
		utils.Error.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(pcListStr, hgId)
	if err != nil {
		utils.Error.Println(err)
	}

}

// =====================================================================================================================
// DELETE
// =====================================================================================================================

// Remove from base by subclass
func (Delete) BySubclass(host string, subClass string, ctx *user.GlobalCTX) error {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM puppet_classes WHERE `host` = ? and `subclass`=?;")
	if err != nil {
		utils.Warning.Println(err)
		return err
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, subClass)
	if err != nil {
		utils.Warning.Printf("%q, DeletePuppetClass", err)
		return err
	}
	return nil
}
