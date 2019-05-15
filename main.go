package main

import (
	"flag"
	"fmt"
	cfg "git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	_ "github.com/go-sql-driver/mysql"
)

var globConf = cfg.Config{}

var (
	webServer bool
	file      string
	conf      string
	host      string
	action    string
)

// =====================
//  Args
// =====================
func init() {
	flag.StringVar(&conf, "conf", "", "Config file, TOML")
	flag.StringVar(&file, "file", "", "File contain hosts divide by new line")
	flag.StringVar(&host, "host", "", "Foreman FQDN")
	flag.StringVar(&action, "action", "", "If specified run one of env|loc|pc|sc|hg|pcu")
	flag.BoolVar(&webServer, "server", false, "Run as web server daemon")
}

func main() {
	flag.Parse()

	// Params and DB =================
	utils.Parser(&globConf, conf)
	utils.InitializeDB(&globConf)
	//utils.InitializeAMQP(&globConf)
	utils.GetHosts(file, &globConf)
	// Logging =======================
	utils.Init(&globConf.Logging.TraceLog,
		&globConf.Logging.AccessLog,
		&globConf.Logging.ErrorLog,
		&globConf.Logging.ErrorLog)

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
