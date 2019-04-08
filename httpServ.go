package main

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

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

func Token() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
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
			claims := &Claims{}

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
			context.Set(r, UserKey, claims.Username)
			f(w, r)
		}
	}
}

func loggingHandlerPOST(msg string, dataStruct interface{}) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			user := context.Get(r, UserKey)
			if user != nil {
				req, _ := r.GetBody()
				decoder := json.NewDecoder(req)
				err := decoder.Decode(dataStruct)
				jsonStr, _ := json.Marshal(dataStruct)
				if err != nil {
					logger.Error.Fatalf("Error on POST HG Logging!: %s", err)
				}
				logger.Info.Printf("%s tringgered %s DATA: %q", user.(string), msg, jsonStr)
			}
			f(w, r)
		}
	}
}

func loggingHandler(msg string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := context.Get(r, UserKey)
			if user != nil {
				logger.Info.Printf("%s : %s", user.(string), msg)
			}
			f(w, r)
		}
	}
}

// our main function
func Server() {
	router := mux.NewRouter()

	// User ===
	router.HandleFunc("/signin", SignIn).Methods("POST")
	router.HandleFunc("/refreshjwt", Refresh).Methods("POST")

	// GET ===
	router.HandleFunc("/", Chain(Index, loggingHandler("root of API"), Token())).Methods("GET")
	// Hosts
	router.HandleFunc("/hosts", Chain(getAllHostsHttp, Token())).Methods("GET")
	// Env
	router.HandleFunc("/env/{host}", Chain(getAllEnv, Token())).Methods("GET")
	// Locations
	router.HandleFunc("/loc", Chain(getAllLocHttp, Token())).Methods("GET")
	// Host Groups
	router.HandleFunc("/hg", Chain(getAllHGListHttp, Token())).Methods("GET")
	router.HandleFunc("/hg/{host}", Chain(getHGListHttp, Token())).Methods("GET")
	router.HandleFunc("/hg/{host}/{swe_id}", Chain(getHGHttp, Token())).Methods("GET")
	router.HandleFunc("/hg/foreman/update/{host}/{hgName}", Chain(getHGUpdateInBaseHttp, Token())).Methods("GET")
	router.HandleFunc("/hg/foreman/get/{host}/{hgName}", Chain(getHGFHttp, Token())).Methods("GET")
	router.HandleFunc("/hg/foreman/check/{host}/{hgName}", Chain(getHGCheckHttp, Token())).Methods("GET")
	router.HandleFunc("/hg/overrides/{hgName}", Chain(getOverridesByHGHttp, Token())).Methods("GET")
	// Locations
	router.HandleFunc("/loc/overrides/{locName}", Chain(getOverridesByLocHttp, Token())).Methods("GET")

	// POST ===
	var dataStruct HGPost
	router.HandleFunc("/hg/upload", Chain(postHGHttp, loggingHandlerPOST("upload HG", &dataStruct), Token())).Methods("POST")
	router.HandleFunc("/hg/check", Chain(postHGCheckHttp, Token())).Methods("POST")
	router.HandleFunc("/hg/update", Chain(postHGUpdateHttp, loggingHandlerPOST("updated HG data", &dataStruct), Token())).Methods("POST")
	router.HandleFunc("/env/check", Chain(postEnvCheckHttp, Token())).Methods("POST")

	// Run Server
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://sjc01-c01-pds10:8086", "http://localhost:8080"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, handler)))
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "nope")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
