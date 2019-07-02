package middleware

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"github.com/gorilla/context"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type key int

const ContextKey key = 1

// Chain applies middleware to a http.HandlerFunc
func Chain(f http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for _, m := range middleware {
		f = m(f)
	}
	return f
}

func GetContext(r *http.Request) *user.GlobalCTX {
	if rv := context.Get(r, ContextKey); rv != nil {
		return rv.(*user.GlobalCTX)
	}
	fmt.Println("Error on getting sessions")
	var void *user.GlobalCTX
	return void
}
