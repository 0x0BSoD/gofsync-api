package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func fillTableSWEState() {
	list := getAllSWE()

	for _, item := range list {
		insertSWEs(item)
	}
}

func checkSWEState() {
	rtSwes := []string{Config.RTPro, Config.RTStage}

	hosts := "./hosts"
	SWElist := getAllSWE()

	// Check RT SWEs
	for _, _host := range rtSwes {
		for _, SWE := range SWElist {
			state := SWEstate(_host, SWE)
			if state {
				insertSWEState(_host, SWE, "OK")
			} else {
				insertSWEState(_host, SWE, "NONE")
			}
		}
	}

	// Check SWEs on hosts
	f, err := os.Open(hosts)
	if err != nil {
		log.Fatalf("Not file: %v\n", err)
	}
	hostsList, _ := ioutil.ReadAll(f)
	sHosts := strings.Split(string(hostsList), "\n")
	for _, _host := range sHosts {
		if !strings.HasPrefix(_host, "#") {
			for _, SWE := range SWElist {
				state := SWEstate(_host, SWE)
				rtPro := SWEstate(Config.RTPro, SWE)
				rtStage := SWEstate(Config.RTStage, SWE)
				if state {
					strState := "OK"
					if rtPro {
						strState += "_PROD"
					}
					if rtStage {
						strState += "_STAGE"
					}
					if !rtPro && !rtStage {
						strState = "NOTINRT_ONHOST"
					}
					insertSWEState(_host, SWE, strState)
				} else {
					strState := "NOT"
					if rtPro {
						strState += "_PROD"
					}
					if rtStage {
						strState += "_STAGE"
					}
					if !rtPro && !rtStage {
						strState = "NOTINRT"
					}
					insertSWEState(_host, SWE, strState)
				}
			}
		}
	}
}
