package API

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	All(host string, ctx *user.GlobalCTX) (Environments, error)
	SmartProxyID(host string, ctx *user.GlobalCTX) int
}

type Get struct{}
