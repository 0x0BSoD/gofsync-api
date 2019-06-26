// Testing go-swagger generation
//
// The purpose of this application is to test go-swagger in a simple GET request.
//
//     Schemes: https
//     Host: sjc01-c01-pds10.c01.ringcentral.com:8086/api/v1/
//     Version: 1.3
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Alexander<alexander.simonov@nordigy.ru>
//
//     Consumes:
//     - text/json
//
//     Produces:
//     - text/json
//
// swagger:meta
package main

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/hosts"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

// our main function
func Server(cfg *models.Config) {
	router := mux.NewRouter()

	// GET =============================================================================================================
	router.HandleFunc("/", middleware.Chain(Index, middleware.Token(cfg))).Methods("GET")

	// Hosts
	router.HandleFunc("/hosts/foreman", middleware.Chain(hostgroups.GetAllHostsHttp, middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hosts/all/hg/{hgName}", middleware.Chain(hosts.ByHostgroupNameHttp, middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hosts/{host}/hg/{hgForemanId}", middleware.Chain(hosts.ByHostgroupHttp, middleware.Token(cfg))).Methods("GET")

	// Env
	router.HandleFunc("/env/{host}", middleware.Chain(environment.GetAll, middleware.Token(cfg))).Methods("GET")

	// Locations
	router.HandleFunc("/loc", middleware.Chain(locations.GetAll, middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/loc/overrides/{host}/{locName}", middleware.Chain(smartclass.GetOverridesByLocHttp, middleware.Token(cfg))).Methods("GET")

	// Puppet Classes
	router.HandleFunc("/pc/{host}", middleware.Chain(puppetclass.GetAll, middleware.Token(cfg))).Methods("GET")

	// Smart Classes
	router.HandleFunc("/sc/{sc_id}", middleware.Chain(smartclass.GetSCDataByIdHttp, middleware.Token(cfg))).Methods("GET")

	// Host Groups
	// GET ===
	router.HandleFunc("/hg", middleware.Chain(hostgroups.GetAllHGListHttp, middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/{host}", middleware.Chain(hostgroups.GetHGListHttp, middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/{host}/{swe_id}", middleware.Chain(hostgroups.GetHGHttp, middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/update/{host}/{hgName}", middleware.Chain(hostgroups.GetHGUpdateInBaseHttp(cfg), middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/get/{host}/{hgName}", middleware.Chain(hostgroups.GetHGFHttp, middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/foreman/check/{host}/{hgName}", middleware.Chain(hostgroups.GetHGCheckHttp, middleware.Token(cfg))).Methods("GET")
	router.HandleFunc("/hg/overrides/{hgName}", middleware.Chain(smartclass.GetOverridesByHGHttp, middleware.Token(cfg))).Methods("GET")
	// POST ===
	router.HandleFunc("/hg/update/{host}", middleware.Chain(hostgroups.Update, middleware.Token(cfg))).Methods("POST")
	router.HandleFunc("/hg/upload", middleware.Chain(hostgroups.Post, middleware.Token(cfg))).Methods("POST")
	//router.HandleFunc("/hg/batch/upload", middleware.Chain(hostgroups.BatchPost, middleware.Token(cfg))).Methods("POST")
	router.HandleFunc("/hg/create/{host}", middleware.Chain(hostgroups.Create, middleware.Token(cfg))).Methods("POST")
	router.HandleFunc("/hg/check", middleware.Chain(hostgroups.PostHGCheckHttp, middleware.Token(cfg))).Methods("POST")

	// POST Other ======================================================================================================
	router.HandleFunc("/", middleware.Chain(Index, middleware.Token(cfg))).Methods("POST")

	// User
	router.HandleFunc("/signin", user.SignIn(cfg)).Methods("POST")
	router.HandleFunc("/refreshjwt", user.Refresh(cfg)).Methods("POST")

	// Checks
	router.HandleFunc("/env/check", middleware.Chain(environment.PostCheck, middleware.Token(cfg))).Methods("POST")

	// Env
	router.HandleFunc("/env/{host}", middleware.Chain(environment.Update, middleware.Token(cfg))).Methods("POST")

	// Loc
	router.HandleFunc("/loc/{host}", middleware.Chain(locations.Update, middleware.Token(cfg))).Methods("POST")

	// Puppet Classes
	router.HandleFunc("/pc/update/{host}", middleware.Chain(puppetclass.Update, middleware.Token(cfg))).Methods("POST")

	// SocketIO ========================================================================================================
	if cfg.Web.SocketActive {
		//hub := utils.NewHub()
		//go hub.Run()
		router.HandleFunc("/ws", utils.Serve(cfg))
	}

	// Run Server
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://sjc01-c01-pds10:8086",
			"http://localhost:8080",
			"ws://localhost:8080",
			"wss://localhost:8080",
			"https://sjc01-c01-pds10:8086",
			"https://sjc01-c01-pds10.c01.ringcentral.com:8086"},
		AllowCredentials: true,
		Debug:            false,
	})
	handler := c.Handler(router)
	bindAddr := fmt.Sprintf(":%d", cfg.Web.Port)
	log.Fatal(http.ListenAndServe(bindAddr, handler))
}

func Index(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	_, err := fmt.Fprintf(w, "I'am a teapot")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
