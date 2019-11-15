package environment

import (
	"encoding/json"
	"fmt"
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

	ctx := middleware.GetContext(r)
	data := DbAll(ctx)

	utils.SendResponse(w, "error on getting HG list: %s", data)
}

func GetByHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := DbByHost(params["host"], ctx)

	utils.SendResponse(w, "error on getting HG: %s", data)
}

func GetSvnInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	data := RemoteGetSVNInfo(ctx)

	utils.SendResponse(w, "error on getting svn info: %s", data)
}

func GetSvnLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	envData := DbGet(params["host"], params["name"], ctx)
	data := RemoteGetSVNLog(params["host"], params["name"], envData.Repo, ctx)

	utils.SendResponse(w, "error on getting svn log: %s", data)
}

func GetSvnInfoHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := RemoteGetSVNInfoHost(params["host"], ctx)

	utils.SendResponse(w, "error on getting svn info on host: %s", data)
}

func GetSvnInfoName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	DirData, err := RemoteDIRGetSVNInfoName(params["host"], params["name"], ctx)
	if err != nil {
		utils.Error.Printf("[svn] error on getting info from server: %s", err)
	}

	envData := DbGet(params["host"], params["name"], ctx)

	UrlData, err := RemoteURLGetSVNInfoName(params["host"], params["name"], envData.Repo, ctx)
	if err != nil {
		utils.Error.Printf("[svn] error on getting info from svn: %s", err)
		http.Error(w, fmt.Sprintf("[svn] error on getting info from svn: %s", err), http.StatusInternalServerError)
		return
	}

	var response struct {
		Directory  SvnInfo `json:"directory"`
		Repository SvnInfo `json:"repository"`
	}

	if len(DirData.Entry.Path) != 0 {
		response.Directory = DirData
	} else {
		response.Directory.Entry.Path = "not exist"
	}

	if len(UrlData.Entry.Path) != 0 {
		response.Repository = UrlData
	} else {
		response.Repository.Entry.Path = "not exist"
	}
	_ = json.NewEncoder(w).Encode(response)
}

func GetSvnRepo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	data := DbGetRepo(params["host"], ctx)

	utils.SendResponse(w, "error on getting svn repo: %s", data)
}

// =====================================================================================================================
// POST
// =====================================================================================================================
func Submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var body EnvCheckP
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		utils.Error.Printf("error on submiting new environment: %s", err)
		return
	}

	fmt.Println(body)

	err = Add(body, ctx)
	if err != nil {
		utils.Error.Printf("error on submiting new environment: %s", err)
		w.WriteHeader(http.StatusNotModified)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// ==========
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("created"))
}

func SvnBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	var b map[string][]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}
	RemoteSVNBatch(b, ctx)

	// ==========
	err = json.NewEncoder(w).Encode(b)
	if err != nil {
		utils.Error.Printf("Error on getting SVN Repo: %s", err)
	}
}

func SetSvnRepo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	DbSetRepo(b.Url, b.Host, ctx)
	err = json.NewEncoder(w).Encode("submitted")
	if err != nil {
		utils.Error.Printf("Error on getting SVN Repo: %s", err)
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
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	out, err := RemoteSVNUpdate(b.Host, b.Environment, ctx)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
		w.WriteHeader(http.StatusNotModified)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
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
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	envData := DbGet(b.Host, b.Environment, ctx)
	out, err := RemoteSVNCheckout(b.Host, b.Environment, envData.Repo, ctx)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
		w.WriteHeader(http.StatusNotModified)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

func PostCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var t EnvCheckP
	err := decoder.Decode(&t)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	data := ID(t.Host, t.Env, ctx)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func ForemanPostCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var t EnvCheckP
	err := decoder.Decode(&t)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	data := ForemanID(t.Host, t.Env, ctx)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}

func ForemanUpdatePCSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := middleware.GetContext(r)
	decoder := json.NewDecoder(r.Body)
	var t SweUpdateParams
	err := decoder.Decode(&t)
	if err != nil {
		utils.Error.Printf("Error on POST EnvCheck: %s", err)
	}

	response, err := ImportPuppetClasses(t, ctx)
	if err != nil {
		utils.Error.Printf("error on import puppet classes: %s", err)
		w.WriteHeader(http.StatusNotModified)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Error.Printf("Error on EnvCheck: %s", err)
	}
}
