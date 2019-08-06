package API

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	All(host string, ctx *user.GlobalCTX) ([]Parameter, error)
	OverridesByID(host string, ForemanID int, ctx *user.GlobalCTX) ([]OverrideValue, error)
	ByID(host string, foremanId int, ctx *user.GlobalCTX) (Parameter, error)
	ByPuppetClassID(host string, pcId int, ctx *user.GlobalCTX) []Parameter
}

type Get struct{}
