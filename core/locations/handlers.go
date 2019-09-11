package locations

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/locations/info"
	"git.ringcentral.com/archops/goFsync/middleware"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// ===============================
// GET
// ===============================
func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)

	var res = make([]AllLocations, len(ctx.Config.Hosts))

	for _, host := range ctx.Config.Hosts {
		dash := info.Get(host, ctx)
		locations, env := DbAll(host, ctx)
		tmp := AllLocations{
			Host:      host,
			Env:       env,
			Dashboard: dash,
			Open:      []bool{false},
		}
		for _, loc := range locations {
			tmp.Locations = append(tmp.Locations, Loc{
				Name:        loc,
				Highlighted: false,
			})
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
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	Sync(params["host"], ctx)
	err := json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}
