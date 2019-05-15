package middleware

import (
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"github.com/gorilla/context"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type key int

const UserKey key = 0
const ConfigKey key = 1

// Chain applies middleware to a http.HandlerFunc
func Chain(f http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for _, m := range middleware {
		f = m(f)
	}
	return f
}

func GetConfig(r *http.Request) *models.Config {
	if rv := context.Get(r, ConfigKey); rv != nil {
		return rv.(*models.Config)
	}
	return nil
}

func GetUser(r *http.Request) string {
	if rv := context.Get(r, ConfigKey); rv != nil {
		return rv.(string)
	}
	return ""
}
