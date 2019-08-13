package API

import (
	"git.ringcentral.com/archops/goFsync/core/hostgroup/DB"
	"git.ringcentral.com/archops/goFsync/core/user"
)

type GMethods interface {
	All(host string, ctx *user.GlobalCTX) []DB.HostGroup
	Parameters(host string, dbID int, hgID int, ctx *user.GlobalCTX) []DB.HostGroupParameter
}

type Get struct{}
