package info

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"strconv"
	"strings"
)

func Get(host string, ctx *user.GlobalCTX) Dashboard {
	var trend string
	var s int
	var f int
	var rf int
	var t int
	var last string
	var dash Dashboard

	stmt, err := ctx.Config.Database.DB.Prepare("select trend, Success, Failed, RFailed, Total, Last from hosts where host = ?")
	if err != nil {
		utils.Warning.Println(err)
	}
	defer utils.DeferCloseStmt(stmt)

	err = stmt.QueryRow(host).Scan(&trend, &s, &f, &rf, &t, &last)
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
		l, _ := strconv.Atoi(splt[0])
		v, _ := strconv.Atoi(splt[1])
		trendStruct.Labels = append(trendStruct.Labels, l)
		trendStruct.Values = append(trendStruct.Values, v)
	}
	fmt.Println(trendStruct)
	dash.Trend = trendStruct

	return dash
}

func Update(host string, data Dashboard, ctx *user.GlobalCTX) {
	stmt, err := ctx.Config.Database.DB.Prepare("UPDATE hosts SET `trend` = ?, `Last` = ?, `Success` = ?, `Failed` = ?, `RFailed` = ?, `Total` = ? WHERE (`host` = ?)")
	if err != nil {
		utils.Warning.Println(err)
	}

	var tmp []string

	for idx, l := range data.Trend.Labels {
		tmp = append(tmp, fmt.Sprintf("%d:%d", l, data.Trend.Values[idx]))
	}

	jsonStr, err := json.Marshal(tmp)
	if err != nil {
		utils.Error.Println(err)
	}

	_, err = stmt.Exec(jsonStr, data.LastHost, data.Success, data.Failed, data.RFailed, data.Summary, host)
	if err != nil {
		utils.Warning.Println(err)
	}
	utils.DeferCloseStmt(stmt)
}
