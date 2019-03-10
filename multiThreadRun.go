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
			getHostGroups(host, count)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Host Groups")
	fmt.Println("=============================")

}

func parallelGetPuppetClasses(sHosts []string, count string) {
	fmt.Println("Getting Puppet Classes")
	var wg sync.WaitGroup

	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			getPuppetClasses(host, count)
		}(host)
	}
	wg.Wait()
	fmt.Println("Complete! Puppet Classes")
	fmt.Println("=============================")
}

func parallelGetSmartClasses(sHosts []string) {
	fmt.Println("Getting Puppet Smart Classes")
	var wg sync.WaitGroup

	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			InsertPuppetSmartClasses(host)
		}(host)
	}
	wg.Wait()
	fmt.Println("Complete! Puppet Smart Classes")
	fmt.Println("=============================")
}

func parallelGetOverrideBase(sHosts []string) {
	fmt.Println("Getting Override Parameters Base")
	var wg sync.WaitGroup

	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			InsertToOverridesBase(host)
		}(host)
	}
	wg.Wait()
	fmt.Println("Complete! Override Parameters Base")
	fmt.Println("=============================")
}

func parallelGetOverrideP(sHosts []string) {
	fmt.Println("Getting Override Parameters")
	var wg sync.WaitGroup

	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			InsertOverridesParameters(host)
		}(host)
	}
	wg.Wait()
	fmt.Println("Complete! Override Parameters")
	fmt.Println("=============================")
}

// =================================================================
// RUN
// =================================================================
func mustRunParr(sHosts []string, count string) {
	actions := Config.Actions
	if stringInSlice("dbinit", actions) {
		dbActions()
	}
	if stringInSlice("locations", actions) {
		parallelGetLoc(sHosts)
	}
	if stringInSlice("swes", actions) {
		parallelGetHostGroups(sHosts, count)
	}
	if stringInSlice("pclasses", actions) {
		parallelGetPuppetClasses(sHosts, count)
	}
	if stringInSlice("sclasses", actions) {
		parallelGetSmartClasses(sHosts)
	}
	if stringInSlice("overridebase", actions) {
		parallelGetOverrideBase(sHosts)
	}
	if stringInSlice("overrideparams", actions) {
		parallelGetOverrideP(sHosts)
	}
	if stringInSlice("swefill", actions) {
		fillSWEtable()
	}
	if stringInSlice("swecheck", actions) {
		crossCheck()
	}
}
