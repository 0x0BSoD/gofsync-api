package locations

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/middleware"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"net/http"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	var res []models.AllLocations
	for _, host := range cfg.Hosts {
		data := DbAll(host, cfg)
		tmp := models.AllLocations{
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
