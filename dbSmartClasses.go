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
func checkSC(parameter string, host string) int64 {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id from smart_classes where host=? and parameter=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(host, parameter).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}
func checkSCO(scID string) []int64 {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id from override_values where sc_id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(scID)
	if err != nil {
		return []int64{-1}
	}
	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}
	return ids
}

// ======================================================
// GET
// ======================================================
// Return (foreman_ids, sc_ids)
func getSCWithOverrides(host string) []SCGetRes {

	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("select id, foreman_id, parameter_type from smart_classes where host=? and override_values_count != 0")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var results []SCGetRes

	rows, err := stmt.Query(host)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var foremanID int
		var iD int
		var pType string
		err = rows.Scan(&iD, &foremanID, &pType)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, SCGetRes{
			ForemanID: foremanID,
			ID:        iD,
			Type:      pType,
		})
	}
	return results
}

// ======================================================
// INSERT
// ======================================================
func insertSC(host string, data SCParameter) int64 {

	db := getDBConn()
	defer db.Close()

	existID := checkSC(data.Parameter, host)

	if existID == -1 {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into smart_classes(host, parameter, parameter_type, foreman_id, override_values_count, dump) values(?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		sJson, _ := json.Marshal(data)

		res, err := stmt.Exec(host, data.Parameter, data.ParameterType, data.ID, data.OverrideValuesCount, sJson)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()

		lastID, err := res.LastInsertId()
		if data.OverrideValuesCount > 0 {
			return lastID
		} else {
			return -1
		}
	}
	return -1
}

// Insert Smart Class override
func insertSCOverride(scId int64, data OverrideValue, pType string) {

	var strData string

	//if !checkSCO(scId) {
	switch pType {
	case "string":
		//fmt.Printf("Str: %s\n", data.Value.(string))
		strData = data.Value.(string)
	case "array":
		var tmpRes []string
		for _, i := range data.Value.([]interface{}) {
			tmpRes = append(tmpRes, i.(string))
		}
		tmpData, _ := json.Marshal(tmpRes)
		strData = string(tmpData)
		//fmt.Println("Array:", strData)
	case "boolean":
		strData = strconv.FormatBool(data.Value.(bool))
		//fmt.Printf("Bool: %s\n", strData)
	case "integer":
		strData = strconv.FormatFloat(data.Value.(float64), 'f', 6, 64)
		//fmt.Printf("Int: %d\n", strData)
	case "hash":
		fmt.Printf("Hash: %T\n", data.Value.(map[string]string))
		log.Fatal("!!! NOT !!!")
	case "real":
		fmt.Printf("Real: %f\n", data.Value.(float64))
		log.Fatal("!!! NOT !!!")
	default:
		log.Fatal("Type not known\b")
	}

	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into override_values(match, value, sc_id, use_puppet_default) values(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.Match, strData, scId, data.UsePuppetDefault)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	//}
}
