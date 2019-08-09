package DB

import (
	"database/sql"
	"git.ringcentral.com/archops/goFsync/core/user"
)

type GMethods interface {
	ID(host string, db *sql.DB) int
	All(ctx *user.GlobalCTX) []ForemanHost
	Environment(host string, ctx *user.GlobalCTX) string
}

type IMethods interface {
	Add(host string, db *sql.DB) int
}

type Get struct{}
type Insert struct{}
