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
	//ByHost(host string, ctx *user.GlobalCTX) []HostGroupJSON
}

type Get struct{}
