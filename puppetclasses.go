package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/foremanGetter/entitys"
	"log"
	"sort"
	"strconv"
)

// ===============
// GET
// ===============
func getPuppetClassesByHostgroup(host string, hostgroupID int) {

	var result entitys.PuppetClasses

	bodyText := getForemanAPI(host, "hostgroups/"+strconv.Itoa(hostgroupID)+"/puppetclasses")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	for _, cl := range result.Results {
		var subclassList []string
		for _, v := range cl {
			subclassList = append(subclassList, v.Name)
		}
	}
}

// ===============
// INSERT
// ===============
func getPuppetClasses(host string, count string) {

	var result entitys.PuppetClasses

	bodyText := getForemanAPI(host, "puppetclasses?per_page="+count+"")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	for index, cl := range result.Results {
		var subclassList []string
		for _, v := range cl {
			subclassList = append(subclassList, v.Name)
		}

		sort.Strings(subclassList)

		for _, subclass := range subclassList {
			insertToPupClasses(host, index, subclass)
		}

	}
}
func InsertToOverridesBase(host string) {

	var result entitys.SCPOverride

	aSC := getAllSClasses(host)

	for _, sc := range aSC {

		bodyText := getForemanAPI(host, "smart_class_parameters/"+strconv.Itoa(sc.SCID)+"")

		err := json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Printf("%q:\n %s\n", err, bodyText)
			return
		}

		insertSCOverride(host, result, sc.SCID)

	}
}
func InsertOverridesParameters(host string) {
	Params := getOverrideAllParamBase(host)
	for _, i := range Params {
		var result entitys.OverrideValuesContainer

		bodyText := getForemanAPI(host, "smart_class_parameters/"+strconv.Itoa(i.ClassID)+"/override_values")
		err := json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Printf("%q:\n %s\n", err, bodyText)
			return
		}

		for _, param := range result.Results {
			insertOverrideP(i.ID, param)
		}

	}
}
func InsertPuppetSmartClasses(host string) {
	puppetClasses := getAllPuppetClasses(host)
	puppetClassesCount := getCountAllPuppetClasses(host)
	fmt.Println(host, puppetClassesCount)
	for _, pClass := range puppetClasses {
		var result entitys.PuppetClassName
		fmt.Println(" puppetclasses/" + pClass)
		bodyText := getForemanAPI(host, "puppetclasses/"+pClass+"")
		err := json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Printf("%q:\n %s\n", err, bodyText)
			return
		}

		for _, sc := range result.SmartClassParameters {
			insSmartClasses(host, pClass, sc.ID, sc.Parameter)
		}

	}
}
