package middleware

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
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

func GetConfig(r *http.Request) models.Session {
	if rv := context.Get(r, ConfigKey); rv != nil {
		return rv.(models.Session)
	}
	fmt.Println("Error on getting sessions")
	return models.Session{}
}

//func GetUser(r *http.Request) string {
//	if rv := context.Get(r, ConfigKey); rv != nil {
//		return rv.(*models.Config).
//	}
//	return ""
//}
