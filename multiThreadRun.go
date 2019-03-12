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

func parallelGetHostGroups(sHosts []string) {
	fmt.Println("Getting Host Groups")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			hg := SWE{}
			hg.Get(host)
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
		parallelGetHostGroups(sHosts)
	}
}

func fullSync(sHosts []string) {
	dbActions()
	parallelGetLoc(sHosts)
	parallelGetHostGroups(sHosts)
}
