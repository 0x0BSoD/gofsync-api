package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/logger"
	"strconv"
)

// ======================================================
// CHECKS
// ======================================================
func checkSC(pc string, parameter string, host string) int64 {

	var id int64
	//fmt.Printf("select id from smart_classes where host=%s and parameter=%s and puppetclass=%s\n", host, parameter, pc)
	stmt, err := globConf.DB.Prepare("select id from smart_classes where host=? and parameter=? and puppetclass=?")
	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(host, parameter, pc).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}
func checkOvr(scId int64, match string) int64 {

	var id int64
	//fmt.Printf("select id from override_values where sc_id=%d and `match`=%s\n", scId, match)
	stmt, err := globConf.DB.Prepare("select id from override_values where sc_id=? and `match`=?")
	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(scId, match).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func getSC(host string, puppetClass string, parameter string) SCGetResAdv {

	var id int
	var foremanId int
	var ovrCount int

	stmt, err := globConf.DB.Prepare("select id, override_values_count, foreman_id from smart_classes where parameter=? and puppetclass=? and host=?")
	if err != nil {
		logger.Warning.Printf("%q, checkSC", err)
	}
	defer stmt.Close()

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
		logger.Warning.Printf("%q, getSCData", err)
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

	var results []SCOParams
	matchStr := fmt.Sprintf("hostgroup=SWE/%s", name)

	stmt, err := globConf.DB.Prepare("select foreman_id, `match`, value, sc_id from override_values where sc_id=? and `match` like ?")
	if err != nil {
		logger.Warning.Printf("%q, getOvrData", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(scId, matchStr)
	if err != nil {
		logger.Warning.Printf("%q, getOvrData", err)
	}

	for rows.Next() {
		var foremanId int
		var match string
		var scID int
		var val string
		err = rows.Scan(&foremanId, &match, &val, &scID)
		if err != nil {
			logger.Warning.Printf("%q, getOvrData", err)
		}
		results = append(results, SCOParams{
			OverrideId:   foremanId,
			SmartClassId: scID,
			Parameter:    parameter,
			Match:        match,
			Value:        val,
		})
	}

	return results
}
func getOverridesHG(hgName string) []OvrParams {

	var results []OvrParams

	qStr := fmt.Sprintf("hostgroup=SWE/%s", hgName)
	stmt, err := globConf.DB.Prepare("select `match`, value, sc_id from override_values where `match` like ?")
	if err != nil {
		logger.Warning.Printf("%q, getOverridesHG", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(qStr)
	if err != nil {
		logger.Warning.Printf("%q, getOverridesHG", err)
	}
	for rows.Next() {
		var smartClassId int
		var value string
		var match string
		err = rows.Scan(&match, &value, &smartClassId)
		if err != nil {
			logger.Warning.Printf("%q, getOverridesHG", err)
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

	var results []OvrParams

	qStr := fmt.Sprintf("location=%s", locName)
	stmt, err := globConf.DB.Prepare("select `match`, value, sc_id from override_values where `match` like ?")
	if err != nil {
		logger.Warning.Printf("%q, getOverridesLoc", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(qStr)
	if err != nil {
		logger.Warning.Printf("%q, getOverridesLoc", err)
	}
	for rows.Next() {
		var smartClassId int
		var value string
		var match string
		err = rows.Scan(&match, &value, &smartClassId)
		if err != nil {
			logger.Warning.Printf("%q, getOverridesLoc", err)
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
			logger.Warning.Printf("%q, insertSC", err)
		}
		defer stmt.Close()

		sJson, _ := json.Marshal(data)
		res, err := stmt.Exec(host, data.PuppetClass.Name, data.Parameter, data.ParameterType, data.ID, data.OverrideValuesCount, sJson)
		if err != nil {
			logger.Warning.Printf("%q, insertSC", err)
		}

		lastId, _ := res.LastInsertId()
		if data.OverrideValuesCount > 0 {
			return lastId
		} else {
			return -1
		}
	} else {
		stmt, err := globConf.DB.Prepare("UPDATE `goFsync`.`smart_classes` SET `override_values_count` = ? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Printf("%q, updateSC", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(data.OverrideValuesCount, existID)
		if err != nil {
			logger.Warning.Printf("%q, updateSC", err)
		}
		if data.OverrideValuesCount > 0 {
			return existID
		} else {
			return -1
		}
	}
}

// Insert Smart Class override
func insertSCOverride(scId int64, data OverrideValue, pType string) {

	var strData string

	// Check value type
	if data.Value != nil {

		fmt.Println(scId, data, pType)

		switch pType {
		case "string":
			strData = data.Value.(string)
		case "array":
			var tmpResInt []string
			var tmpData string
			switch data.Value.(type) {
			case string:
				logger.Warning.Printf("Type Not Match!! Type: %s, Val: %s, Match: %s", pType, data.Value, data.Match)
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
			var tmpData string
			switch data.Value.(type) {
			case string:
				logger.Warning.Printf("Type Not Match!! Type: %s, Val: %s, Match: %s", pType, data.Value, data.Match)
				tmpData = data.Value.(string)
			default:
				tmpData = strconv.FormatBool(data.Value.(bool))
			}
			strData = string(tmpData)
		case "integer":
			switch data.Value.(type) {
			case string:
				logger.Warning.Printf("Type Not Match!! Type: %s, Val: %s, Match: %s", pType, data.Value, data.Match)
				strData = data.Value.(string)
			default:
				strData = strconv.FormatFloat(data.Value.(float64), 'f', 6, 64)
			}
		case "hash":
			logger.Warning.Printf("Type inversion not implemented. Type: %s, Val: %s, Match: %s", pType, data.Value, data.Match)
		case "real":
			logger.Warning.Printf("Type inversion not implemented. Type: %s, Val: %s, Match: %s", pType, data.Value, data.Match)
		default:
			logger.Warning.Printf("Type not known, Type: %s, Val: %s, Match: %s", pType, data.Value, data.Match)
		}
	}
	existId := checkOvr(scId, data.Match)
	if existId == -1 {
		stmt, err := globConf.DB.Prepare("insert into override_values(foreman_id, `match`, value, sc_id, use_puppet_default) values(?, ?,?,?,?)")
		if err != nil {
			logger.Warning.Printf("%q, insertSCOverride", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(data.ID, data.Match, strData, scId, data.UsePuppetDefault)
		if err != nil {
			logger.Warning.Printf("%q, insertSCOverride", err)
		}
	} else {
		stmt, err := globConf.DB.Prepare("UPDATE `goFsync`.`override_values` SET `value` = ?, `foreman_id`=? WHERE (`id` = ?)")
		if err != nil {
			logger.Warning.Printf("%q, Prepare updateSCOverride data: %q, %d", err, strData, existId)
		}
		defer stmt.Close()

		_, err = stmt.Exec(strData, data.ID, existId)
		if err != nil {
			logger.Warning.Printf("%q, Exec updateSCOverride data: %q, %d", err, strData, existId)
		}
	}
}

// ======================================================
// UPDATE
// ======================================================
