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
	"log"
	"net/http"

	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/hosts"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/puppetclass"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// our main function
func Server(ctx *user.GlobalCTX) {
	router := mux.NewRouter()
	// GET =============================================================================================================

	// SocketIO ========================================================================================================
	router.HandleFunc("/ws", utils.WSServe(ctx))

	// Hosts
	router.HandleFunc("/hosts/foreman", middleware.Chain(hostgroups.GetAllHostsHttp, middleware.Token(ctx))).Methods("GET")
	//router.HandleFunc("/hosts/all/hg/{hgName}", middleware.Chain(hosts.ByHostgroupNameHttp, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/hosts/{host}/hg/{hgForemanId}", middleware.Chain(hosts.ByHostgroupHttp, middleware.Token(ctx))).Methods("GET")

	// Env
	//// Svn
	router.HandleFunc("/env/svn/info/all", middleware.Chain(environment.GetSvnInfo, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/env/svn/repo/{host}", middleware.Chain(environment.GetSvnRepo, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/env/svn/info/{host}", middleware.Chain(environment.GetSvnInfoHost, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/env/svn/info/{host}/{name}", middleware.Chain(environment.GetSvnInfoName, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/env/svn/log/{host}/{name}", middleware.Chain(environment.GetSvnLog, middleware.Token(ctx))).Methods("GET")
	//// Foreman
	router.HandleFunc("/env/{host}", middleware.Chain(environment.GetByHost, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/env", middleware.Chain(environment.GetAll, middleware.Token(ctx))).Methods("GET")
	// POST ===
	//// Svn
	router.HandleFunc("/env/svn/update", middleware.Chain(environment.SvnUpdate, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/env/svn/checkout", middleware.Chain(environment.SvnCheckout, middleware.Token(ctx))).Methods("POST")
	//router.HandleFunc("/env/svn/foreman", middleware.Chain(environment.ForemanUpdatePCSource, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/env/svn/repo", middleware.Chain(environment.SetSvnRepo, middleware.Token(ctx))).Methods("POST")
	//// Foreman
	router.HandleFunc("/env/foreman/check", middleware.Chain(environment.ForemanPostCheck, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/env/db/check", middleware.Chain(environment.PostCheck, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/env/db/update", middleware.Chain(environment.Update, middleware.Token(ctx))).Methods("POST")

	// Locations
	router.HandleFunc("/loc", middleware.Chain(locations.GetAll, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/loc/overrides/{host}/{locName}", middleware.Chain(smartclass.GetOverridesByLocHttp, middleware.Token(ctx))).Methods("GET")
	// POST ===
	router.HandleFunc("/loc/submit", middleware.Chain(hostgroups.SubmitLocation, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/loc/{host}", middleware.Chain(locations.Update, middleware.Token(ctx))).Methods("POST")

	// Puppet Classes
	router.HandleFunc("/pc/{host}", middleware.Chain(puppetclass.GetAll, middleware.Token(ctx))).Methods("GET")
	// POST ===
	router.HandleFunc("/pc/update/{host}", middleware.Chain(puppetclass.Update, middleware.Token(ctx))).Methods("POST")

	// Smart Classes
	router.HandleFunc("/sc/{sc_id}", middleware.Chain(smartclass.GetSCDataByIdHttp, middleware.Token(ctx))).Methods("GET")

	// Host Groups
	// GET ===
	router.HandleFunc("/hg", middleware.Chain(hostgroups.GetAllHGListHttp, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/hg/git/commit/{host}/{swe_id}", middleware.Chain(hostgroups.CommitGitHttp, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/hg/foreman/update/{host}/{hgName}", middleware.Chain(hostgroups.GetHGUpdateInBaseHttp, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/hg/foreman/get/{host}/{hgName}", middleware.Chain(hostgroups.GetHGFHttp, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/hg/foreman/check/{host}/{hgName}", middleware.Chain(hostgroups.GetHGCheckHttp, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/hg/overrides/{hgName}", middleware.Chain(smartclass.GetOverridesByHGHttp, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/hg/{host}/{swe_id}", middleware.Chain(hostgroups.GetHGHttp, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/hg/{host}", middleware.Chain(hostgroups.GetHGListHttp, middleware.Token(ctx))).Methods("GET")
	// POST ===
	router.HandleFunc("/hg/check", middleware.Chain(hostgroups.PostHGCheckHttp, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/hg/upload", middleware.Chain(hostgroups.Post, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/hg/batch/upload", middleware.Chain(hostgroups.BatchPost, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/hg/update/{host}", middleware.Chain(hostgroups.Update, middleware.Token(ctx))).Methods("POST")
	router.HandleFunc("/hg/create/{host}", middleware.Chain(hostgroups.Create, middleware.Token(ctx))).Methods("POST")

	router.HandleFunc("/", middleware.Chain(Index, middleware.Token(ctx))).Methods("GET")
	router.HandleFunc("/", middleware.Chain(Index, middleware.Token(ctx))).Methods("POST")

	// User
	// POST ===
	router.HandleFunc("/signin", user.SignIn(ctx)).Methods("POST")
	router.HandleFunc("/refreshjwt", user.Refresh(ctx)).Methods("POST")

	// Run Server
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
			"ws://localhost:8080",
			"wss://localhost:8080",
			"ws://localhost:8000",
			"wss://localhost:8000",
			"https://sjc01-c01-pds10:8086",
			"https://sjc01-c01-pds10.c01.ringcentral.com:8086",
			"ws://sjc01-c01-pds10:8086",
			"wss://sjc01-c01-pds10.c01.ringcentral.com:8086",
		},
		AllowCredentials: true,
		Debug:            false,
	})
	handler := c.Handler(router)
	bindAddr := fmt.Sprintf(":%d", ctx.Config.Web.Port)
	log.Fatal(http.ListenAndServe(bindAddr, handler))
}

func Index(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	_, err := fmt.Fprintf(w, "I'am a teapot")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
