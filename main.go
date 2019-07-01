package main

import (
	"flag"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/user"
	cfg "git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"strings"
)

var globConf = cfg.Config{}

var (
	webServer bool
	test      bool
	file      string
	conf      string
	action    string
)

// =====================
//  Args
// =====================
func init() {
	flag.BoolVar(&test, "test", false, "run compare")
	flag.StringVar(&conf, "conf", "", "Config file, TOML")
	flag.StringVar(&file, "hosts", "", "File contain hosts divide by new line")
	flag.StringVar(&action, "action", "", "If specified run one of env|loc|pc|sc|hg|pcu")
	//flag.BoolVar(&globConf.Web.SocketActive, "socket", false, "Run socket server")
	flag.BoolVar(&webServer, "server", false, "Run as web server daemon")
}

func main() {
	flag.Parse()
	//fmt.Println(globConf.Web.SocketActive)
	// Params and DB =================
	utils.Parser(&globConf, conf)
	utils.InitializeDB(&globConf)
	globConf.Sessions = user.CreateHub()
	// Logging =======================
	utils.Init(&globConf.Logging.TraceLog,
		&globConf.Logging.AccessLog,
		&globConf.Logging.ErrorLog,
		&globConf.Logging.ErrorLog)

	utils.GetHosts(file, &globConf)
	//utils.InitRedis(&globConf)
	hostgroups.StoreHosts(&globConf)

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
		Server(&globConf)
	} else if test {
		utils.GetInfo()
		//hostgroups.Compare(&globConf)
	} else {
		session := user.Start(&cfg.Claims{Username: "srv_foreman"}, "fake", &globConf)
		if strings.Contains(action, ",") {
			actions := strings.Split(action, ",")
			for _, a := range actions {
				switch a {
				case "loc":
					locSync(&session)
				case "env":
					envSync(&session)
				case "pc":
					puppetClassSync(&session)
				case "sc":
					smartClassSync(&session)
				case "hg":
					hostGroupsSync(&session)
				case "pcu":
					puppetClassUpdate(&session)
				}
			}
		} else {
			switch action {
			case "loc":
				locSync(&session)
			case "env":
				envSync(&session)
			case "pc":
				puppetClassSync(&session)
			case "sc":
				smartClassSync(&session)
			case "hg":
				hostGroupsSync(&session)
			case "pcu":
				puppetClassUpdate(&session)
			default:
				fullSync(&session)
			}
		}
	}
}
