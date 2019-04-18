package locations

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"net/http"
)

func GetAllLocHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res []models.AllLocations
		w.Header().Set("Content-Type", "application/json")
		for _, host := range cfg.Hosts {
			data := GetAllLocNames(host, cfg)
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
}
