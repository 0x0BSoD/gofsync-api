package user

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// Create the SignIn handler
func SignIn(ctx *GlobalCTX) http.HandlerFunc {
	var jwtKey = []byte(ctx.Config.Web.JWTSecret)
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		// Get the JSON body and decode into credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("401"))
		}

		// Get the expected user
		user, err := LdapGet(creds.Username, creds.Password, ctx)

		// If a password exists for the given user
		// AND, if it is the same as the password we received, the we can move ahead
		// if NOT, then we return an "Unauthorized" status
		if err != nil {
			//utils.GetErrorContext(err)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
		}

		//TODO: set cfg to ctx
		//context.Set(r, 1, cfg)
		// Pass current user creds for API auth
		ctx.Config.Api.Username = creds.Username
		ctx.Config.Api.Password = creds.Password

		// Declare the expiration time of the token
		// here, we have kept it as 24 minutes or 96 hours
		expirationTime := 24
		if creds.RememberMe {
			expirationTime = 96
		}

		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Username:   creds.Username,
			RememberMe: creds.RememberMe,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: time.Now().Add(time.Duration(expirationTime) * time.Hour).Unix(),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		//_, err = cfg.Web.Redis.Do("SETEX", token, string(expirationTime), creds.Username)
		//if err != nil {
		// If there is an error in setting the cache, return an internal server error
		//w.WriteHeader(http.StatusInternalServerError)
		//return
		//}
		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500"))
		}

		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: time.Now().Add(time.Duration(expirationTime) * time.Hour),
			Path:    "/",
		})

		ctx.Set(claims, tokenString)
		_, _ = w.Write([]byte(user))

	}
}

func Refresh(ctx *GlobalCTX) http.HandlerFunc {
	var jwtKey = []byte(ctx.Config.Web.JWTSecret)
	return func(w http.ResponseWriter, r *http.Request) {
		// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		// (END) The code up-till this point is the same as the first part of the `Welcome` route

		// We ensure that a new token is not issued until enough time has elapsed
		// In this case, a new token will only be issued if the old token is within
		// 30 seconds of expiry. Otherwise, return a bad request status
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("400"))
			return
		}

		// Now, create a new token for the current use, with a renewed expiration time
		expirationTime := time.Now().Add(5 * time.Minute)
		if claims.RememberMe {
			expirationTime = time.Now().Add(24 * time.Hour)
		}
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set the new token as the users `token` cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}
}
