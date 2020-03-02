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
	selectID      = "select id from hosts where name=?"
	selectAll     = "select id, name, env from hosts"
	selectHostEnv = "select env from hosts where id=?"

	insert = "insert into hosts (name) values(?)"
)

// ======================================================
// GET
// ======================================================

// Return DB ID for puppet master host parameter
func ForemanHostID(name string, cfg *models.Config) int {
	stmt, err := cfg.Database.DB.Prepare(selectID)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)
	var id int
	err = stmt.QueryRow(name).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// Return Environment for puppet master host
func PuppetHostEnv(hostID int, ctx *user.GlobalCTX) string {
	stmt, err := ctx.Config.Database.DB.Prepare(selectHostEnv)
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	var hostEnv string
	err = stmt.QueryRow(hostID).Scan(&hostEnv)
	if err != nil {
		utils.Warning.Println(err)
		return ""
	}

	return hostEnv
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
				ID:   ID,
				Name: name,
				Env:  env,
			})
		}

	}
	return result
}

func getTrends(hostID int, ctx *user.GlobalCTX) Dashboard {
	var (
		trend string
		s     int
		f     int
		rf    int
		t     int
		last  string
		dash  Dashboard
	)

	stmt, err := ctx.Config.Database.DB.Prepare("select trend, success, failed, rFailed, total, last from hosts where id = ?")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(hostID).Scan(&trend, &s, &f, &rf, &t, &last)
	if err != nil {
		utils.Warning.Printf("%q, getDashboardData", err)
	}

	dash.Summary = t
	dash.Success = s
	dash.RFailed = rf
	dash.Failed = f

	trendStruct := make(map[int]int)

	var trendStr []string
	_ = json.Unmarshal([]byte(trend), &trendStr)

	for _, i := range trendStr {
		splt := strings.Split(i, ":")

		l, _ := strconv.Atoi(splt[0])
		v, _ := strconv.Atoi(splt[1])

		trendStruct[l] = v

	}
	dash.Trend = trendStruct

	hostsStruct := make(map[string]bool)
	var hostsStr []string
	_ = json.Unmarshal([]byte(last), &hostsStr)

	for _, i := range hostsStr {
		splt := strings.Split(i, ":")

		v, _ := strconv.ParseBool(splt[1])

		hostsStruct[splt[0]] = v
		fmt.Println(splt[0], v)
	}

	fmt.Println(hostsStruct)

	dash.LastHosts = hostsStruct

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
	stmt, err := ctx.Config.Database.DB.Prepare("update hosts set `trend` = ?, `Last` = ?, `Success` = ?, `Failed` = ?, `RFailed` = ?, `Total` = ? where (`id` = ?)")
	if err != nil {
		utils.Warning.Println(err)
	}

	var tmpTrends []string
	var tmpLastHost []string

	for h, c := range data.Trend {
		tmpTrends = append(tmpTrends, fmt.Sprintf("%d:%d", h, c))
	}

	for h, c := range data.LastHosts {
		tmpLastHost = append(tmpLastHost, fmt.Sprintf("%s:%t", h, c))
	}

	jsonStrTrend, _ := json.Marshal(tmpTrends)
	jsonStrHosts, _ := json.Marshal(tmpLastHost)

	_, err = stmt.Exec(jsonStrTrend, jsonStrHosts, data.Success, data.Failed, data.RFailed, data.Summary, hostID)
	if err != nil {
		utils.Warning.Println(err)
	}
	utils.DeferCloseStmt(stmt)
}
