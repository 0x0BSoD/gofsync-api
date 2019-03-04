package main

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/foremanGetter/entitys"
	"strings"
	"sync"
)

func parallelGetLoc(sHosts []string) {
	fmt.Println("Getting Locations")
	var wg sync.WaitGroup
	hostsCount := len(sHosts)
	wg.Add(hostsCount)
	for i := 0; i < hostsCount; i++ {
		host := sHosts[i]

		go func(i int, host string) {
			if !strings.HasPrefix(host, "#") {
				getLocations(host)
				fmt.Println(host, " Done.")
			}
			defer wg.Done()
		}(i, host)

	}
	wg.Wait()
	fmt.Println("Complete! Getting Locations")

}

func parallelGetHostGroups(sHosts []string, count string) {
	fmt.Println("Getting Host Groups")

	var wg sync.WaitGroup
	hostsCount := len(sHosts)
	// TODO: Make FIFO mq
	var results []entitys.ChanSWE
	c := make(chan entitys.ChanSWE)
	wg.Add(hostsCount)

	for i := 0; i < hostsCount; i++ {
		host := sHosts[i]
		go func(i int, host string, c chan entitys.ChanSWE) {
			if !strings.HasPrefix(host, "#") {
				go getHostGroups(host, count)
			}
		}(i, host, c)
	}
	results = append(results, <-c)

	wg.Wait()
	fmt.Println("Complete! Host Groups")
	fmt.Println("=============================")
	for _, SWE := range results {
		fmt.Println(SWE)
	}

}
