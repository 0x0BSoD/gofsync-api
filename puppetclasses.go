package main

import (
	"encoding/json"
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

	aSC := getAllSClasses()

	for _, sc := range aSC {

		bodyText := getForemanAPI(host, "smart_class_parameters/"+strconv.Itoa(sc.SCID)+"")

		err := json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Printf("%q:\n %s\n", err, bodyText)
			return
		}

		insertSCOverride(result, sc.SCID)

	}
}
func InsertOverridesParameters(host string) {
	Params := getOverrideAllParamBase()
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
	for _, pClass := range puppetClasses {
		var result entitys.PuppetClassName

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
