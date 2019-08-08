package DB

import (
	"git.ringcentral.com/archops/goFsync/core/environment/API"
	"git.ringcentral.com/archops/goFsync/core/user"
)

type GMethods interface {
	ID(host, env string, ctx *user.GlobalCTX) int
	ForemanID(host, env string, ctx *user.GlobalCTX) int

	All(ctx *user.GlobalCTX) map[string][]API.Environment
	ByName(host string, env string, ctx *user.GlobalCTX) API.Environment
	ByHost(host string, ctx *user.GlobalCTX) []string

	Repo(host string, ctx *user.GlobalCTX) string
}

type IMethods interface {
	Add(host, env, state string, foremanId int, codeInfo SvnInfo, ctx *user.GlobalCTX)
}

type UMethods interface {
	SetRepo(repo, host string, ctx *user.GlobalCTX)
	SetState(state, host, name string, ctx *user.GlobalCTX)
}

type DMethods interface {
	ByName(host string, env string, ctx *user.GlobalCTX)
}

type Get struct{}
type Insert struct{}
type Update struct{}
type Delete struct{}
