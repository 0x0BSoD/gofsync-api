package middleware

import (
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"net/http"
)

func Token(cfg *models.Config) Middleware {
	var jwtKey = []byte(cfg.Web.JWTSecret)
	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Check base
			err := cfg.Database.DB.Ping()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - DB Error"))
			}

			// Do middleware things
			// We can obtain the session token from the requests cookies, which come with every request
			c, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					// If the cookie is not set, return an unauthorized status
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("401 - Unauthorized"))
					return
				}
				// For any other type of error, return a bad request status
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - BadRequest"))
				return
			}

			// Get the JWT string from the cookie
			tknStr := c.Value

			// Initialize a new instance of `Claims`
			claims := &models.Claims{}

			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - Unauthorized"))
				return
			}
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("401 - Unauthorized"))
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - BadRequest"))
				return
			}
			// Call the next middleware/handler in chain and set user in ctx
			context.Set(r, UserKey, cfg.Api.Username)
			context.Set(r, ConfigKey, user.Start(claims, tknStr, cfg))
			f(w, r)
		}
	}
}
