package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"time"
)

func mustRun(hosts []string) {

	overallT := time.Now()
	actions := Config.Actions

	for _, host := range hosts {

		fmt.Println(host)
		fmt.Println("=============================")

		s := spinner.New(spinner.CharSets[25], 100*time.Millisecond)
		st := time.Now()

		if stringInSlice("dbinit", actions) {
			s.Suffix = " Creating DB..."
			s.Start()
			dbActions()
			s.Stop()
			s.FinalMSG = "Complete! Creating DB worked: " + getDeltaTime(st) + "\n"
		}

		if stringInSlice("locations", actions) {
			s.Restart()
			st = time.Now()
			s.Suffix = " Getting Locations..."
			s.Start()
			getLocations(host)
			s.Stop()
			s.FinalMSG = "Complete! Getting Locations worked: " + getDeltaTime(st) + "\n"
		}

		if stringInSlice("swes", actions) {
			s.Restart()
			st = time.Now()
			s.Suffix = " Getting Host Groups..."
			s.Start()
			getHostGroups(host, count)
			s.Stop()
			s.FinalMSG = "Complete! Getting Host Groups worked: " + getDeltaTime(st) + "\n"
		}

		if stringInSlice("pclasses", actions) {
			s.Restart()
			st = time.Now()
			s.Suffix = " Getting Puppet Classes..."
			s.Start()
			getPuppetClasses(host, count)
			s.Stop()
			s.FinalMSG = "Complete! Getting Puppet Classes worked: " + getDeltaTime(st) + "\n"
		}

		if stringInSlice("sclasses", actions) {
			s.Restart()
			st = time.Now()
			s.Suffix = " Filling Smart Classes table..."
			s.Start()
			InsertPuppetSmartClasses(host)
			s.Stop()
			s.FinalMSG = "Complete! Filling Smart Classes table worked: " + getDeltaTime(st) + "\n"
		}
		if stringInSlice("overridebase", actions) {
			s.Restart()
			st = time.Now()
			s.Suffix = " Filling Smart Classes Base table..."
			s.Start()
			InsertToOverridesBase(host)
			s.Stop()
			s.FinalMSG = "Complete! Filling Smart Classes Base table worked: " + getDeltaTime(st) + "\n"
		}

		if stringInSlice("overrideparams", actions) {
			s.Restart()
			st = time.Now()
			s.Suffix = " Filling Smart Classes Overrides parameters table..."
			s.Start()
			InsertOverridesParameters(host)
			s.Stop()
			s.FinalMSG = "Complete! Filling Smart Classes Overrides parameters table worked: " + getDeltaTime(st) + "\n"
		}

		if stringInSlice("oprah", actions) {
			getSWEsByState(host, "NOTINRT_ONHOST")
		}

		fmt.Println()
		sOverall := getDeltaTime(overallT)
		fmt.Println("Done by ", sOverall)
		fmt.Println("=============================")
		fmt.Println()
	}

	if stringInSlice("swefill", actions) {
		fillSWEtable()
	}
	if stringInSlice("swecheck", actions) {
		crossCheck()
	}
}

func fillSWEtable() {
	s := spinner.New(spinner.CharSets[25], 100*time.Millisecond)
	st := time.Now()
	fmt.Println("Actions for all instances")
	s.Restart()
	st = time.Now()
	s.Suffix = " Filling SWE table..."
	s.Start()
	fillTableSWEState()
	s.Stop()
	s.FinalMSG = "Complete! Filling SWE table worked: " + getDeltaTime(st) + "\n"
}

func crossCheck() {
	s := spinner.New(spinner.CharSets[25], 100*time.Millisecond)
	s.Restart()
	st := time.Now()
	s.Suffix = " Checking SWE..."
	s.Start()
	checkSWEState()
	s.Stop()
	s.FinalMSG = "Complete! Checking SWE worked: " + getDeltaTime(st) + "\n"
}
