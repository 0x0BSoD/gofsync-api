package utils

import (
	"database/sql"
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func InitializeDB(cfg *models.Config) {
	var connectionString string
	if cfg.Database.Host != "" {
		connectionString = fmt.Sprintf("%s:%s@%s/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.DBName)
	} else {
		connectionString = fmt.Sprintf("%s:%s@/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName)
	}
	var err error
	cfg.Database.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	err = cfg.Database.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Database.DB.SetMaxIdleConns(140)
	cfg.Database.DB.SetMaxOpenConns(100)
	cfg.Database.DB.SetConnMaxLifetime(time.Second * 10)
}
