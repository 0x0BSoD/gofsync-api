package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

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
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
