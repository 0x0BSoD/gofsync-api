package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type aLocations struct {
	Host      string   `json:"host"`
	Locations []string `json:"locations"`
}

func getAllLocHttp(w http.ResponseWriter, r *http.Request) {
	var res []aLocations
	//getAllLocations
	w.Header().Set("Content-Type", "application/json")
	for _, host := range globConf.Hosts {
		data := getAllLocNames(host)
		tmp := aLocations{
			Host: host,
		}
		for _, loc := range data {
			tmp.Locations = append(tmp.Locations, loc)
		}
		res = append(res, tmp)
	}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("Error on getting HG list: %s", err)
	}
}
