package DB

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/hosts"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"time"
)

// ======================================================
// IDS
// ======================================================

// Return DB ID for puppet master host parameter
func HID(host string, cfg *models.Config) int {
	stmt, err := cfg.Database.DB.Prepare("select id from hosts where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)
	var id int
	err = stmt.QueryRow(host).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================

// Return all puppet master hosts with environments
func PuppetHosts(ctx *user.GlobalCTX) []hosts.ForemanHost {
	var result []hosts.ForemanHost
	stmt, err := ctx.Config.Database.DB.Prepare("select host, env from hosts")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query()
	if err != nil {
		logger.Warning.Println(err)
	}
	for rows.Next() {
		var name string
		var env string
		err = rows.Scan(&name, &env)
		if err != nil {
			logger.Error.Println(err)
		}
		if logger.StringInSlice(name, ctx.Config.Hosts) {
			result = append(result, hosts.ForemanHost{
				Name: name,
				Env:  env,
			})
		}
	}
	return result
}

// Return Environment for puppet master host
func PuppetHostEnv(host string, ctx *user.GlobalCTX) string {
	stmt, err := ctx.Config.Database.DB.Prepare("select env from hosts where host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var hostEnv string
	err = stmt.QueryRow(host).Scan(&hostEnv)
	if err != nil {
		logger.Warning.Println(err)
		return ""
	}

	return hostEnv
}

// Get host group by DB ID

// ======================================================
// INSERT
// ======================================================

// Insert/Update host group
func Insert(name, host, data, sweStatus string, foremanId int, ctx *user.GlobalCTX) int {
	hgID := ID(name, host, ctx)
	if hgID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into hg(name, host, dump, created_at, updated_at, foreman_id, pcList, status) values(?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		res, err := stmt.Exec(name, host, data, time.Now(), time.Now(), foremanId, "NULL", sweStatus)
		if err != nil {
			return -1
		}

		lastID, _ := res.LastInsertId()
		return int(lastID)
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE hg SET  `status` = ?, `foreman_id` = ?, `updated_at` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(sweStatus, foremanId, time.Now(), hgID)
		if err != nil {
			return -1
		}

		return hgID
	}
}

// Insert/Update host group parameters
func InsertParameters(sweId int, p HostGroupP, ctx *user.GlobalCTX) {
	PID := PID(sweId, p.Name, ctx)
	if PID == -1 {
		stmt, err := ctx.Config.Database.DB.Prepare("insert into hg_parameters(hg_id, foreman_id, name, `value`, priority) values(?, ?, ?, ?, ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(sweId, p.ID, p.Name, p.Value, p.Priority)
		if err != nil {
			logger.Warning.Println(err)
		}
	} else {
		stmt, err := ctx.Config.Database.DB.Prepare("UPDATE hg_parameters SET `foreman_id` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)

		_, err = stmt.Exec(p.ID, PID)
		if err != nil {
			logger.Warning.Println(err)
		}
	}
}

// Insert puppet master host
func InsertHost(host string, cfg *models.Config) {
	if id := HID(host, cfg); id == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into hosts (host) values(?)")
		if err != nil {
			logger.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)
		_, err = stmt.Exec(host)
		if err != nil {
			logger.Warning.Println(err)
		}
	}
}

//func insertState(hgName, host, state string, ctx *user.GlobalCTX) {
//	ID := StateID(hgName, ctx)
//	if ID == -1 {
//		q := fmt.Sprintf("insert into hg_state (host_group, `%s`) values(?, ?)", host)
//		stmt, err := ctx.Config.Database.DB.Prepare(q)
//		if err != nil {
//			logger.Warning.Println(err)
//		}
//		defer utils.DeferCloseStmt(stmt)
//		_, err = stmt.Exec(hgName, state)
//		if err != nil {
//			logger.Warning.Println(err)
//		}
//	} else {
//		q := fmt.Sprintf("UPDATE `goFsync`.`hg_state` SET `%s` = ? WHERE (`id` = ?)", host)
//		stmt, err := ctx.Config.Database.DB.Prepare(q)
//		if err != nil {
//			logger.Warning.Println(err)
//		}
//		defer utils.DeferCloseStmt(stmt)
//		_, err = stmt.Exec(state, ID)
//		if err != nil {
//			logger.Warning.Println(err)
//		}
//	}
//
//}

// ======================================================
// UPDATE
// ======================================================

// ======================================================
// DELETE
// ======================================================
//func DeleteHGbyId(hgId int, ctx *user.GlobalCTX) {
//	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM hg WHERE id=?")
//	if err != nil {
//		logger.Warning.Println(err)
//	}
//	defer utils.DeferCloseStmt(stmt)
//
//	_, err = stmt.Exec(hgId)
//	if err != nil {
//		logger.Warning.Println(err)
//	}
//}

func Delete(foremanID int, host string, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("DELETE FROM hg WHERE foreman_id=? AND host=?")
	if err != nil {
		logger.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	_, err = stmt.Exec(foremanID, host)
	if err != nil {
		logger.Warning.Println(err)
	}
}
