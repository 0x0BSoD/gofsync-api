package main

import (
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/core/environment"
	"git.ringcentral.com/alexander.simonov/goFsync/core/hostgroups"
	"git.ringcentral.com/alexander.simonov/goFsync/core/locations"
	"git.ringcentral.com/alexander.simonov/goFsync/core/puppetclass"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/core/user"
	"git.ringcentral.com/alexander.simonov/goFsync/middleware"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

// our main function
func Server(cfg *models.Config) {
	router := mux.NewRouter()

	// GET ===
	router.HandleFunc("/", middleware.Chain(Index, middleware.Token(cfg))).Methods("GET")
	// Hosts
	router.HandleFunc("/hosts", middleware.Chain(hostgroups.GetAllHostsHttp(cfg), middleware.Token(cfg))).Methods("GET")
	// Env
	router.HandleFunc("/env/{host}", middleware.Chain(environment.GetAllEnv(cfg), middleware.Token(cfg))).Methods("GET")
	// Locations
	router.HandleFunc("/loc", middleware.Chain(locations.GetAllLocHttp(cfg), middleware.Token(cfg))).Methods("GET")
	// Puppet Classes
	router.HandleFunc("/pc/{host}", middleware.Chain(puppetclass.GetAllPCHttp(cfg), middleware.Token(cfg))).Methods("GET")
	// Smart Classes
	router.HandleFunc("/sc/{sc_id}", middleware.Chain(smartclass.GetSCDataByIdHttp(cfg), middleware.Token(cfg))).Methods("GET")
	// Host Groups
	router.HandleFunc("/hg", middleware.Chain(hostgroups.GetAllHGListHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/{host}", middleware.Chain(hostgroups.GetHGListHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/{host}/{swe_id}", middleware.Chain(hostgroups.GetHGHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/update/{host}/{hgName}", middleware.Chain(hostgroups.GetHGUpdateInBaseHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/get/{host}/{hgName}", middleware.Chain(hostgroups.GetHGFHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/check/{host}/{hgName}", middleware.Chain(hostgroups.GetHGCheckHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/overrides/{hgName}", middleware.Chain(smartclass.GetOverridesByHGHttp(cfg), middleware.Token(cfg))).Methods("GET")
	// Locations
	router.HandleFunc("/loc/overrides/{host}/{locName}", middleware.Chain(smartclass.GetOverridesByLocHttp(cfg), middleware.Token(cfg))).Methods("GET")

	// POST ===
	// User ===
	router.HandleFunc("/signin", user.SignIn(cfg)).Methods("POST")
	router.HandleFunc("/refreshjwt", user.Refresh(cfg)).Methods("POST")

	router.HandleFunc("/hg/upload", middleware.Chain(hostgroups.PostHGHttp(cfg), middleware.Token(cfg))).Methods("POST")
	router.HandleFunc("/hg/check", middleware.Chain(hostgroups.PostHGCheckHttp(cfg), middleware.Token(cfg))).Methods("POST")
	router.HandleFunc("/env/check", middleware.Chain(environment.PostEnvCheckHttp(cfg), middleware.Token(cfg))).Methods("POST")

	// SocketIO
	router.HandleFunc("/ws", utils.Serve(cfg))

	// Run Server
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://sjc01-c01-pds10:8086", "http://localhost:8080", "ws://localhost:8080",
			"https://sjc01-c01-pds10:8086", "https://sjc01-c01-pds10.c01.ringcentral.com:8086"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	handler := c.Handler(router)
	bindAddr := fmt.Sprintf(":%d", cfg.Web.Port)
	//log.Fatal(http.ListenAndServe(bindAddr, handlers.LoggingHandler(os.Stdout, handler)))
	log.Fatal(http.ListenAndServe(bindAddr, handler))
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	_, err := fmt.Fprintf(w, "I'am a teapot")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
