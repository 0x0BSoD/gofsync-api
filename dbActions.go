package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

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
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ======================================================
func getAllSWE() []string {

	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
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

func insSmartClasses(host string, class string, parID int, params string) {

	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

func getAllSClasses() []entitys.Result {

	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, class, param, id_in_puppethost from sc_params")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var res []entitys.Result

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

func getAllPuppetClasses() []string {
	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select subclass from puppet_classes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var classes []string

	for rows.Next() {
		var puppetClass string
		err = rows.Scan(&puppetClass)
		if err != nil {
			log.Fatal(err)
		}
		classes = append(classes, puppetClass)
	}
	return classes
}

func insertToSWE(name string, host string, data string) bool {

	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	exist := checkSWE(name, host, db)
	if !exist {
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
	exist = checkSWE(name, host, db)

	if exist {
		return false
	} else {
		return true
	}
}

func insertToLocations(host string, list string) bool {

	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
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
	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
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
}

func insertSCOverride(data entitys.SCPOverrideForBase) {
	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into over_params(name, override, type, default_value, override_value_order, override_values, class_id) values(?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	ovJson, _ := json.Marshal(data.OverrideValues)
	defJson, _ := json.Marshal(data.DefaultValue)
	_, err = stmt.Exec(data.Name, data.Override, data.ValidatorType, defJson, data.OverrideValueOrder, ovJson, data.ClassID)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func insertSWEState(host string, swe string, state string) {
	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
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

func SWEstate(host string, swe string) bool {

	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	exist := checkSWE(swe, host, db)
	return exist
}

func insertToPupClasses(class string, list string) bool {

	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into puppet_classes(class, subclass) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(class, list)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	return true
}

func createSmartClassBase(host string) {
	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := fmt.Sprintf(`
		CREATE TABLE "smart_class_%s" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
						 "name" varchar NOT NULL, 
						 "host" varchar NOT NULL, 
						 "dump" text NOT NULL, 
						 "created_at" datetime NOT NULL, 
						 "updated_at" datetime NOT NULL);
	CREATE INDEX "index_swes_on_name" ON "swes" ("name");
	CREATE INDEX "index_swes_on_host" ON "swes" ("host");	
	`, host)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func dbActions() {

	if _, err := os.Stat("./gofSync.db"); os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", "./gofSync.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		sqlStmt := `
	CREATE TABLE "swes" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
						 "name" varchar NOT NULL, 
						 "host" varchar NOT NULL, 
						 "dump" text NOT NULL, 
						 "created_at" datetime NOT NULL, 
						 "updated_at" datetime NOT NULL);
	CREATE INDEX "index_swes_on_name" ON "swes" ("name");
	CREATE INDEX "index_swes_on_host" ON "swes" ("host");	

	CREATE TABLE "bts" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
						"name" varchar NOT NULL, 
						"vcenter" varchar NOT NULL, 
						"info" varchar, 
						"present" boolean NOT NULL, 
						"created_at" datetime NOT NULL, 
						"updated_at" datetime NOT NULL);
    CREATE INDEX "index_bts_on_name" ON "bts" ("name");
    CREATE INDEX "index_bts_on_vcenter" ON "bts" ("vcenter");

	CREATE TABLE "locations" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
							  "host" varchar NOT NULL, 
							  "list" text NOT NULL, 
							  "created_at" datetime NOT NULL, 
							  "updated_at" datetime NOT NULL);
	CREATE INDEX "index_locations_on_host" ON "locations" ("host");
	
	CREATE TABLE "lts" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
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
	
	CREATE TABLE "replications" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
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

	CREATE TABLE "swes_state" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
							   "id_swe" INTEGER NOT NULL,
							   "hosts" varchar NOT NULL);
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
