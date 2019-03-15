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
						 	                         "locList" text,
						 	                         "created_at" datetime NOT NULL, 
						 	                         "updated_at" datetime NOT NULL);

			CREATE INDEX "index_hg_on_name" ON "hg" ("name");
			CREATE INDEX "index_hg_on_host" ON "hg" ("host");

			CREATE TABLE IF NOT EXISTS "locations" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
													"host" varchar NOT NULL, 
							  						"loc" text NOT NULL);
			CREATE INDEX "index_locations_on_host" ON "locations" ("host");

			CREATE TABLE IF NOT EXISTS "environments" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
													"host" varchar NOT NULL, 
							  						"env" text NOT NULL);
			CREATE INDEX "index_environments_on_host" ON "environments" ("host");

			CREATE TABLE IF NOT EXISTS hg_parameters("id" INTEGER NOT NULL CONSTRAINT hg_parameters_pk PRIMARY KEY AUTOINCREMENT,
                                                     "hg_id" INTEGER,
                                                     "name"     TEXT,
                                                     "value"    TEXT,
                                                     "priority" INTEGER);
            CREATE UNIQUE INDEX hg_parameters_id_uindex ON hg_parameters (id);

            CREATE TABLE puppet_classes("id" INTEGER NOT NULL CONSTRAINT puppet_classes_pk PRIMARY KEY AUTOINCREMENT,
            							"host" TEXT,
            	                        "class" TEXT,
            	                        "subclass" TEXT,
            	                        "sc_ids" TEXT,
            	                        "env_ids" TEXT,
            	                        "hg_ids" TEXT
            	                        
            );
            CREATE UNIQUE INDEX puppet_classes_id_uindex on puppet_classes (id);

            CREATE TABLE smart_classes("id" INTEGER NOT NULL CONSTRAINT smart_classes_pk PRIMARY KEY AUTOINCREMENT,
            						   "host" TEXT,
            	                       "parameter" TEXT,
            	                       "parameter_type" TEXT,
            	                       "override_values_count" INTEGER,
            	                       "foreman_id" INTEGER,
            	                       "dump" TEXT
            );
            CREATE UNIQUE INDEX smart_classes_id_uindex on smart_classes (id);

            CREATE TABLE override_values("id" INTEGER NOT NULL CONSTRAINT override_values_pk PRIMARY KEY AUTOINCREMENT,
            	                       "match" TEXT,
            	                       "value" TEXT,
            	                       "sc_id" INTEGER,
            	                       "use_puppet_default" TEXT
            );
            CREATE UNIQUE INDEX override_values_id_uindex on override_values (id);
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
