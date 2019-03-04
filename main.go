package main

import (
	"flag"
	"fmt"
	"git.ringcentral.com/alexander.simonov/foremanGetter/entitys"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
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

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found...")
	} else {
		dbFile = viper.GetString("DB.db_file")
		username = viper.GetString("API.username")
		pass = viper.GetString("API.password")
	}

	auth := entitys.Auth{
		Username: username,
		Pass:     pass,
		DBFile:   dbFile,
	}
	Config = auth
}

func getDeltaTime(start time.Time) string {
	delta := time.Since(start)
	res := fmt.Sprint(delta.String())
	return res
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
			sHosts := strings.Split(string(hosts), "\n")
			if parallel {
				// RT
				getRTHostGroups("rt.stage.ringcentral.com")
				getRTHostGroups("rt.ringcentral.com")
				// Foremans
				parallelGetLoc(sHosts)
				parallelGetHostGroups(sHosts, count)
			} else {
				getRTHostGroups("rt.stage.ringcentral.com")
				getRTHostGroups("rt.ringcentral.com")
				for _, host := range sHosts {
					if !strings.HasPrefix(host, "#") {
						initialRun(host)
					}
					//crossCheck()
				}
			}
		} else {
			initialRun(host)
			crossCheck()
		}
	}

}
