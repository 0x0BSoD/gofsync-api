package puppetclass

import (
	"fmt"
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
// STATEMENTS
// ======================================================
var (
	selectID        = "select id from puppet_classes where host_id=? and subclass=?"
	selectForemanID = "select foreman_id from puppet_classes where host_id=? and subclass=?"
	selectAll       = "select id, foreman_id, class, subclass, sc_ids from puppet_classes where host_id=?"
	selectByName    = "select id, class, sc_ids, env_ids, foreman_id from puppet_classes where subclass=? and host_id=?"
	selectByID      = "select class, subclass, sc_ids, env_ids from puppet_classes where id=?"

	insert = "insert into puppet_classes(host_id, class, subclass, foreman_id, sc_ids, env_ids) values(?,?,?,?,?,?)"
	update = "update puppet_classes set sc_ids=?, env_ids=? where host_id=? and foreman_id=?"

	updateHG = "update hg set pcList=? where id=?"

	deletePC = "delete from puppet_classes where (`host_id` = ? and `foreman_id`=?);"
)

// ======================================================
// CHECKS
// ======================================================
func ID(hostID int, subclass string, ctx *user.GlobalCTX) int {
	var id int
	stmt, err := ctx.Config.Database.DB.Prepare(selectID)
	if err != nil {
		logger.Warning.Printf("%q, checkPC", err)
	}
	defer utils.DeferCloseStmt(stmt)
	err = stmt.QueryRow(hostID, subclass).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

func ForemanID(hostID int, subclass string, ctx *user.GlobalCTX) int {
	var foremanID int
	stmt, err := ctx.Config.Database.DB.Prepare(selectForemanID)
	if err != nil {
		logger.Warning.Printf("%q, checkPC", err)
	}
	defer utils.DeferCloseStmt(stmt)
	err = stmt.QueryRow(hostID, subclass).Scan(&foremanID)
	if err != nil {
		return -1
	}
	return foremanID
}

// ======================================================
// GET
// ======================================================
func DbAll(hostID int, ctx *user.GlobalCTX) []PCintId {
	var res []PCintId
	stmt, err := ctx.Config.Database.DB.Prepare(selectAll)
	if err != nil {
		logger.Warning.Printf("%q, getByNamePC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query(hostID)
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

func DbByName(hostID int, subclass string, ctx *user.GlobalCTX) PC {
	var class string
	var sCIDs string
	var envIDs string
	var foremanId int
	var id int

	stmt, err := ctx.Config.Database.DB.Prepare(selectByName)
	if err != nil {
		logger.Warning.Printf("%q, getByNamePC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(subclass, hostID).Scan(&id, &class, &sCIDs, &envIDs, &foremanId)
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
func DbByID(pID int, ctx *user.GlobalCTX) PC {
	var class string
	var subclass string
	var sCIDs string
	var envIDs string

	stmt, err := ctx.Config.Database.DB.Prepare(selectByID)
	if err != nil {
		logger.Warning.Printf("%q, getPC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(pID).Scan(&class, &subclass, &sCIDs, &envIDs)
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
func DbInsert(hostID, foremanID int, class, subclass string, ctx *user.GlobalCTX) int {
	existID := ID(hostID, subclass, ctx)
	if existID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare(insert)
		if err != nil {
			logger.Warning.Printf("%q, insertPC", err)
		}
		defer utils.DeferCloseStmt(stmt)

		res, err := stmt.Exec(hostID, class, subclass, foremanID, "NULL", "NULL")
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
func DbUpdate(hostID int, puppetClass smartclass.PCSCParameters, ctx *user.GlobalCTX) {
	var strScList []string
	var strEnvList []string

	//sort.Slice(puppetClass.SmartClassParameters, func(i, j int) bool {
	//	return puppetClass.SmartClassParameters[i].ID < puppetClass.SmartClassParameters[j].ID
	//})
	//sort.Slice(puppetClass.Environments, func(i, j int) bool {
	//	return puppetClass.Environments[i].ID < puppetClass.Environments[j].ID
	//})

	for _, i := range puppetClass.SmartClassParameters {
		scID := smartclass.ScID(hostID, i.ID, ctx)

		//fmt.Printf("%d\t%s\t%s\t%s\n", scID, host, puppetClass.Name, i.Parameter)

		if scID != -1 {
			strScList = append(strScList, strconv.Itoa(scID))
		}
	}

	for _, i := range puppetClass.Environments {
		envID := environment.ID(hostID, i.Name, ctx)
		if envID != -1 {
			strEnvList = append(strEnvList, strconv.Itoa(envID))
		}
	}

	//fmt.Printf("update puppet_classes set sc_ids='%s', env_ids='%s' where host='%s' and foreman_id='%d'\n", strings.Join(strScList, ","),
	//	strings.Join(strEnvList, ","),
	//	host,
	//	puppetClass.ID)
	stmt, err := ctx.Config.Database.DB.Prepare(update)
	if err != nil {
		logger.Warning.Printf("%q, updatePC", err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(
		strings.Join(strScList, ","),
		strings.Join(strEnvList, ","),
		hostID,
		puppetClass.ID)
	if err != nil {
		logger.Warning.Printf("%q, updatePC", err)
	}

}

func DbUpdatePcID(hgID int, pcList []int, ctx *user.GlobalCTX) {

	var strPcList []string

	for _, i := range pcList {
		if i != 0 {
			strPcList = append(strPcList, utils.String(i))
		}
	}
	pcListStr := strings.Join(strPcList, ",")
	stmt, err := ctx.Config.Database.DB.Prepare(updateHG)
	if err != nil {
		logger.Error.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(pcListStr, hgID)
	if err != nil {
		logger.Error.Println(err)
	}

}

// ======================================================
// DELETE
// ======================================================
func DeletePuppetClass(hostID, foremanID int, ctx *user.GlobalCTX) {
	fmt.Println(hostID, foremanID)
	stmt, err := ctx.Config.Database.DB.Prepare(deletePC)
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Query(hostID, foremanID)
	if err != nil {
		logger.Warning.Printf("%q, DeletePuppetClass", err)
	}
}
