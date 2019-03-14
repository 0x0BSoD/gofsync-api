package main

import (
	"flag"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	webServer bool
	file      string
	synConf   string
	host      string
	count     string
	parallel  bool
	tosync    bool
)

type Config struct {
	Actions  []string
	RTPro    string
	RTStage  string
	Username string
	Pass     string
	Port     int
	DBFile   string
	PerPage  int
}

var globConf Config

// =====================
//  Args
// =====================
func init() {
	const (
		//defaultCount   = "10"
		//usageCount     = "Pulled items"
		usageWebServer = "Run as web server daemon"
	)
	flag.StringVar(&synConf, "synconf", "", "Config file for sync, TOML")
	flag.BoolVar(&webServer, "server", false, usageWebServer)
	flag.BoolVar(&parallel, "parallel", false, "Parallel run")
	flag.StringVar(&file, "file", "", "File contain hosts divide by new line")
	flag.StringVar(&host, "host", "", "Foreman FQDN")
	flag.BoolVar(&tosync, "sync", false, "Sync Foreman[s] -synconf required")
}

// Return Auth structure with Username and Password for Foreman api
func configParser() {
	viper.SetConfigName("config")
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
		}
	}
}

func main() {
	flag.Parse()
	configParser()
	if webServer {
		Server()
	} else {
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
				if !strings.HasPrefix(i, "#") {
					sHosts = append(sHosts, i)
				}
			}
			// =========================
			//for _, host := range hosts {
			//fmt.Println(host)
			//}
			if parallel {
				fullSync(sHosts)
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

}
