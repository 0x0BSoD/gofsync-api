package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func dbActions() {

	if _, err := os.Stat("./gofSync.db"); os.IsNotExist(err) {

		db := getDBConn()
		defer db.Close()

		sqlStmt := `
			CREATE TABLE IF NOT EXISTS "hg" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
						 	                         "name" varchar NOT NULL, 
						 	                         "host" varchar NOT NULL, 
						 	                         "dump" text NOT NULL,
						 	                         "pcList" text,
						 	                         "created_at" datetime NOT NULL, 
						 	                         "updated_at" datetime NOT NULL);

			CREATE INDEX "index_hg_on_name" ON "hg" ("name");
			CREATE INDEX "index_hg_on_host" ON "hg" ("host");

			CREATE TABLE IF NOT EXISTS "locations" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
													"host" varchar NOT NULL, 
							  						"list" text NOT NULL, 
							  						"created_at" datetime NOT NULL, 
							  						"updated_at" datetime NOT NULL);
			CREATE INDEX "index_locations_on_host" ON "locations" ("host");

			CREATE TABLE IF NOT EXISTS hg_parameters("id" INTEGER NOT NULL CONSTRAINT hg_parameters_pk PRIMARY KEY AUTOINCREMENT,
                                                     "hg_id" INTEGER,
                                                     "name"     TEXT,
                                                     "value"    TEXT,
                                                     "priority" INTEGER);
            CREATE UNIQUE INDEX hg_parameters_id_uindex ON hg_parameters (id);

            CREATE TABLE puppet_classes("id" INTEGER NOT NULL CONSTRAINT puppet_classes_pk PRIMARY KEY AUTOINCREMENT,
            							"host" TEXT,
            	                        "class" TEXT,
            	                        "subclass" TEXT
            );
            CREATE UNIQUE INDEX puppet_classes_id_uindex on puppet_classes (id);
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
