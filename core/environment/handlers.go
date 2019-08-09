package environment

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/environment/DB"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// =====================================================================================================================
// GET
// =====================================================================================================================

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	var gDB DB.Get
	ctx := middleware.GetContext(r)
	data := gDB.All(ctx)

	// ======
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetByHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	var gDB DB.Get
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := gDB.ByHost(params["host"], ctx)

	// =======
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetSvnInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	data := RemoteGetSVNInfo(ctx)

	// ======
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetSvnLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)
	envData := gDB.ByName(params["host"], params["name"], ctx)
	data := RemoteGetSVNLog(params["host"], params["name"], envData.Repo, ctx)

	// ====
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetSvnInfoHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := RemoteGetSVNInfoHost(params["host"], ctx)

	// ====
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting HG list: %s", err)
	}
}

func GetSvnInfoName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)
	DirData, err := RemoteDIRGetSVNInfoName(params["host"], params["name"], ctx)
	if err != nil {
		utils.Error.Printf("Error on getting HG list: %s", err)
	}
	envData := gDB.ByName(params["host"], params["name"], ctx)
	UrlData, err := RemoteURLGetSVNInfoName(params["host"], params["name"], envData.Repo, ctx)
	if err != nil {
		utils.Error.Printf("Error on getting HG list: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	// ===========
	_ = json.NewEncoder(w).Encode(struct {
		Directory  DB.SvnInfo `json:"directory"`
		Repository DB.SvnInfo `json:"repository"`
	}{
		Directory:  DirData,
		Repository: UrlData,
	})
}

func GetSvnRepo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	params := mux.Vars(r)
	data := gDB.Repo(params["host"], ctx)

	// =========
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on getting SVN Repo: %s", err)
	}
}

// =====================================================================================================================
// POST
// =====================================================================================================================

func SetSvnRepo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	var uDB DB.Update
	var b struct {
		Host string `json:"host"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}
	ctx := middleware.GetContext(r)
	uDB.SetRepo(b.Url, b.Host, ctx)

	// ===========
	err = json.NewEncoder(w).Encode("submitted")
	if err != nil {
		utils.Error.Printf("Error on getting SVN Repo: %s", err)
	}
}

func SvnUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)

	var b struct {
		Host        string `json:"host"`
		Environment string `json:"environment"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	RemoteSVNUpdate(b.Host, b.Environment, ctx)

	err = json.NewEncoder(w).Encode("submitted")
	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func SvnCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	var b struct {
		Host        string `json:"host"`
		Environment string `json:"environment"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	envData := gDB.ByName(b.Host, b.Environment, ctx)
	RemoteSVNCheckout(b.Host, b.Environment, envData.Repo, ctx)

	// ========
	err = json.NewEncoder(w).Encode("submitted")
	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)

	var b struct {
		Host string `json:"host"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		utils.Error.Printf("Error on POST EnvUpdate: %s", err)
	}

	Sync(b.Host, ctx)
	err = json.NewEncoder(w).Encode("submitted")
	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var gDB DB.Get
	decoder := json.NewDecoder(r.Body)
	var t EnvCheckParameters
	err := decoder.Decode(&t)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	dataBase := gDB.ID(t.Host, t.Env, ctx)
	foreman := gDB.ForemanID(t.Host, t.Env, ctx)

	err = json.NewEncoder(w).Encode(CheckResponse{
		ID:        dataBase,
		ForemanID: foreman,
	})

	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}

//func ForemanUpdatePCSource(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	// VARS
//	ctx := middleware.GetContext(r)
//	decoder := json.NewDecoder(r.Body)
//	var t SweUpdateParams
//	err := decoder.Decode(&t)
//	if err != nil {
//		utils.Error.Printf("Error on POST EnvCheck: %s", err)
//	}
//
//	ImportPuppetClasses(t, ctx)
//
//	//data := DbForemanID(t.Host, t.Env, ctx)
//	err = json.NewEncoder(w).Encode("triggered")
//	if err != nil {
//		utils.Error.Printf("Error on EnvCheck: %s", err)
//	}
//}
