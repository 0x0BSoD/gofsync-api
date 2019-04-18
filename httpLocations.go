package main

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"net/http"
)

func getAllLocHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res []models.GetAllLocations
		//getAllLocations
		w.Header().Set("Content-Type", "application/json")
		for _, host := range globConf.Hosts {
			data := getAllLocNames(host, cfg)
			tmp := models.GetAllLocations{
				Host: host,
			}
			for _, loc := range data {
				tmp.Locations = append(tmp.Locations, loc)
			}
			res = append(res, tmp)
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Error.Printf("Error on getting all locations: %s", err)
		}
	}
}
