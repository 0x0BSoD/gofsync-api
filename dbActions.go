package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

func checkSwe(name string, host string, db *sql.DB) bool {

	stmt, err := db.Prepare("select id from swes where name= ? and host = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(name, host).Scan(&id)
	if err != nil {
		return false
	}

	//rows, err := db.Query("select id from swes where name='" + name + "'")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for rows.Next() {
	//	var id int
	//	err = rows.Scan(&id)
	//	if err != nil {
	//		log.Println("select id from swes where name='" + name + "'")
	//		log.Fatal(err)
	//		return false
	//	}
	//	fmt.Printf(string(id))
	//}
	return true
}

func insertToSWE(name string, host string, data string) bool {

	db, err := sql.Open("sqlite3", "./gofSync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	exist := checkSwe(name, host, db)
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
	exist = checkSwe(name, host, db)

	if exist {
		return false
	} else {
		return true
	}
}

func dbActions() {
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
}
