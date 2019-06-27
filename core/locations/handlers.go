package locations

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/models"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// ===============================
// GET
// ===============================
func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	session := middleware.GetConfig(r)
	var res []models.AllLocations

	for _, host := range session.Config.Hosts {
		locs, env := DbAll(host, &session)
		tmp := models.AllLocations{
			Host: host,
			Env:  env,
		}
		for _, loc := range locs {
			tmp.Locations = append(tmp.Locations, loc)
		}
		res = append(res, tmp)
	}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Error.Printf("Error on getting all locations: %s", err)
	}
}

// ===============================
// POST
// ===============================
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	session := middleware.GetConfig(r)
	params := mux.Vars(r)

	Sync(params["host"], &session)
	err := json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}
