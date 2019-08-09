package API

import (
	"git.ringcentral.com/archops/goFsync/core/user"
	"net/url"
)

type GMethods interface {
	ByHostGroup(hgName string, params url.Values, ctx *user.GlobalCTX) map[string][]Host
}

type Get struct{}
