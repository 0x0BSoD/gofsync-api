package API

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	All(host string, ctx *user.GlobalCTX) (map[string][]PuppetClass, error)
	ByHostGroupID(host string, hgID int, bdId int, ctx *user.GlobalCTX) (map[string][]PuppetClass, error)
}

type Get struct{}
