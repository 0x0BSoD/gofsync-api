package main

import (
	"fmt"
	"log"
	"sync"
)

// =================================================================
// RUN
// =================================================================
func fullSync() {
	//dbActions()
	parallelGetLoc(globConf.Hosts)
	parallelGetEnv(globConf.Hosts)
	parallelGetPuppetClasses(globConf.Hosts)
	parallelGetSmartClasses(globConf.Hosts)
	parallelGetHostGroups(globConf.Hosts)
	parallelUpdatePC(globConf.Hosts)
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

			result, err := locations(host)
			if err != nil {
				log.Printf("Error on getting Locations:\n%q", err)
			}

			for _, loc := range result.Results {
				insertToLocations(host, loc.Name, loc.ID)
			}
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

			result, err := environments(host)
			if err != nil {
				log.Printf("Error on getting Environments:\n%q", err)
			}

			for _, env := range result.Results {
				insertToEnvironments(host, env.Name, env.ID)
			}
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

			result, err := getAllPC(host)
			if err != nil {
				log.Printf("Error on getting Puppet classes:\n%q", err)
			}

			for className, subClasses := range result {
				for _, subClass := range subClasses {
					insertPC(host, className, subClass.Name, subClass.ID)
				}
			}
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! PuppetClasses")
	fmt.Println("=============================")
}

func parallelGetHostGroups(sHosts []string) {
	fmt.Println("Getting HostGroups")

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

	fmt.Println("Complete! Getting HostGroups")
	fmt.Println("=============================")
}

func parallelGetSmartClasses(sHosts []string) {
	fmt.Println("Getting Smart Classes and Overrides")

	var wg sync.WaitGroup
	for _, host := range sHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			fmt.Println("==> ", host)

			result, err := smartClasses(host)
			if err != nil {
				log.Printf("Error on getting Smart Classes and Overrides:\n%q", err)
			}

			for _, i := range result {
				lastID := insertSC(host, i)
				if lastID != -1 {
					// Getting data by Foreman Smart Class ID
					ovrResult := scOverridesById(host, i.ID)
					for _, ovr := range ovrResult {
						// Storing data by internal SmartClass ID
						insertSCOverride(lastID, ovr, i.ParameterType)
					}
				}
			}
		}(host)
	}
	wg.Wait()

	fmt.Println("Complete! Getting Smart Classes and Overrides")
	fmt.Println("=============================")
}

func parallelUpdatePC(sHosts []string) {
	fmt.Println("Getting ids for Puppet Class from target foreman host")

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

	fmt.Println("Complete!")
	fmt.Println("=============================")
}
