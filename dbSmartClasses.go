package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// ======================================================
// CHECKS
// ======================================================
func checkSC(pc string, parameter string, host string) int64 {

	stmt, err := globConf.DB.Prepare("select id from smart_classes where host=? and parameter=? and puppetclass=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(host, parameter, pc).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func getSC(host string, puppetClass string, parameter string) SCGetResAdv {

	stmt, err := globConf.DB.Prepare("select id, override_values_count, foreman_id from smart_classes where parameter=? and puppetclass=? and host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	var foremanId int
	var ovrCount int
	err = stmt.QueryRow(parameter, puppetClass, host).Scan(&id, &ovrCount, &foremanId)
	if err != nil {
		return SCGetResAdv{}
	}

	return SCGetResAdv{
		ID:                  id,
		ForemanId:           foremanId,
		Name:                parameter,
		OverrideValuesCount: ovrCount,
	}
}
func getSCData(scID int) SCGetResAdv {

	stmt, err := globConf.DB.Prepare("select id, parameter, override_values_count, foreman_id from smart_classes where id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	var foremanId int
	var paramName string
	var ovrCount int

	err = stmt.QueryRow(scID).Scan(&id, &paramName, &ovrCount, &foremanId)
	if err != nil {
		return SCGetResAdv{}
	}

	return SCGetResAdv{
		ID:                  id,
		ForemanId:           foremanId,
		Name:                paramName,
		OverrideValuesCount: ovrCount,
	}
}
func getOvrData(scId int, name string, parameter string) []SCOParams {

	stmt, err := globConf.DB.Prepare("select `match`, value, sc_id from override_values where sc_id=? and `match` like ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var results []SCOParams
	matchStr := fmt.Sprintf("hostgroup=SWE/%s", name)

	rows, err := stmt.Query(scId, matchStr)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var match string
		var scdi int
		var val string
		err = rows.Scan(&match, &val, &scdi)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, SCOParams{
			SmartClassId: scdi,
			Parameter:    parameter,
			Match:        match,
			Value:        val,
		})
	}

	return results
}
func getOverridesHG(hgName string) []OvrParams {

	qStr := fmt.Sprintf("hostgroup=SWE/%s", hgName)
	stmt, err := globConf.DB.Prepare("select `match`, value, sc_id from override_values where `match` like ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var results []OvrParams

	rows, err := stmt.Query(qStr)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var smartClassId int
		var value string
		var match string
		err = rows.Scan(&match, &value, &smartClassId)
		if err != nil {
			log.Fatal(err)
		}
		scData := getSCData(smartClassId)
		results = append(results, OvrParams{
			SmartClassName: scData.Name,
			Value:          value,
		})
	}

	return results
}
func getOverridesLoc(locName string) []OvrParams {

	qStr := fmt.Sprintf("location=%s", locName)
	stmt, err := globConf.DB.Prepare("select `match`, value, sc_id from override_values where `match` like ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var results []OvrParams

	rows, err := stmt.Query(qStr)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var smartClassId int
		var value string
		var match string
		err = rows.Scan(&match, &value, &smartClassId)
		if err != nil {
			log.Fatal(err)
		}
		scData := getSCData(smartClassId)
		results = append(results, OvrParams{
			SmartClassName: scData.Name,
			Value:          value,
		})
	}

	return results
}

// ======================================================
// INSERT
// ======================================================
func insertSC(host string, data SCParameter) int64 {

	existID := checkSC(data.PuppetClass.Name, data.Parameter, host)

	if existID == -1 {
		stmt, err := globConf.DB.Prepare("insert into smart_classes(host, puppetclass, parameter, parameter_type, foreman_id, override_values_count, dump) values(?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		sJson, _ := json.Marshal(data)

		res, err := stmt.Exec(host, data.PuppetClass.Name, data.Parameter, data.ParameterType, data.ID, data.OverrideValuesCount, sJson)
		if err != nil {
			log.Fatal(err)
		}

		lastId, _ := res.LastInsertId()
		if data.OverrideValuesCount > 0 {
			return lastId
		} else {
			return -1
		}
	}
	return -1
}

// Insert Smart Class override
func insertSCOverride(scId int64, data OverrideValue, pType string) {

	var strData string

	// Check value type
	if data.Value != nil {
		switch pType {
		case "string":
			strData = data.Value.(string)
		case "array":
			var tmpResInt []string
			var tmpData string
			switch data.Value.(type) {
			case string:
				log.Println("Type Not Match!!")
				log.Println(pType, " || ", data.Value, " || ", data.Match)
				tmpData = data.Value.(string)
			default:
				for _, i := range data.Value.([]interface{}) {
					tmpResInt = append(tmpResInt, i.(string))
				}
				strIng, _ := json.Marshal(tmpResInt)
				tmpData = string(strIng)
			}
			strData = string(tmpData)
		case "boolean":
			strData = strconv.FormatBool(data.Value.(bool))
		case "integer":
			switch data.Value.(type) {
			case string:
				log.Println("Type Not Match!!")
				log.Println(pType, " || ", data.Value, " || ", data.Match)
				strData = data.Value.(string)
			default:
				strData = strconv.FormatFloat(data.Value.(float64), 'f', 6, 64)
			}
		case "hash":
			fmt.Printf("Hash: %T\n", data.Value.(map[string]string))
			log.Fatal("!!! NOT !!!")
		case "real":
			fmt.Printf("Real: %f\n", data.Value.(float64))
			log.Fatal("!!! NOT !!!")
		default:
			log.Fatal("Type not known\b")
		}
	}

	stmt, err := globConf.DB.Prepare("insert into override_values(`match`, value, sc_id, use_puppet_default) values(?,?,?,?)")
	if err != nil {
		fmt.Println(data.Match, strData, scId, data.UsePuppetDefault)
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.Match, strData, scId, data.UsePuppetDefault)
	if err != nil {
		fmt.Println(data.Match, strData, scId, data.UsePuppetDefault)
		log.Fatal(err)
	}

}
