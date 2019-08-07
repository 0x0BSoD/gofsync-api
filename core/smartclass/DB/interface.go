package DB

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	ID(host string, pc string, parameter string, ctx *user.GlobalCTX) int
	IDByForemanID(host string, foremanID int, ctx *user.GlobalCTX) int
	ByID(scID int, ctx *user.GlobalCTX) (SmartClass, error)
	GetSC(host string, puppetClass string, parameter string, ctx *user.GlobalCTX) (SmartClass, error)
	Override(scID int, name string, parameter string, ctx *user.GlobalCTX) (Override, error)
	OverridesByMatch(host, matchParameter string, ctx *user.GlobalCTX) map[string][]Override
}

type Get struct{}
