package main

import (
	"flag"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"strings"
	"sync"
)

var globConf models.Config
var globSession user.GlobalCTX

var (
	webServer bool
	dashupd   bool
	repo      bool
	file      string
	conf      string
	action    string
)

// =====================
//  Args
// =====================
func init() {
	flag.BoolVar(&dashupd, "dashupd", false, "Update dashboard data")
	flag.BoolVar(&repo, "repo", false, "Init Git repo")
	flag.StringVar(&conf, "conf", "", "Config file, TOML")
	flag.StringVar(&file, "hosts", "", "File contain hosts divide by new line")
	flag.StringVar(&action, "action", "", "If specified run one of env|loc|pc|sc|hg|pcu")
	//flag.BoolVar(&globConf.Web.SocketActive, "socket", false, "Run socket server")
	flag.BoolVar(&webServer, "server", false, "Run as web server daemon")
}

func main() {

	gl := &sync.Mutex{}
	globSession.GlobalLock = gl

	flag.Parse()
	// Params and DB =================
	utils.Parser(&globConf, conf)
	utils.InitializeDB(&globConf)
	globSession.Sessions = user.CreateHub()
	// Logging =======================
	utils.Init(&globConf.Logging.TraceLog,
		&globConf.Logging.AccessLog,
		&globConf.Logging.ErrorLog,
		&globConf.Logging.ErrorLog)

	utils.GetHosts(file, &globConf)
	//utils.InitRedis(&globConf)
	hostgroups.StoreHosts(&globConf)

	// Set global config to global sessions container
	globSession.Config = globConf

	if webServer {
		hello := `
|￣￣￣￣￣￣￣￣|
| goFsync_api    |
|＿＿＿＿＿＿＿＿|
(\__/) ||
(•ㅅ•) ||
/ 　 づ`
		fmt.Println(hello)
		fmt.Printf("running on port %d\n", globConf.Web.Port)
		go startScheduler(&globSession)
		Server(&globSession)
	} else if dashupd {
		DashboardUpdate(&globSession)
	} else if repo {
		utils.InitRepo(&globSession)
	} else {
		globSession.Set(&user.Claims{Username: "srv_foreman"}, "fake")
		if strings.Contains(action, ",") {
			actions := strings.Split(action, ",")
			for _, a := range actions {
				switch a {
				case "loc":
					locSync(&globSession)
				case "env":
					envSync(&globSession)
				case "pc":
					puppetClassSync(&globSession)
				case "sc":
					smartClassSync(&globSession)
				case "hg":
					hostGroupsSync(&globSession)
				case "pcu":
					puppetClassUpdate(&globSession)
				}
			}
		} else {
			switch action {
			case "loc":
				locSync(&globSession)
			case "env":
				envSync(&globSession)
			case "pc":
				puppetClassSync(&globSession)
			case "sc":
				smartClassSync(&globSession)
			case "hg":
				hostGroupsSync(&globSession)
			case "pcu":
				puppetClassUpdate(&globSession)
			default:
				fullSync(&globSession)
			}
		}
	}
}
