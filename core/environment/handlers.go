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

func GetSvnLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	envData := DbGet(params["host"], params["name"], ctx)
	data := RemoteGetSVNLog(params["host"], params["name"], envData.Repo, ctx)
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

	DirData, err := RemoteDIRGetSVNInfoName(params["host"], params["name"], ctx)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
	}

	envData := DbGet(params["host"], params["name"], ctx)
	UrlData, err := RemoteURLGetSVNInfoName(params["host"], params["name"], envData.Repo, ctx)
	if err != nil {
		logger.Error.Printf("Error on getting HG list: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	_ = json.NewEncoder(w).Encode(struct {
		Directory  logger.SvnInfo `json:"directory"`
		Repository logger.SvnInfo `json:"repository"`
	}{
		Directory:  DirData,
		Repository: UrlData,
	})
}

func GetSvnRepo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := DbGetRepo(params["host"], ctx)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error.Printf("Error on getting SVN Repo: %s", err)
	}
}

// ===============================
// POST
// ===============================
func SetSvnRepo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var b struct {
		Url string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		logger.Error.Printf("Error on POST EnvCheck: %s", err)
	}
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	DbSetRepo(b.Url, params["host"], ctx)
	err = json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on getting SVN Repo: %s", err)
	}
}

func SvnUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)

	var b struct {
		Host        string `json:"host"`
		Environment string `json:"environment"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		logger.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	RemoteSVNUpdate(b.Host, b.Environment, ctx)

	err = json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func SvnCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)

	var b struct {
		Host        string `json:"host"`
		Environment string `json:"environment"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		logger.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	envData := DbGet(b.Host, b.Environment, ctx)
	RemoteSVNCheckout(b.Host, b.Environment, envData.Repo, ctx)

	err = json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}

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
