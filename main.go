package main

import (
	"flag"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
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
	flag.BoolVar(&globConf.Web.SocketActive, "socket", false, "Run socket server")
	flag.BoolVar(&webServer, "server", false, "Run as web server daemon")
}

func main() {
	flag.Parse()
	fmt.Println(globConf.Web.SocketActive)
	// Params and DB =================
	utils.Parser(&globConf, conf)
	utils.InitializeDB(&globConf)
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
		hostgroups.Compare(&globConf)
	} else {
		if strings.Contains(action, ",") {
			actions := strings.Split(action, ",")
			for _, a := range actions {
				switch a {
				case "loc":
					locSync(&globConf)
				case "env":
					envSync(&globConf)
				case "pc":
					puppetClassSync(&globConf)
				case "sc":
					smartClassSync(&globConf)
				case "hg":
					hostGroupsSync(&globConf)
				case "pcu":
					puppetClassUpdate(&globConf)
				}
			}
		} else {
			switch action {
			case "loc":
				locSync(&globConf)
			case "env":
				envSync(&globConf)
			case "pc":
				puppetClassSync(&globConf)
			case "sc":
				smartClassSync(&globConf)
			case "hg":
				hostGroupsSync(&globConf)
			case "pcu":
				puppetClassUpdate(&globConf)
			default:
				fullSync(&globConf)
			}
		}
	}
}
