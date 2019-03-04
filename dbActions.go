package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/foremanGetter/entitys"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

// ======================================================
// HELPERS
// ======================================================
func getDBConn() *sql.DB {

	db, err := sql.Open("sqlite3", Config.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// ======================================================
// CHECKS
// ======================================================
func checkSWE(name string, host string, db *sql.DB) bool {

	stmt, err := db.Prepare("select id from swes where name=? and host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func checkPC(subclass string, host string, db *sql.DB) bool {

	stmt, err := db.Prepare("select id from puppet_classes where subclass=? and host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(subclass, host).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func checkSCInsert(class string, param string, host string, db *sql.DB) bool {

	stmt, err := db.Prepare("select id from sc_params where param=? and host=? and class=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(param, host, class).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func checkSWEInStateTable(name string, db *sql.DB) bool {

	stmt, err := db.Prepare("select id from swes_state where swe_name=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func SWEstate(host string, swe string) bool {
	db := getDBConn()
	defer db.Close()

	exist := checkSWE(swe, host, db)
	return exist
}

// ======================================================
// GET
// ======================================================
func getAllSWE() []string {

	db := getDBConn()
	defer db.Close()

	rows, err := db.Query("select name from swes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var swes []string

	for rows.Next() {
		var swe string
		err = rows.Scan(&swe)
		if err != nil {
			log.Fatal(err)
		}
		if !stringInSlice(swe, swes) {
			swes = append(swes, swe)
		}
	}
	return swes
}
func getAllSClasses(host string) []entitys.Result {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id, class, param, id_in_puppethost from sc_params where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var res []entitys.Result
	rows, err := stmt.Query(host)
	for rows.Next() {

		var classId int
		var className string
		var puppetSCOverrides string
		var puppetSCOverridesID int

		err = rows.Scan(&classId, &className, &puppetSCOverrides, &puppetSCOverridesID)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, entitys.Result{
			ClassID:           classId,
			ClassName:         className,
			PuppetSCOverrides: puppetSCOverrides,
			SCID:              puppetSCOverridesID,
		})
	}
	return res
}
func getCountAllPuppetClasses(host string) int {
	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select COUNT(*) from puppet_classes where host=? group by host")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var res int

	err = stmt.QueryRow(host).Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	//for _, class := range puppetClass {
	//	classes = append(classes, class)
	//}
	return res
}
func getAllPuppetClasses(host string) []string {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select subclass from puppet_classes where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var classes []string

	rows, err := stmt.Query(host)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var puppetClass string
		err = rows.Scan(&puppetClass)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(puppetClass)
		classes = append(classes, puppetClass)
	}
	//for _, class := range puppetClass {
	//	classes = append(classes, class)
	//}
	return classes
}

func getOverrideAllParamBase(host string) []entitys.IDs {

	db := getDBConn()
	defer db.Close()

	stmt, err := db.Prepare("select id, class_id from over_params_base where host=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var res []entitys.IDs
	rows, err := stmt.Query(host)

	for rows.Next() {

		var opID int
		var cID int
		err = rows.Scan(&opID, &cID)
		if err != nil {
			log.Fatal(err)
		}

		res = append(res, entitys.IDs{
			ID:      opID,
			ClassID: cID,
		})
	}
	return res
}

// ======================================================
// INSERT
// ======================================================
func insSmartClasses(host string, class string, parID int, params string) {

	db := getDBConn()
	defer db.Close()

	if !checkSCInsert(class, params, host, db) {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into sc_params(host, class, id_in_puppethost, param) values(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(host, class, parID, params)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
	}
}

func insertToSWE(name string, host string, data string) bool {

	db := getDBConn()
	defer db.Close()

	if ! checkSWE(name, host, db) {
		fmt.Println(host, " == ", name)

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into swes(name, host, dump, created_at, updated_at) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, host, data, time.Now(), time.Now())
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
	}
	//exist := checkSWE(name, host, db)
	//
	//if exist {
	//	return false
	//} else {
	//	return true
	//}
	return true
}

func insertToLocations(host string, list string) bool {

	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into locations(host, list, created_at, updated_at) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(host, list, time.Now(), time.Now())
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	return true
}

func insertSWEs(swe string) {

	db := getDBConn()
	defer db.Close()
	//if checkSWEInStateTable(swe, db) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into swes_state(swe_name, check_date) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(swe, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
	//} else {
	//	log.Printf("Nope! %s exist", swe)
	//}

}

func insertSCOverride(host string, data entitys.SCPOverride, classId int) {

	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into over_params_base(host, parameter, description, override, parameter_type, default_value, use_puppet_default, required, validator_type, validator_rule, merge_overrides, avoid_duplicates, override_value_order, override_values_count, class_id) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	mDef, _ := json.Marshal(data.DefaultValue)
	mOrd, _ := json.Marshal(data.OverrideValueOrder)

	_, err = stmt.Exec(host, data.Parameter, data.Description,
		data.Override, data.ParameterType,
		mDef, data.UsePuppetDefault,
		data.Required, data.ValidatorType,
		data.ValidatorRule, data.MergeOverrides,
		data.AvoidDuplicates, mOrd,
		data.OverrideValuesCount, classId)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func insertOverrideP(baseId int, data *entitys.OverrideValues) {
	db := getDBConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into override_params(over_params_base_id, match, value, use_puppet_default) values(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	mVal, _ := json.Marshal(data.Value)
	_, err = stmt.Exec(baseId, data.Match, mVal, data.UsePuppetDefault)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func insertSWEState(host string, swe string, state string) {

	db := getDBConn()
	defer db.Close()

	q := fmt.Sprintf("update swes_state set `%s` = '%s' where swe_name = '%s'", host, state, swe)
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func insertToPupClasses(host string, class string, list string) bool {

	db := getDBConn()
	defer db.Close()

	if !checkPC(host, list, db) {
		fmt.Println(host, " Class: ", list)
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into puppet_classes(host, class, subclass) values(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(host, class, list)
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()
	}
	return true
}

// =================================================
// CREATIONS
// =================================================
func dbActions() {

	if _, err := os.Stat("./gofSync.db"); os.IsNotExist(err) {

		db := getDBConn()
		defer db.Close()

		sqlStmt := `CREATE TABLE IF NOT EXISTS "swes" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
						 "name" varchar NOT NULL, 
						 "host" varchar NOT NULL, 
						 "dump" text NOT NULL, 
						 "created_at" datetime NOT NULL, 
						 "updated_at" datetime NOT NULL);

CREATE INDEX "index_swes_on_name" ON "swes" ("name");
CREATE INDEX "index_swes_on_host" ON "swes" ("host");

CREATE TABLE IF NOT EXISTS "bts" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
						"name" varchar NOT NULL, 
						"vcenter" varchar NOT NULL, 
						"info" varchar, 
						"present" boolean NOT NULL, 
						"created_at" datetime NOT NULL, 
						"updated_at" datetime NOT NULL);

CREATE INDEX "index_bts_on_name" ON "bts" ("name");
CREATE INDEX "index_bts_on_vcenter" ON "bts" ("vcenter");

CREATE TABLE IF NOT EXISTS "locations" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
							  "host" varchar NOT NULL, 
							  "list" text NOT NULL, 
							  "created_at" datetime NOT NULL, 
							  "updated_at" datetime NOT NULL);

CREATE INDEX "index_locations_on_host" ON "locations" ("host");

CREATE TABLE IF NOT EXISTS "lts" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
						"name" varchar NOT NULL, 
						"vcenter" varchar NOT NULL, 
						"cluster" varchar NOT NULL, 
						"datastore" varchar, 
						"info" varchar, 
						"dc" varchar NOT NULL, 
						"present" boolean NOT NULL, 
						"created_at" datetime NOT NULL, 
						"updated_at" datetime NOT NULL);

CREATE INDEX "index_lts_on_name" ON "lts" ("name");
CREATE INDEX "index_lts_on_vcenter" ON "lts" ("vcenter");

CREATE TABLE IF NOT EXISTS "replications" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
								 "arguments" varchar NOT NULL, 
								 "cat" integer DEFAULT 0, 
								 "status" integer DEFAULT 0, 
								 "result" integer, 
								 "log" text,
								 "basetpl" varchar NOT NULL, 
								 "vcenter" varchar NOT NULL, 
								 "user" varchar NOT NULL, 
								 "created_at" datetime NOT NULL, 
								 "updated_at" datetime NOT NULL);

CREATE TABLE swes_state
(
	id integer not null
		constraint swes_state_pk
			primary key autoincrement,
	swe_name text not null,
	check_date Datetime not null,
	"rt.ringcentral.com" text, 
	"rt.stage.ringcentral.com" text, 
	"spb01-puppet.lab.nordigy.ru" text, 
	"xmn02-puppet.lab.nordigy.ru" text, 
	"sjc01-puppet.ringcentral.com" text, 
	"sjc02-puppet.ringcentral.com" text, 
	"sjc06-puppet.ringcentral.com" text, 
	"sjc10-puppet.ringcentral.com" text, 
	"iad01-puppet.ringcentral.com" text, 
	"ams01-puppet.ringcentral.com" text, 
	"ams03-puppet.ringcentral.com" text, 
	"zrh01-puppet.ringcentral.com" text);

CREATE UNIQUE INDEX swes_state_id_uindex
	on swes_state (id);

CREATE TABLE locations_state
(
	id                             integer  not null
		constraint locations_state_pk
			primary key autoincrement,
	swe_name                       text     not null,
	check_date                     Datetime not null,
	"spb01-puppet.lab.nordigy.ru"  text,
	"xmn02-puppet.lab.nordigy.ru"  text,
	"sjc01-puppet.ringcentral.com" text,
	"sjc02-puppet.ringcentral.com" text,
	"sjc06-puppet.ringcentral.com" text,
	"sjc10-puppet.ringcentral.com" text,
	"iad01-puppet.ringcentral.com" text,
	"ams01-puppet.ringcentral.com" text,
	"ams03-puppet.ringcentral.com" text,
	"zrh01-puppet.ringcentral.com" text);

CREATE UNIQUE INDEX locations_state_id_uindex
	on locations_state (id);

CREATE TABLE puppet_classes
(
	id integer not null
		constraint puppet_classes_pk
			primary key autoincrement,
	host text,
	class text,
	subclass text
);

CREATE UNIQUE INDEX puppet_classes_id_uindex
	on puppet_classes (id);

CREATE TABLE IF NOT EXISTS "sc_params"
(
	id integer not null
		constraint sc_params_pk
			primary key autoincrement,
	host text,
	class text,
	param text,
	id_in_puppethost integer
);

CREATE UNIQUE INDEX sc_params_id_uindex
	on sc_params (id);

CREATE TABLE IF NOT EXISTS "over_params_base"
(
	id integer not null
		constraint over_params_pk
			primary key autoincrement,
	parameter text,
	host text,
	description text,
	override int,
	parameter_type text,
	default_value text,
	use_puppet_default int,
	required int,
	validator_type text,
	validator_rule text,
	merge_overrides int,
	avoid_duplicates int,
	override_value_order text,
	override_values_count int
, class_id int);

CREATE UNIQUE INDEX over_params_id_uindex
	on "over_params_base" (id);

CREATE TABLE override_params
(
	id integer not null
		constraint override_params_pk
			primary key autoincrement,
	over_params_base_id int,
	match text,
	value text
, use_puppet_default int);

CREATE UNIQUE INDEX override_params_id_uindex
	on override_params (id);

/* No STAT tables available */
	`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return
		}
	} else {
		fmt.Println("Base file exist")
	}
}
