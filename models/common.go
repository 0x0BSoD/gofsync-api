package models

import (
	"database/sql"
	"log"
)

type Config struct {
	Hosts      map[string]int
	MasterHost string
	Git        struct {
		Repo      string
		Directory string
		Token     string
	}
	Api struct {
		Username   string
		Password   string
		GetPerPage int
	}
	RackTables struct {
		Production string
		Stage      string
	}
	Database struct {
		Host     string
		Provider string
		Username string
		Password string
		DBName   string
		DB       *sql.DB
	}
	Web struct {
		JWTSecret string
		Port      int
	}
	Logging struct {
		TraceLog  string
		ErrorLog  string
		AccessLog string
	}
	LDAP struct {
		BindUser       string
		BindPassword   string
		LdapServer     []string
		LdapServerPort int
		BaseDn         string
		MatchStr       string
	}
}

func (c *Config) deferCloseStmt(conn *sql.Stmt) {
	if err := conn.Close(); err != nil {
		log.Println("error on closing DB connection: ", err)
	}
}

// DBGetOne
func (c *Config) DBGetOne(query string, resultCallback interface{}, params ...interface{}) error {

	stmt, err := c.Database.DB.Prepare(query)
	if err != nil {
		log.Printf("[Q] %s\n%q", query, err)
	}
	defer c.deferCloseStmt(stmt)

	err = stmt.QueryRow(params...).Scan(resultCallback)
	if err != nil {
		return err
	}

	return nil
}

// Response from API wrapper
type Response struct {
	StatusCode int
	Body       []byte
	RequestUri string
}

type Step struct {
	Status  string      `json:"status,omitempty"`
	State   string      `json:"state,omitempty"`
	Item    string      `json:"item,omitempty"`
	Actions string      `json:"actions,omitempty"`
	Host    string      `json:"host,omitempty"`
	Counter interface{} `json:"counter,omitempty"`
	Total   int         `json:"total,omitempty"`
}

type WSMessage struct {
	Broadcast bool        `json:"broadcast"`
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}
