package foremans

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"strconv"
	"strings"
)

// ======================================================
// STATEMENTS
// ======================================================
var (
	selectID    = "select id from hosts where name=?"
	selectAll   = "select id, name, env from hosts"
	selectStats = "select trend, Success, Failed, RFailed, Total, Last from hosts where id = ?"

	insert = "insert into hosts (name) values(?)"

	updateTrends = "update hosts set `trend` = ?, `Last` = ?, `Success` = ?, `Failed` = ?, `RFailed` = ?, `Total` = ? where (`id` = ?)"
)

// ======================================================
// GET
// ======================================================

// Return DB ID for puppet master host parameter
func ForemanHostID(host string, cfg *models.Config) int {
	stmt, err := cfg.Database.DB.Prepare(selectID)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)
	var id int
	err = stmt.QueryRow(host).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// Return all puppet master hosts with environments
func PuppetHosts(ctx *user.GlobalCTX) []ForemanHost {
	var result []ForemanHost
	stmt, err := ctx.Config.Database.DB.Prepare(selectAll)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	rows, err := stmt.Query()
	if err != nil {
		utils.Warning.Println(err)
	}
	for rows.Next() {
		var ID int
		var name string
		var env string
		err = rows.Scan(&ID, &name, &env)
		if err != nil {
			utils.Error.Println(err)
		}

		if _, ok := ctx.Config.Hosts[name]; ok {
			result = append(result, ForemanHost{
				ID:        ID,
				Name:      name,
				Env:       env,
				Dashboard: getTrends(ID, ctx),
			})
		}

	}
	return result
}

func getTrends(hostID int, ctx *user.GlobalCTX) Dashboard {
	var trend string
	var s int
	var f int
	var rf int
	var t int
	var last string
	var dash Dashboard

	stmt, err := ctx.Config.Database.DB.Prepare(selectStats)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(hostID).Scan(&trend, &s, &f, &rf, &t, &last)
	if err != nil {
		utils.Warning.Printf("%q, getDashboardData", err)
	}

	dash.LastHost = last
	dash.Summary = t
	dash.Success = s
	dash.RFailed = rf
	dash.Failed = f

	var trendStruct Trend
	var trendStr []string
	_ = json.Unmarshal([]byte(trend), &trendStr)

	for _, i := range trendStr {
		splt := strings.Split(i, ":")
		if len(splt) >= 2 {
			l, _ := strconv.Atoi(splt[0])
			v, _ := strconv.Atoi(splt[1])
			trendStruct.Labels = append(trendStruct.Labels, l)
			trendStruct.Values = append(trendStruct.Values, v)
		} else {
			trendStruct.Labels = append(trendStruct.Labels, 0)
			trendStruct.Values = append(trendStruct.Values, 0)
		}

	}
	//fmt.Println(trendStruct)
	dash.Trend = trendStruct

	return dash
}

// ======================================================
// INSERT
// ======================================================

// Insert puppet master host
func InsertHost(name string, cfg *models.Config) int {
	ID := ForemanHostID(name, cfg)
	if ID == -1 {
		stmt, err := cfg.Database.DB.Prepare(insert)
		if err != nil {
			utils.Warning.Println(err)
		}
		defer utils.DeferCloseStmt(stmt)
		resp, err := stmt.Exec(name)
		if err != nil {
			utils.Warning.Println(err)
		}
		id64, err := resp.LastInsertId()

		return int(id64)
	} else {
		return ID
	}
}

// ======================================================
// UPDATE
// ======================================================

func UpdateTrends(hostID int, data Dashboard, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare(updateTrends)
	if err != nil {
		utils.Warning.Println(err)
	}

	var tmp = make([]string, len(data.Trend.Labels))
	for idx, l := range data.Trend.Labels {
		tmp = append(tmp, fmt.Sprintf("%d:%d", l, data.Trend.Values[idx]))
	}
	jsonStr, err := json.Marshal(tmp)
	if err != nil {
		utils.Error.Println(err)
	}

	_, err = stmt.Exec(jsonStr, data.LastHost, data.Success, data.Failed, data.RFailed, data.Summary, hostID)
	if err != nil {
		utils.Warning.Println(err)
	}
	utils.DeferCloseStmt(stmt)
}
