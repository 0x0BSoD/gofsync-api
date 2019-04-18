package main

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/middleware"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

// our main function
func Server(cfg *models.Config) {
	router := mux.NewRouter()

	// User ===
	router.HandleFunc("/signin", SignIn(cfg)).Methods("POST")
	router.HandleFunc("/refreshjwt", Refresh(cfg)).Methods("POST")

	// GET ===
	router.HandleFunc("/", middleware.Chain(Index, middleware.Token(cfg))).Methods("GET")
	// Hosts
	router.HandleFunc("/hosts", middleware.Chain(getAllHostsHttp(cfg), middleware.Token(cfg))).Methods("GET")
	// Env
	router.HandleFunc("/env/{host}", middleware.Chain(getAllEnv(cfg), middleware.Token(cfg))).Methods("GET")
	// Locations
	router.HandleFunc("/loc", middleware.Chain(getAllLocHttp(cfg), middleware.Token(cfg))).Methods("GET")
	// Host Groups
	router.HandleFunc("/hg", middleware.Chain(getAllHGListHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/{host}", middleware.Chain(getHGListHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/{host}/{swe_id}", middleware.Chain(getHGHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/update/{host}/{hgName}", middleware.Chain(getHGUpdateInBaseHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/get/{host}/{hgName}", middleware.Chain(getHGFHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/check/{host}/{hgName}", middleware.Chain(getHGCheckHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/overrides/{hgName}", middleware.Chain(getOverridesByHGHttp(cfg), middleware.Token(cfg))).Methods("GET")
	// Locations
	router.HandleFunc("/loc/overrides/{locName}", middleware.Chain(getOverridesByLocHttp(cfg), middleware.Token(cfg))).Methods("GET")

	// POST ===
	router.HandleFunc("/hg/upload", middleware.Chain(postHGHttp(cfg), middleware.Token(cfg))).Methods("POST")
	router.HandleFunc("/hg/check", middleware.Chain(postHGCheckHttp(cfg), middleware.Token(cfg))).Methods("POST")
	router.HandleFunc("/env/check", middleware.Chain(postEnvCheckHttp(cfg), middleware.Token(cfg))).Methods("POST")

	// Run Server
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://sjc01-c01-pds10:8086", "http://localhost:8080",
			"https://sjc01-c01-pds10:8086", "https://sjc01-c01-pds10.c01.ringcentral.com:8086"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	handler := c.Handler(router)
	bindAddr := fmt.Sprintf(":%d", cfg.Web.Port)
	log.Fatal(http.ListenAndServe(bindAddr, handlers.LoggingHandler(os.Stdout, handler)))
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "nope")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
