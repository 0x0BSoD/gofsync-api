package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

// our main function
func Server() {
	router := mux.NewRouter()

	// GET
	router.HandleFunc("/", Index).Methods("GET")
	// Hosts
	router.HandleFunc("/hosts", getAllHostsHttp).Methods("GET")
	// Host Groups
	router.HandleFunc("/hg", getAllHGListHttp).Methods("GET")
	router.HandleFunc("/hg/{host}", getHGListHttp).Methods("GET")
	router.HandleFunc("/hg/{host}/{swe_id}", getHGHttp).Methods("GET")
	router.HandleFunc("/hg/overrides/{hgName}", getOverridesByHGHttp).Methods("GET")
	// Locations
	router.HandleFunc("/loc/overrides/{locName}", getOverridesByLocHttp).Methods("GET")

	// POST
	router.HandleFunc("/send/hg", postHGHttp).Methods("POST")
	router.HandleFunc("/hg/check", postHGCheckHttp).Methods("POST")
	router.HandleFunc("/env/check", postEnvCheckHttp).Methods("POST")

	// Run Server
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Welcome to the HomePage!")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
