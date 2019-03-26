package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// ======================================================
// CHECKS
// ======================================================
func checkEnv(host string, env string) int {

	stmt, err := globConf.DB.Prepare("select id from environments where host=? and env=?")
	if err != nil {
		log.Fatal(err)
	}

	var id int
	err = stmt.QueryRow(host, env).Scan(&id)
	if err != nil {
		stmt.Close()
		return -1
	}
	stmt.Close()
	return id
}
func checkPostEnv(host string, env string) int {

	stmt, err := globConf.DB.Prepare("select foreman_id from environments where host=? and env=?")
	if err != nil {
		log.Fatal(err)
	}

	var id int
	err = stmt.QueryRow(host, env).Scan(&id)
	if err != nil {
		stmt.Close()
		return -1
	}
	stmt.Close()
	return id
}

// ======================================================
// GET
// ======================================================
func getEnvList(host string) []string {

	stmt, err := globConf.DB.Prepare("select env from environments where host=?")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query(host)
	if err != nil {
		log.Fatal(err)
	}

	var list []string
	for rows.Next() {
		var env string
		err = rows.Scan(&env)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, env)
	}

	stmt.Close()

	return list
}

// ======================================================
// INSERT
// ======================================================
func insertToEnvironments(host string, env string, foremanId int) {

	eId := checkEnv(host, env)
	if eId == -1 {
		stmt, err := globConf.DB.Prepare("insert into environments(host, env, foreman_id) values(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(host, env, foremanId)
		if err != nil {
			log.Fatal(err)
		}
		stmt.Close()
	}
}
