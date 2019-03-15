package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// our main function
func Server() {
	router := mux.NewRouter()
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/send/hg/{shost}/{vhost}/{swe_id}", postHGHttp).Methods("GET")
	router.HandleFunc("/hg/{host}", getHGListHttp).Methods("GET")
	router.HandleFunc("/hg/{host}/{swe_id}", getHGHttp).Methods("GET")
	router.HandleFunc("/hg/overrides/{hgName}", getOverridesByHGHttp).Methods("GET")
	router.HandleFunc("/loc/overrides/{locName}", getOverridesByLocHttp).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Welcome to the HomePage!")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
