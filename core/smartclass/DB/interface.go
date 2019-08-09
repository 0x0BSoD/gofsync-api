package DB

import (
	"git.ringcentral.com/archops/goFsync/core/smartclass/API"
	"git.ringcentral.com/archops/goFsync/core/user"
)

type GMethods interface {
	ID(host, pc, parameter string, ctx *user.GlobalCTX) int
	IDByForemanID(host string, foremanID int, ctx *user.GlobalCTX) int
	OverrideID(scID, foremanID int, ctx *user.GlobalCTX) int

	ForemanIDs(host string, ctx *user.GlobalCTX) []int
	OverridesFIDsBySmartClassID(scId int, ctx *user.GlobalCTX) []int

	ByID(scID int, ctx *user.GlobalCTX) (SmartClass, error)
	ByParameter(host, puppetClass, parameter string, ctx *user.GlobalCTX) (SmartClass, error)

	Overrides(scID int, ctx *user.GlobalCTX) ([]Override, error)
	OverrideByMatch(scID int, matchParameter string, ctx *user.GlobalCTX) (Override, error)
	OverridesByMatch(host, matchParameter string, ctx *user.GlobalCTX) map[string][]Override
}

type IMethods interface {
	Add(host string, data API.Parameter, ctx *user.GlobalCTX)
	AddOverride(scId int, data OverrideValue, pType string, ctx *user.GlobalCTX)
}

type DMethods interface {
	SmartClass(host string, foremanId int, ctx *user.GlobalCTX) error
	Override(scID int, foremanID int, ctx *user.GlobalCTX) error
}

type Get struct{}
type Insert struct{}
type Delete struct{}
