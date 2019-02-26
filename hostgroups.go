package main

import (
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
	"log"
)

func pPrintCommit(result entitys.SWEs, commit bool, host string) {
	for _, item := range result.Results {
		fmt.Println(host + "  ==================================================")
		fmt.Println("ID              :  ", item.ID)
		fmt.Println("Name            :  ", item.Name)
		fmt.Println("Title           :  ", item.Title)
		fmt.Println("EnvironmentName :  ", item.EnvironmentName)
		//getPuppetClasses(host, item.ID)

		sJson, _ := json.Marshal(item)

		if commit {
			if insertToSWE(item.Name, host, string(sJson)) {
				fmt.Println("  ==================================================")
				fmt.Println(item.Name + "  INSERTED")
				fmt.Println("  ==================================================")
			}
		}

		fmt.Println()
	}
}

func getHostGroups(host string, count string) {

	var result entitys.SWEs
	fmt.Printf("Getting from %s \n", host)

	bodyText := getAPI(host, "hostgroups?format=json&per_page="+count+"&search=label+~+SWE%2F")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	pPrintCommit(result, true, host)

}
