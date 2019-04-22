package environment

import (
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

// ======================================================
// CHECKS
// ======================================================
func CheckEnv(host string, env string, cfg *models.Config) int {

	var id int

	stmt, err := cfg.Database.DB.Prepare("select id from environments where host=? and env=?")
	if err != nil {
		logger.Warning.Printf("%q, checkEnv", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(host, env).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}
func CheckPostEnv(host string, env string, cfg *models.Config) int {

	var id int

	stmt, err := cfg.Database.DB.Prepare("select foreman_id from environments where host=? and env=?")
	if err != nil {
		logger.Warning.Printf("%q, checkPostEnv", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(host, env).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

// ======================================================
// GET
// ======================================================
func GetEnvList(host string, cfg *models.Config) []string {

	var list []string

	stmt, err := cfg.Database.DB.Prepare("select env from environments where host=?")
	if err != nil {
		logger.Warning.Printf("%q, getEnvList", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(host)
	if err != nil {
		return list
	}

	for rows.Next() {
		var env string
		err = rows.Scan(&env)
		if err != nil {
			logger.Error.Printf("%q, getEnvList", err)
		}
		list = append(list, env)
	}

	return list
}

// ======================================================
// INSERT
// ======================================================
func InsertToEnvironments(host string, env string, foremanId int, cfg *models.Config) {

	eId := CheckEnv(host, env, cfg)
	if eId == -1 {
		stmt, err := cfg.Database.DB.Prepare("insert into environments(host, env, foreman_id) values(?, ?, ?)")
		if err != nil {
			logger.Warning.Printf("%q, insertToEnvironments", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(host, env, foremanId)
		if err != nil {
			logger.Warning.Printf("%q, insertToEnvironments", err)
		}
	}
}
