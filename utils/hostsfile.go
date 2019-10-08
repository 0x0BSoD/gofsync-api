package utils

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/tatsushid/go-fastping"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func GetHosts(file string, cfg *models.Config) {
	var hosts []byte
	var f io.Reader
	var err error

	f, err = os.Open(file)
	if err != nil {
		Warning.Println("hosts parameter not found try to use 'hosts' file in the same dir")
		f, err = os.Open("./hosts")
		if err != nil {
			Error.Fatal("hosts file not found ...")
		}
	}

	hosts, _ = ioutil.ReadAll(f)
	tmpHosts := strings.Split(string(hosts), "\n")
	var sHosts []string

	p := fastping.NewPinger()
	for _, i := range tmpHosts {
		if !strings.HasPrefix(i, "#") && len(i) > 0 {
			if !StringInSlice(i, sHosts) {
				sHosts = append(sHosts, i)
			}
			err = p.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	cfg.Hosts = sHosts
}
