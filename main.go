package main

import (
	"flag"
	"fmt"
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
	var globConf Config
// =====================
//  Args
// =====================
func init() {
	const (
		defaultCount   = "10"
		usageCount     = "Pulled items"
		usageWebServer = "Run as web server daemon"
	)
	flag.StringVar(&synConf, "synconf", "", "Config file for sync, TOML")
	flag.StringVar(&count, "count", defaultCount, usageCount)
	flag.BoolVar(&webServer, "server", false, usageWebServer)
	flag.BoolVar(&parallel, "parallel", false, "Parallel run")
	flag.StringVar(&file, "file", "", "File contain hosts divide by new line")
	flag.StringVar(&host, "host", "", "Foreman FQDN")
	flag.BoolVar(&tosync, "sync", false, "Sync Foreman[s] -synconf required")
}

// Return Auth structure with Username and Password for Foreman api
func configParser() {
	var dbFile string
	var username string
	var pass string
	var actions []string
	var rtPro string
	var rtStage string

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found...")
	} else {
		dbFile   = viper.GetString("DB.db_file")
		username = viper.GetString("API.username")
		pass     = viper.GetString("API.password")
		actions  = viper.GetStringSlice("RUNNING.actions")
		rtPro    = viper.GetString("RT.pro")
		rtStage  = viper.GetString("RT.stage")
	}

	globConf = Config{
		Username: username,
		Pass:     pass,
		DBFile:   dbFile,
		Actions:  actions,
		RTPro:    rtPro,
		RTStage:  rtStage,
	}
}

func main() {
	flag.Parse()
	configParser()
	if tosync {
		_, err := os.Stat(synConf)
		if err != nil {
			log.Fatalf("Fatal error Sync config file: %s \nParam -synconf required", err)
		}
		name := strings.Split(synConf, ".")
		viper.SetConfigName(name[0])
		viper.AddConfigPath(".")
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error Sync config file: %s \n", err))
		}
		//sSource := viper.GetString("source.host")
		//sDest := viper.GetStringSlice("dest.host")
		//sSWEs := viper.GetStringSlice("params.swes")
		//sState := viper.GetStringSlice("params.swes_state")
		//runSync(sSource, sDest, sState, sSWEs)

	} else if webServer {
		log.Fatal("Not implemented\n")
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
			rtSWEs := RTSWE{}
			jData := rtSWEs.Get("rt-sndbx.lab.nordigy.ru").ToJSON()
			fmt.Print(jData)
			//for _, host := range hosts {
				//fmt.Println(host)
			//}
		//	if parallel {
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
