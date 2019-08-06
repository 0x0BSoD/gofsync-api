package DB

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	ID(host string, pc string, parameter string, ctx *user.GlobalCTX) int
	IDByForemanID(host string, foremanID int, ctx *user.GlobalCTX) int
	ByID(scID int, ctx *user.GlobalCTX) SmartClass
}

type Get struct{}
