package foremans

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"sort"
	"strconv"
	"strings"
)

// ===============================
// GET
// ===============================
func ApiReportsDaily(host string, ctx *user.GlobalCTX) Dashboard {
	response, err := utils.ForemanAPI("GET", host, "reports?search=reported+%3D+Today&order=reported+DESC", "", ctx)
	if err != nil {
		utils.Error.Println(err)
	}

	var r Container
	err = json.Unmarshal(response.Body, &r)
	if err != nil {
		utils.Error.Println(err)
	}

	var dashboard Dashboard
	if len(r.Results) > 0 {

		dashboard.LastHost = r.Results[0].HostName

		s := 0
		rf := 0
		f := 0
		t := 0

		trendMap := make(map[int]int)

		for _, item := range r.Results {
			if item.Status.Failed != 0 {
				f++
			} else if item.Status.FailedRestarts != 0 {
				rf++
			} else {
				s++
			}

			_time := strings.Split(item.ReportedAt, "T")[1]
			hour := strings.Split(_time, ":")[0]
			sInt, _ := strconv.Atoi(hour)
			trendMap[sInt]++
			t++
		}

		for hour := range trendMap {
			dashboard.Trend.Labels = append(dashboard.Trend.Labels, hour)
		}
		sort.Ints(dashboard.Trend.Labels)

		for _, hour := range dashboard.Trend.Labels {
			dashboard.Trend.Values = append(dashboard.Trend.Values, trendMap[hour])
		}

		dashboard.Failed = f
		dashboard.RFailed = rf
		dashboard.Success = s
		dashboard.Summary = t
	}

	return dashboard
}
