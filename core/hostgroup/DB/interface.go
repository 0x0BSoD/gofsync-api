package DB

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	ID(name, host string, ctx *user.GlobalCTX) int
	ParameterID(hgId int, name string, ctx *user.GlobalCTX) int
	ForemanID(name, host string, ctx *user.GlobalCTX) int
	ForemanIDs(host string, ctx *user.GlobalCTX) []int
	ByName(host, name string, ctx *user.GlobalCTX) HostGroupJSON
	//All(ctx *user.GlobalCTX) []HostGroupJSON
	//ByHost(host string, ctx *user.GlobalCTX) []HostGroupJSON
}

type Get struct{}
