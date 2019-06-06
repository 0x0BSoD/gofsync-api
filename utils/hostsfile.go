package utils

import (
	"fmt"
	mod "git.ringcentral.com/archops/goFsync/models"
	"github.com/tatsushid/go-fastping"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func GetHosts(file string, cfg *mod.Config) {
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

		p := fastping.NewPinger()
		for _, i := range tmpHosts {
			if !strings.HasPrefix(i, "#") && len(i) > 0 {
				ra, err := net.ResolveIPAddr("ip4:icmp", i)
				if err != nil {
					fmt.Println(err)
					continue
				}
				p.AddIPAddr(ra)
				p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
					//fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
					if !StringInSlice(i, sHosts) {
						sHosts = append(sHosts, i)
					}
				}
				p.OnIdle = func() {
					fmt.Println("finish")
				}
				err = p.Run()
				if err != nil {
					fmt.Println(err)
				}
				//sHosts = append(sHosts, i)
			}
		}
		fmt.Println(sHosts)
		cfg.Hosts = sHosts
	} else {
		log.Fatal("Hosts file not found...")
	}
}
