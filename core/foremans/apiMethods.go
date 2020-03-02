package foremans

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
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

		success := 0
		rFailed := 0
		failed := 0
		total := 0

		hostsMap := make(map[string]bool)
		trendMap := make(map[int]int)
		for i := 1; i <= 24; i++ {
			trendMap[i] = 0
		}

		for _, item := range r.Results {
			if item.Status.Failed != 0 {
				failed++
				hostsMap[item.HostName] = false
			} else if item.Status.FailedRestarts != 0 {
				rFailed++
				hostsMap[item.HostName] = true
			} else {
				success++
				hostsMap[item.HostName] = true
			}

			_time := strings.Split(item.ReportedAt, "T")[1]
			hour := strings.Split(_time, ":")[0]
			sInt, _ := strconv.Atoi(hour)
			trendMap[sInt]++
			total++
		}

		dashboard.LastHosts = hostsMap
		dashboard.Trend = trendMap
		dashboard.Failed = failed
		dashboard.RFailed = rFailed
		dashboard.Success = success
		dashboard.Summary = total
	}

	return dashboard
}
