package main

import (
	"fmt"
	"sync"
)

func parallelGetLoc(sHosts []string) {
	fmt.Println("Getting Locations")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			getLocations(host)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Getting Locations")
	fmt.Println("=============================")
}

func parallelGetHostGroups(sHosts []string, count string) {
	fmt.Println("Getting Host Groups")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			hg := SWE{}
			hg.Get(host, count).Save(host)
			//getHostGroups(host, count)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Host Groups")
	fmt.Println("=============================")

}

// =================================================================
// RUN
// =================================================================
func mustRunParr(sHosts []string, count string) {
	actions := globConf.Actions
	if stringInSlice("dbinit", actions) {
		dbActions()
	}
	if stringInSlice("locations", actions) {
		parallelGetLoc(sHosts)
	}
	if stringInSlice("swes", actions) {
		parallelGetHostGroups(sHosts, count)
	}
}

func fullSync(sHosts []string, count string) {
	dbActions()
	parallelGetLoc(sHosts)
	parallelGetHostGroups(sHosts, count)
}
