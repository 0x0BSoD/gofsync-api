package main

import (
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
	"log"
	"strconv"
)

func getPuppetClasses(host string, classID int) {
	//spaces := 10
	var result entitys.PuppetClasses

	fmt.Printf("Getting %d class.\n", classID)

	bodyText := getAPI(host, "hostgroups/"+strconv.Itoa(classID)+"/puppetclasses")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	for index, cl := range result.Results {
		fmt.Printf("%s ====\n", index)
		for _, v := range cl {
			fmt.Println("    ID          :  ", v.ID)
			fmt.Println("    Name        :  ", v.Name)
			fmt.Println("    CreatedAt   :  ", v.CreatedAt)
			fmt.Println("    UpdatedAt   :  ", v.UpdatedAt)
		}
	}
}
