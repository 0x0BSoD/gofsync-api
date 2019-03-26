package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	webServer bool
	file      string
	conf      string
	host      string
	count     string
	parallel  bool
	tosync    bool
)

type Config struct {
	Actions  []string
	Hosts    []string
	RTPro    string
	RTStage  string
	Username string
	Pass     string
	Port     int
	DBFile   string
	PerPage  int
	DbInit   string
	DB       *sql.DB
}

var globConf Config

// =====================
//  Args
// =====================
func (a *Config) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	a.DB.SetMaxIdleConns(140)
	a.DB.SetMaxOpenConns(100)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&conf, "conf", "", "Config file, TOML")
	flag.BoolVar(&webServer, "server", false, "Run as web server daemon")
	flag.BoolVar(&parallel, "parallel", false, "Parallel run")
	flag.StringVar(&file, "file", "", "File contain hosts divide by new line")
	flag.StringVar(&host, "host", "", "Foreman FQDN")
	flag.BoolVar(&tosync, "sync", false, "Sync Foreman[s] -synconf required")
}

// Return Auth structure with Username and Password for Foreman api
func configParser() {
	viper.SetConfigName(conf)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found...")
	} else {
		globConf = Config{
			Username: viper.GetString("API.username"),
			Pass:     viper.GetString("API.password"),
			DBFile:   viper.GetString("DB.db_file"),
			Actions:  viper.GetStringSlice("RUNNING.actions"),
			RTPro:    viper.GetString("RT.pro"),
			RTStage:  viper.GetString("RT.stage"),
			PerPage:  viper.GetInt("RUNNING.per_page_def"),
			DbInit:   viper.GetString("DB.init_file"),
		}
		globConf.Initialize(viper.GetString("DB.db_user"),
			viper.GetString("DB.db_password"),
			viper.GetString("DB.db_schema"))
	}
}

func getHosts(file string) {
	if len(file) > 0 {
		// Get hosts from file
		var hosts []byte
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("Not file: %v\n", err)
		}
		hosts, _ = ioutil.ReadAll(f)
		tmpHosts := strings.Split(string(hosts), "\n")
		var sHosts []string
		for _, i := range tmpHosts {
			if !strings.HasPrefix(i, "#") && len(i) > 0 {
				sHosts = append(sHosts, i)
				fmt.Println(i)

			}
		}
		globConf.Hosts = sHosts
	} else {
		log.Fatal("")
	}
}

func main() {
	flag.Parse()
	configParser()
	getHosts(file)
	if webServer {
		Server()
	} else {
		// =========================
		//for _, host := range hosts {
		//fmt.Println(host)
		//}
		if parallel {
			//fullSync(globConf.Hosts)
			saveHGToJson()
		}
		//		// Foremans
		//		mustRunParr(sHosts, count)
		//		// RT
		//		getRTHostGroups("rt.stage.ringcentral.com")
		//		getRTHostGroups("rt.ringcentral.com")
		//	} else {
		//
		//		// Foremans
		//		mustRun(sHosts)
		//		// RT
		//		getRTHostGroups("rt.stage.ringcentral.com")
		//		getRTHostGroups("rt.ringcentral.com")
		//	}
		//} else {
		//fmt.Println(host)
		//kostyl := []string{host}
		//mustRun(kostyl)
	}

}
