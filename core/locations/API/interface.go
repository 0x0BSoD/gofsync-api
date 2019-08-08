package API

import "git.ringcentral.com/archops/goFsync/core/user"

type GMethods interface {
	All(host string, ctx *user.GlobalCTX) (Locations, error)
}

type Get struct{}
