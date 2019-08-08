package API

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	All(host string, ctx *user.GlobalCTX) (map[string][]PuppetClass, error)
	ByHostGroupID(host string, hgID int, bdId int, ctx *user.GlobalCTX) (map[string][]PuppetClass, error)
	ByID(host string, pcId int, ctx *user.GlobalCTX) (map[string][]PuppetClass, error)
}

type UMethods interface {
	SmartClassIDs(host string, ctx *user.GlobalCTX)
}

type Get struct{}
type Update struct{}
