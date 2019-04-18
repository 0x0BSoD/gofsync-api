package main

import (
	"flag"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/core/hostgroups"
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
)

// =====================
//  Args
// =====================
func init() {
	flag.StringVar(&conf, "conf", "", "Config file, TOML")
	flag.StringVar(&file, "file", "", "File contain hosts divide by new line")
	flag.StringVar(&host, "host", "", "Foreman FQDN")
	flag.BoolVar(&webServer, "server", false, "Run as web server daemon")
}

func main() {
	flag.Parse()

	// Params and DB =================
	utils.Parser(&globConf, conf)
	utils.InitializeDB(&globConf)
	utils.GetHosts(file, &globConf)
	// Logging =======================
	utils.Init(&globConf.Logging.TraceLog,
		&globConf.Logging.AccessLog,
		&globConf.Logging.ErrorLog,
		&globConf.Logging.ErrorLog)

	if webServer {
		hello := `
|￣￣￣￣￣￣￣ |
| goFsync_api   |
|＿＿＿＿＿ _＿_|
(\__/) ||
(•ㅅ•) ||
/ 　 づ`
		fmt.Println(hello)
		fmt.Printf("running on port %d\n", globConf.Web.Port)
		Server(&globConf)
	} else {
		fullSync(&globConf)
		hostgroups.SaveHGToJson(&globConf)
	}
}
