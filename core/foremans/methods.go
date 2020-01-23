package foremans

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/tatsushid/go-fastping"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func Sync(hostname string, ctx *user.GlobalCTX) {

	// Step LOG to stdout ======================
	fmt.Println(utils.PrintJsonStep(models.Step{
		Actions: "Getting Locations Info",
		Host:    hostname,
	}))
	// =========================================

	// from foreman
	locationsResult := ApiReportsDaily(hostname, ctx)
	UpdateTrends(ctx.Config.Hosts[hostname], locationsResult, ctx)

	// Socket Broadcast ---
	//ctx.Broadcast(models.WSMessage{
	//	Broadcast: true,
	//	Operation: "dashboardUpdate",
	//	Data: models.Step{
	//		Host:   hostname,
	//		Status: "updated",
	//	},
	//})
	// ---

}

func StoreHosts(file string, cfg *models.Config) {
	var hosts []byte
	var f io.Reader
	var err error
	cfg.Hosts = make(map[string]int)

	f, err = os.Open(file)
	if err != nil {
		utils.Warning.Println("hosts parameter not found try to use 'hosts' file in the same dir")
		f, err = os.Open("./hosts")
		if err != nil {
			utils.Error.Fatal("hosts file not found ...")
		}
	}

	hosts, _ = ioutil.ReadAll(f)
	tmpHosts := strings.Split(string(hosts), "\n")
	var sHosts []string

	p := fastping.NewPinger()
	for _, i := range tmpHosts {
		if !strings.HasPrefix(i, "#") && len(i) > 0 {
			if !utils.StringInSlice(i, sHosts) { // ?
				sHosts = append(sHosts, i)
				ID := InsertHost(i, cfg)
				cfg.Hosts[i] = ID
			}
			err = p.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
