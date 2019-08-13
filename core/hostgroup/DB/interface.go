package DB

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	ID(name, host string, ctx *user.GlobalCTX) int
	ParameterID(hgId int, name string, ctx *user.GlobalCTX) int
	ForemanID(name, host string, ctx *user.GlobalCTX) int
	ForemanIDs(host string, ctx *user.GlobalCTX) []int
	ByID(ID int, ctx *user.GlobalCTX) (HostGroupJSON, error)
	ByName(host, name string, ctx *user.GlobalCTX) (HostGroupJSON, error)
	List(ctx *user.GlobalCTX) []string
	ListByHost(host string, ctx *user.GlobalCTX) []HostGroupJSON
	Parameters(hgID int, ctx *user.GlobalCTX) []HostGroupParameter
	HostEnvironment(host string, ctx *user.GlobalCTX) string
}

type IMethods interface {
	Add(name, host, data, sweStatus string, foremanId int, ctx *user.GlobalCTX) int
	Parameter(hgID int, parameter HostGroupParameter, ctx *user.GlobalCTX)
}

type DMethods interface {
	ByID(foremanID int, host string, ctx *user.GlobalCTX)
}

type Get struct{}
type Insert struct{}
type Delete struct{}
