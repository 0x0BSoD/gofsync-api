package main

import (
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
	"log"
	"sort"
	"strconv"
)

// ===============
// GET
// ===============
func getPuppetClasses(host string, count string) {

	var result entitys.PuppetClasses

	bodyText := getAPI(host, "puppetclasses?per_page="+count+"")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	for index, cl := range result.Results {
		fmt.Printf("%s ====\n", index)
		var subclassList []string
		for _, v := range cl {
			subclassList = append(subclassList, v.Name)
		}

		sort.Strings(subclassList)
		fmt.Println(subclassList)

		for _, subclass := range subclassList {
			insertToPupClasses(index, subclass)
		}

	}
}

func getPuppetClassesByHostgroup(host string, hostgroupID int) {

	var result entitys.PuppetClasses

	fmt.Printf("Getting %d class.\n", hostgroupID)

	bodyText := getAPI(host, "hostgroups/"+strconv.Itoa(hostgroupID)+"/puppetclasses")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	for index, cl := range result.Results {
		fmt.Printf("%s ====\n", index)
		var subclassList []string
		for _, v := range cl {
			subclassList = append(subclassList, v.Name)
		}
		fmt.Println(subclassList)
	}
}

// ===============
// INSERT
// ===============
func InsertToOverridesBase(host string) {

	var result entitys.SCPOverride

	aSC := getAllSClasses()

	for _, sc := range aSC {

		bodyText := getAPI(host, "smart_class_parameters/"+strconv.Itoa(sc.SCID)+"")

		err := json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Printf("%q:\n %s\n", err, bodyText)
			return
		}

		fmt.Println("PARAM   : ", result.Parameter)
		fmt.Println("DESC    :", result.Description)
		fmt.Println("OVERR   :", result.Override)
		fmt.Println("PTYPE   :", result.ParameterType)
		fmt.Println("DEFV    : ", result.DefaultValue)
		fmt.Println("USEPDEF : ", result.UsePuppetDefault)
		fmt.Println("REQ     : ", result.Required)
		fmt.Println("VALIDT  : ", result.ValidatorType)
		fmt.Println("VALIDR  : ", result.ValidatorRule)
		fmt.Println("MOVERR  : ", result.MergeOverrides)
		fmt.Println("AVOIDD  : ", result.AvoidDuplicates)
		fmt.Println("OVERRVO : ", result.OverrideValueOrder)
		fmt.Println("OVERRVC : ", result.OverrideValuesCount)

		insertSCOverride(result, sc.SCID)

		fmt.Println()
	}
}
func InsertOverridesParameters(host string) {
	Params := getOverrideAllParamBase()
	for _, i := range Params {
		var result entitys.OverrideValuesContainer

		fmt.Println(i)

		bodyText := getAPI(host, "smart_class_parameters/"+strconv.Itoa(i.ClassID)+"/override_values")
		err := json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Printf("%q:\n %s\n", err, bodyText)
			return
		}

		for _, param := range result.Results {
			fmt.Println(param.Match)
			insertOverrideP(i.ID, param)
		}

	}
}
func InsertPuppetSmartClasses(host string) {
	puppetClasses := getAllPuppetClasses()
	for _, pClass := range puppetClasses {
		var result entitys.PuppetClassName

		bodyText := getAPI(host, "puppetclasses/"+pClass+"")

		err := json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Printf("%q:\n %s\n", err, bodyText)
			return
		}

		fmt.Println("Name  :  ", result.Name)

		for _, sc := range result.SmartClassParameters {
			fmt.Println(" SmartClassParameter :  ", sc.Parameter)
			insSmartClasses(host, pClass, sc.ID, sc.Parameter)
		}
		fmt.Println()
	}
}
