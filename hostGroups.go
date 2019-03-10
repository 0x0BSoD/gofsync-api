package main

import (
	"encoding/json"
	"fmt"
	"log"
)
// ===============================
// TYPES & VARS
// ===============================
// For Getting SWE from RackTables
type RTSWE struct {
	Name string `json:"name"`
	BaseTpl string `json:"basetpl"`
	OsVersion string `json:"osversion"`
	SWEStatus string `json:"swestatus"`
}
type RTSWES []RTSWE

// ===============================
// METHODS
// ===============================
// Get SWE from RackTables
func (swe RTSWE) Get(host string) RTSWES {
	var r RTSWES
	body := RTAPI("GET", host,
		          "api/rchwswelookups/search?q=name~.*&fields=name,osversion,basetpl,swestatus&format=json")

	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("%q:\n %s\n", err, body)
		return []RTSWE{}
	}
	return r
}
// Return as JSON str
func (swes RTSWES) ToJSON() string {
	sJson, _ := json.Marshal(swes)
	return string(sJson)
}
func (swes RTSWES) String() {
	for _,i := range swes {
		fmt.Println("Name: ", i.Name)
		fmt.Println("Name: ", i.BaseTpl)
		fmt.Println("Name: ", i.OsVersion)
		fmt.Println("Name: ", i.SWEStatus)
		fmt.Println()
	}
}
//func getHostGroups(host string, count string) {
//	var data entitys.SWEs
//
//	bodyText := getForemanAPI(host, "hostgroups?format=json&per_page="+count+"&search=label+~+SWE%2F")
//
//	err := json.Unmarshal(bodyText, &data)
//	if err != nil {
//		log.Fatalf("%q:\n %s\n", err, bodyText)
//	}
//
//	for _, i := range data.Results {
//		sJson, _ := json.Marshal(i)
//		lastId := insertToSWE(i.Name, host, string(sJson))
//		if lastId != -1 {
//			getSWEParams(host, count, lastId, i.ID)
//		}
//	}
//
//	//return  pPrintCommit(result, host)
//}
//
//func getSWEParams(host string, count string, dbID int64, sweID int) {
//	var data entitys.SWEParameterContainer
//
//	bodyText := getForemanAPI(host, "hostgroups/"+strconv.Itoa(sweID)+"/parameters?format=json&per_page="+count)
//
//	err := json.Unmarshal(bodyText, &data)
//	if err != nil {
//		log.Fatalf("%q:\n %s\n", err, bodyText)
//	}
//
//	for _, i := range data.Results {
//		insertToSWEParams(dbID, i.Name, i.Value, i.Priority)
//		fmt.Println(i.Name)
//		fmt.Println(i.Value)
//		fmt.Println()
//	}
//}
