package main

import (
	"fmt"
	"sync"
)

// =================================================================
// RUN
// =================================================================
func mustRunParr(sHosts []string) {
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
	//parallelGetLoc(sHosts)
	//parallelGetEnv(sHosts)
	//parallelGetPuppetClasses(sHosts)
	//parallelGetSmartClasses(sHosts)
	//parallelGetHostGroups(sHosts)
	parallelUpdatePC(sHosts)
}

// =================================================================
// Functions
// =================================================================
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

func parallelGetEnv(sHosts []string) {
	fmt.Println("Getting Environments")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)

		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			getEnvironment(host)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Getting Environments")
	fmt.Println("=============================")
}

func parallelGetPuppetClasses(sHosts []string) {
	fmt.Println("Getting PuppetClasses")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			getAllPC(host)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! PuppetClasses")
	fmt.Println("=============================")
}

func parallelGetHostGroups(sHosts []string) {
	fmt.Println("Getting PuppetClasses")

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

func parallelGetSmartClasses(sHosts []string) {
	fmt.Println("Getting Smart Classes")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			insertSmartClasses(host)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Smart Classes")
	fmt.Println("=============================")
}

func parallelGetSCOverrides(sHosts []string) {
	fmt.Println("Getting Smart Classes Overrides")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			insertSCOverrides(host)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Smart Classes Overrides")
	fmt.Println("=============================")
}

func parallelUpdatePC(sHosts []string) {
	fmt.Println("Getting Smart Classes Parameters For PC")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			insertSCByPC(host)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Smart Classes Parameters For PC")
	fmt.Println("=============================")
}

func parallelUpdateHG(sHosts []string) {
	fmt.Println("Getting Smart Classes Parameters For PC")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)
			insertSCByPC(host)
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Smart Classes Parameters For PC")
	fmt.Println("=============================")
}
