package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/entitys"
	"log"
)

// Params:
// If only OK - sync all with OK state
// PROD -  sync only with prod state
// STAGE -  sync only with stage state
//  ["STAGE", "PROD"] - only if have prod and stage state
// ALL - sync if even state equal NOTINRT
// i.g. ["OK"] or ["PROD"] or ["STAGE", "PROD"] or ["ALL"]
func runSync(source string, targets []string, params []string, swes []string) {

	if len(swes) > 0 {
		fmt.Println(source)
		fmt.Println(targets)
		fmt.Println(params)
		fmt.Println(swes)
	} else {
		var strParams string

		if stringInSlice("OK", params) {
			strParams = "OK%"
		} else if stringInSlice("STAGE", params) && stringInSlice("PROD", params) {
			strParams = "OK_PROD_STAGE"
		} else if stringInSlice("STAGE", params) && !stringInSlice("PROD", params) {
			strParams = "OK_STAGE"
		} else if !stringInSlice("STAGE", params) && stringInSlice("PROD", params) {
			strParams = "OK_PROD"
		} else if stringInSlice("STAGE", params) {
			strParams = "%"
		} else {
			log.Fatal("Error in params")
		}

		//var SourceSWEs []entitys.SWE

		strsSWEs := getSWEsByState(source, strParams)

		for _, strSWE := range strsSWEs {
			var SWE entitys.SWE
			err := json.Unmarshal([]byte(strSWE), &SWE)
			if err != nil {
				log.Fatalf("%q:\n %s\n", err, strSWE)
			}
			fmt.Println(SWE.Name)
			fmt.Println(SWE.Ancestry)
			//SourceSWEs = append(SourceSWEs, SWE)
			//}

			//for _, target := range tragets {
			//
			//	var TargetSWEs []entitys.SWE
			//
			//	strsSWEs := getAllSWEByHost(target)
			//	for _, strSWE := range strsSWEs{
			//		err := json.Unmarshal([]byte(strSWE), &SWE)
			//		if err != nil {
			//			log.Fatalf("%q:\n %s\n", err, strSWE)
			//		}
			//		TargetSWEs = append(TargetSWEs, SWE)
			//	}
			//
			//}
		}
	}

}
