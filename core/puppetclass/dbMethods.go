package puppetclass

import (
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sort"
	"strconv"
	"strings"
)

// ======================================================
// CHECKS
// ======================================================
func DbID(subclass string, host string, ctx *user.GlobalCTX) int {

	var id int

	stmt, err := ctx.Config.Database.DB.Prepare("select id from puppet_classes where host=? and subclass=?")
	if err != nil {
		logger.Warning.Printf("%q, checkPC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, subclass).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

func ForemanID(subclass string, host string, ctx *user.GlobalCTX) int {

	var foremanID int

	stmt, err := ctx.Config.Database.DB.Prepare("select foreman_id from puppet_classes where host=? and subclass=?")
	if err != nil {
		logger.Warning.Printf("%q, checkPC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host, subclass).Scan(&foremanID)
	if err != nil {
		return -1
	}
	return foremanID
}

// ======================================================
// GET
// ======================================================
func DbAll(host string, ctx *user.GlobalCTX) []PCintId {

	var res []PCintId
	stmt, err := ctx.Config.Database.DB.Prepare("SELECT id, foreman_id, class, subclass, sc_ids from goFsync.puppet_classes where host=?;")
	if err != nil {
		logger.Warning.Printf("%q, getByNamePC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(host)
	if err != nil {
		return []PCintId{}
	}

	for rows.Next() {
		var foremanId int
		var class string
		var subclass string
		var scIds string
		var _id int
		err := rows.Scan(&_id, &foremanId, &class, &subclass, &scIds)
		if err != nil {
			logger.Warning.Println("No result while getting puppet classes")
		}

		if scIds != "" {
			intScIds := logger.Integers(scIds)
			res = append(res, PCintId{
				ID:        _id,
				ForemanId: foremanId,
				Class:     class,
				Subclass:  subclass,
				SCIDs:     intScIds,
			})
		} else {
			res = append(res, PCintId{
				ForemanId: foremanId,
				Class:     class,
				Subclass:  subclass,
			})
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ForemanId < res[j].ForemanId
	})

	return res
}

func DbByName(subclass string, host string, ctx *user.GlobalCTX) PC {

	var class string
	var sCIDs string
	var envIDs string
	var foremanId int
	var id int

	stmt, err := ctx.Config.Database.DB.Prepare("select id, class, sc_ids, env_ids, foreman_id from puppet_classes where subclass=? and host=?")
	if err != nil {
		logger.Warning.Printf("%q, getByNamePC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(subclass, host).Scan(&id, &class, &sCIDs, &envIDs, &foremanId)
	if err != nil {
		return PC{}
	}

	return PC{
		ID:        id,
		ForemanId: foremanId,
		Class:     class,
		Subclass:  subclass,
		SCIDs:     sCIDs,
	}
}
func DbByID(pId int, ctx *user.GlobalCTX) PC {

	var class string
	var subclass string
	var sCIDs string
	var envIDs string

	stmt, err := ctx.Config.Database.DB.Prepare("select class, subclass, sc_ids, env_ids from puppet_classes where id=?")
	if err != nil {
		logger.Warning.Printf("%q, getPC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(pId).Scan(&class, &subclass, &sCIDs, &envIDs)
	if err != nil {
		logger.Warning.Printf("%q, getPC", err)
	}

	return PC{
		Class:    class,
		Subclass: subclass,
		SCIDs:    sCIDs,
	}
}

// ======================================================
// INSERT
// ======================================================
func DbInsert(host string, class string, subclass string, foremanId int, ctx *user.GlobalCTX) int {

	existID := DbID(subclass, host, ctx)
	if existID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into puppet_classes(host, class, subclass, foreman_id, sc_ids, env_ids) values(?,?,?,?,?,?)")
		if err != nil {
			logger.Warning.Printf("%q, insertPC", err)
		}
		defer utils.DeferCloseStmt(stmt)

		res, err := stmt.Exec(host, class, subclass, foremanId, "NULL", "NULL")
		if err != nil {
			logger.Warning.Printf("%q, checkPC", err)
		}

		lastID, _ := res.LastInsertId()
		return int(lastID)
	} else {
		return existID
	}
}

// ======================================================
// UPDATE
// ======================================================
func DbUpdate(host string, puppetClass smartclass.PCSCParameters, ctx *user.GlobalCTX) {
	var strScList []string
	var strEnvList []string

	//sort.Slice(puppetClass.SmartClassParameters, func(i, j int) bool {
	//	return puppetClass.SmartClassParameters[i].ID < puppetClass.SmartClassParameters[j].ID
	//})
	//sort.Slice(puppetClass.Environments, func(i, j int) bool {
	//	return puppetClass.Environments[i].ID < puppetClass.Environments[j].ID
	//})

	for _, i := range puppetClass.SmartClassParameters {
		scID := smartclass.CheckSCByForemanId(host, i.ID, ctx)

		//fmt.Printf("%d\t%s\t%s\t%s\n", scID, host, puppetClass.Name, i.Parameter)

		if scID != -1 {
			strScList = append(strScList, strconv.Itoa(int(scID)))
		}
	}

	for _, i := range puppetClass.Environments {
		envID := environment.ID(host, i.Name, ctx)
		if envID != -1 {
			strEnvList = append(strEnvList, strconv.Itoa(int(envID)))
		}
	}

	//fmt.Printf("update puppet_classes set sc_ids='%s', env_ids='%s' where host='%s' and foreman_id='%d'\n", strings.Join(strScList, ","),
	//	strings.Join(strEnvList, ","),
	//	host,
	//	puppetClass.ID)
	stmt, err := ctx.Config.Database.DB.Prepare("update puppet_classes set sc_ids=?, env_ids=? where host=? and foreman_id=?")
	if err != nil {
		logger.Warning.Printf("%q, updatePC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(
		strings.Join(strScList, ","),
		strings.Join(strEnvList, ","),
		host,
		puppetClass.ID)
	if err != nil {
		logger.Warning.Printf("%q, updatePC", err)
	}

}

func DbUpdatePcID(hgId int, pcList []int, ctx *user.GlobalCTX) {

	var strPcList []string

	for _, i := range pcList {
		if i != 0 {
			strPcList = append(strPcList, utils.String(i))
		}
	}
	pcListStr := strings.Join(strPcList, ",")
	stmt, err := ctx.Config.Database.DB.Prepare("update hg set pcList=? where id=?")
	if err != nil {
		logger.Error.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(pcListStr, hgId)
	if err != nil {
		logger.Error.Println(err)
	}

}

// ======================================================
// DELETE
// ======================================================
func DeletePuppetClass(host string, foremanID int, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM puppet_classes WHERE (`host` = ? and `foreman_id`=?);")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(host, foremanID)
	if err != nil {
		logger.Warning.Printf("%q, DeletePuppetClass", err)
	}
}
