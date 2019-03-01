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
		//fmt.Println(item)
		insertSWEs(item)
	}
}

func checkSWEState() {

	hosts := "./hosts"
	SWElist := getAllSWE()

	f, err := os.Open(hosts)
	if err != nil {
		log.Fatalf("Not file: %v\n", err)
	}
	hostsList, _ := ioutil.ReadAll(f)
	sHosts := strings.Split(string(hostsList), "\n")
	for _, _host := range sHosts {
		if !strings.HasPrefix(_host, "#") {
			//fmt.Println(_host)
			for _, SWE := range SWElist {
				state := SWEstate(_host, SWE)
				if state {
					//fmt.Printf("HOST: %s | SWE: %s ==> OK", _host, SWE)
					//fmt.Println()
					insertSWEState(_host, SWE, "OK")
				} else {
					//fmt.Printf("HOST: %s | SWE: %s ==> NOT_SYNC", _host, SWE)
					//fmt.Println()
					insertSWEState(_host, SWE, "NOT_SYNC")
				}
			}
		}
	}
}
