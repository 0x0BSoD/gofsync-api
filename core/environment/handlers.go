package environment

import (
	"encoding/json"
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

	data := DbAll(ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetByHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	data := DbByHost(params["host"], ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetSvnInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)

	data := RemoteGetSVNInfo(ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetSvnInfoHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	data := RemoteGetSVNInfoHost(params["host"], ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetSvnInfoName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	data := RemoteGetSVNInfoName(params["host"], params["name"], ctx)
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

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	Sync(params["host"], ctx)
	err := json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func PostCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var t EnvCheckP
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	data := DbID(t.Host, t.Env, ctx)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func ForemanPostCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var t EnvCheckP
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	data := DbForemanID(t.Host, t.Env, ctx)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func ForemanUpdatePCSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var t SweUpdateParams
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	ImportPuppetClasses(t, ctx)

	//data := DbForemanID(t.Host, t.Env, ctx)
	err = json.NewEncoder(w).Encode("triggered")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}
