package middleware

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc
type key int

const UserKey key = 0

// Chain applies middleware to a http.HandlerFunc
func Chain(f http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for _, m := range middleware {
		f = m(f)
	}
	return f
}
