package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Config struct {
	Actions  []string
	Hosts    []string
	RTPro    string
	RTStage  string
	Username string
	Pass     string
	Port     int
	DBFile   string
	PerPage  int
	DbInit   string
	DB       *sql.DB
}

var (
	webServer bool
	file      string
	conf      string
	host      string
	globConf  Config
)

// =====================
//  DB Init
// =====================
func (a *Config) Initialize(user, password, dbName string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbName)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	a.DB.SetMaxIdleConns(140)
	a.DB.SetMaxOpenConns(100)
	a.DB.SetConnMaxLifetime(time.Second * 10)
	if err != nil {
		log.Fatal(err)
	}
}

// =====================
//  Args
// =====================
func init() {
	flag.StringVar(&conf, "conf", "", "Config file, TOML")
	flag.StringVar(&file, "file", "", "File contain hosts divide by new line")
	flag.StringVar(&host, "host", "", "Foreman FQDN")
	flag.BoolVar(&webServer, "server", false, "Run as web server daemon")
}

func main() {
	flag.Parse()
	configParser()
	getHosts(file)
	if webServer {
		Server()
	} else {
		fullSyncv2()
		//saveHGToJson()
	}
}
