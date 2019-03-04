package main

import (
	"flag"
	"git.ringcentral.com/alexander.simonov/foremanGetter/entitys"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	webServer bool
	file      string
	host      string
	count     string
	parallel  bool
)
var Config entitys.Auth

// =====================
//  Args
// =====================
func init() {
	const (
		defaultCount   = "10"
		usageCount     = "Pulled items"
		usageWebServer = "Run as web server daemon"
	)
	flag.StringVar(&count, "count", defaultCount, usageCount)
	flag.BoolVar(&webServer, "server", false, usageWebServer)
	flag.BoolVar(&parallel, "parallel", false, "Parallel run")
	flag.StringVar(&file, "file", "", "File contain hosts divide by new line")
	flag.StringVar(&host, "host", "", "Foreman FQDN")
}

// Return Auth structure with Username and Password for Foreman api
func configParser() {
	var dbFile string
	var username string
	var pass string
	var actions []string

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found...")
	} else {
		dbFile = viper.GetString("DB.db_file")
		username = viper.GetString("API.username")
		pass = viper.GetString("API.password")
		actions = viper.GetStringSlice("RUNNING.actions")
	}

	auth := entitys.Auth{
		Username: username,
		Pass:     pass,
		DBFile:   dbFile,
		Actions:  actions,
	}
	Config = auth
}

func main() {
	flag.Parse()
	configParser()

	if webServer {
		log.Fatal("Not implemented\n")
	} else {
		if len(file) > 0 {

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

			if parallel {
				// Foremans
				mustRunParr(sHosts, count)
				// RT
				//getRTHostGroups("rt.stage.ringcentral.com")
				//getRTHostGroups("rt.ringcentral.com")
			} else {

				for _, host := range sHosts {
					mustRun(host)
					getRTHostGroups("rt.stage.ringcentral.com")
					getRTHostGroups("rt.ringcentral.com")
					//crossCheck()
				}
			}
		} else {
			mustRun(host)
			crossCheck()
		}
	}

}
