package utils

import (
	mod "git.ringcentral.com/archops/goFsync/models"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GetHosts(file string, globConf *mod.Config) {
	if file != "" && len(file) > 0 {
		// Get hosts from file
		var hosts []byte
		f, err := os.Open(file)
		if err != nil {
			Error.Fatal("Hosts file not found...")
		}

		hosts, _ = ioutil.ReadAll(f)
		tmpHosts := strings.Split(string(hosts), "\n")
		var sHosts []string

		for _, i := range tmpHosts {
			if !strings.HasPrefix(i, "#") && len(i) > 0 {
				sHosts = append(sHosts, i)
			}
		}
		globConf.Hosts = sHosts
	} else {
		log.Fatal("Hosts file not found...")
	}
}
