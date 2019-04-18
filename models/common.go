package models

import "database/sql"

type Config struct {
	Hosts []string
	Api   struct {
		Username   string
		Password   string
		GetPerPage int
	}
	RackTables struct {
		Production string
		Stage      string
	}
	Database struct {
		Provider string
		Username string
		Password string
		DBName   string
		DB       *sql.DB
	}
	Web struct {
		Port      int
		JWTSecret string
	}
	Logging struct {
		TraceLog  string
		ErrorLog  string
		AccessLog string
	}
	LDAP struct {
		BindUser       string
		BindPassword   string
		LdapServer     string
		LdapServerPort int
		BaseDn         string
		MatchStr       string
	}
}

// Response from API wrapper
type Response struct {
	StatusCode int
	Body       []byte
	RequestUri string
}
