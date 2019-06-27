package environment

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
	params := mux.Vars(r)

	data := DbAll(params["host"], &session)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
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

func PostCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	session := middleware.GetConfig(r)
	decoder := json.NewDecoder(r.Body)
	var t models.EnvCheckP
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	data := DbID(t.Host, t.Env, &session)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}
