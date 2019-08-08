package DB

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	ID(host, loc string, ctx *user.GlobalCTX) int
	ForemanIDs(host string, ctx *user.GlobalCTX) []int
	All(host string, ctx *user.GlobalCTX) ([]string, string)
}

type IMethods interface {
	Add(host, loc string, foremanId int, ctx *user.GlobalCTX)
}

type DMethods interface {
	ByName(host, loc string, ctx *user.GlobalCTX)
}

type Get struct{}
type Insert struct{}
type Delete struct{}
