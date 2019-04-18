package utils

import (
	"database/sql"
	"fmt"
	cfg "git.ringcentral.com/alexander.simonov/goFsync/models"
	"time"
)

func InitializeDB(cfg *cfg.Config) {
	connectionString := fmt.Sprintf("%s:%s@/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName)
	var err error
	fmt.Println(connectionString)
	cfg.Database.DB, err = sql.Open("mysql", connectionString)
	cfg.Database.DB.SetMaxIdleConns(140)
	cfg.Database.DB.SetMaxOpenConns(100)
	cfg.Database.DB.SetConnMaxLifetime(time.Second * 10)
	if err != nil {
		Error.Fatal(err)
	}
}
