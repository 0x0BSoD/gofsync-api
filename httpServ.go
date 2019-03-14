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
	router.HandleFunc("/swe/{host}", getHGListHttp).Methods("GET")
	router.HandleFunc("/swe/{host}/{swe_id}", getHGHttp).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Welcome to the HomePage!")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
