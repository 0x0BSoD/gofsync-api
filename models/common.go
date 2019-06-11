package models

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

type Config struct {
	Hosts      []string
	MasterHost string
	Api        struct {
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
		Port      int
		JWTSecret string
		RunSocket bool
		Socket    *websocket.Conn
		Logged    bool
		Redis     redis.Conn
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
	AMQP struct {
		Host     string
		Username string
		Password string
		Port     int
		Channel  *amqp.Channel
		Queue    *amqp.Queue
	}
}

// Response from API wrapper
type Response struct {
	StatusCode int
	Body       []byte
	RequestUri string
}

type Step struct {
	State   string `json:"state,omitempty"`
	Actions string `json:"actions,omitempty"`
	Host    string `json:"host,omitempty"`
	Counter int    `json:"counter,omitempty"`
	Total   int    `json:"total,omitempty"`
}
