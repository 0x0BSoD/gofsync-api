package main

import (
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
	"log"
	"sort"
	"strconv"
	"strings"
)

func getPuppetClassesbyHostgroup(host string, hostgroupID int) {

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
			//fmt.Println("    ID          :  ", v.ID)
			//fmt.Println("    Name        :  ", v.Name)
			subclassList = append(subclassList, v.Name)
			//fmt.Println("    CreatedAt   :  ", v.CreatedAt)
			//fmt.Println("    UpdatedAt   :  ", v.UpdatedAt)
		}
		fmt.Println(subclassList)
	}
}

func getAllPuppetSmartClasses(host string) {
	puppetClasses := getAllPuppetClasses()
	for _, pClass := range puppetClasses {
		getPuppetSmartClasses(host, pClass)
	}
}

func getPuppetSmartClasses(host string, class string) {
	var result entitys.PuppetClassName

	bodyText := getAPI(host, "puppetclasses/"+class+"")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	fmt.Println("Name  :  ", result.Name)

	var params []string

	for _, sc := range result.SmartClassParameters {
		fmt.Println(" SmartClassParameter :  ", sc.Parameter)
		params = append(params, sc.Parameter)
	}
	fmt.Println()
	insSmartClasses(result.Name, host, class, strings.Join(params, ","))
}

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
			//fmt.Println("    ID          :  ", v.ID)
			//fmt.Println("    Name        :  ", v.Name)
			subclassList = append(subclassList, v.Name)
			//fmt.Println("    CreatedAt   :  ", v.CreatedAt)
			//fmt.Println("    UpdatedAt   :  ", v.UpdatedAt)
		}

		sort.Strings(subclassList)
		fmt.Println(subclassList)

		for _, subclass := range subclassList {
			insertToPupClasses(index, subclass)
		}

	}
}
