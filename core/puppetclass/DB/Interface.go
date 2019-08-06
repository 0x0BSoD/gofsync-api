package DB

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	ID(subclass string, host string, ctx *user.GlobalCTX) int
	All(host string, ctx *user.GlobalCTX) []PuppetClass
	ByName(subclass string, host string, ctx *user.GlobalCTX) PuppetClass
	ByID(pId int, ctx *user.GlobalCTX) PuppetClass
}

type IMethods interface {
	Insert(host string, class string, subclass string, foremanId int, ctx *user.GlobalCTX) int
}

type DMethods interface {
	BySubclass(host string, subClass string, ctx *user.GlobalCTX) error
}

type UMethods interface {
	ByID(host string, parameters Parameters, ctx *user.GlobalCTX) int
	HostGroupIDs(hgId int, pcList []int, ctx *user.GlobalCTX)
}

type Insert struct{}
type Get struct{}
type Delete struct{}
type Update struct{}
