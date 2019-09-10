package environment

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// ===============
// GET
// ===============
func ApiAll(host string, ctx *user.GlobalCTX) (Environments, error) {

	var result Environments
	bodyText, err := utils.ForemanAPI("GET", host, "environments", "", ctx)
	if err != nil {
		return Environments{}, err
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return Environments{}, err
	}

	return result, nil
}

func ApiGetSmartProxy(host string, ctx *user.GlobalCTX) int {
	var result SmartProxies
	bodyText, err := utils.ForemanAPI("GET", host, "smart_proxies", "", ctx)
	if err != nil {
		return -1
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return -1
	}

	return result.Results[0].ID
}

// ===============
// POST
// ===============
func Add(p EnvCheckP, ctx *user.GlobalCTX) error {
	jDataBase, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{
		Name: p.Env,
	})
	response, err := utils.ForemanAPI("POST", p.Host, "environments", string(jDataBase), ctx)
	if err != nil {
		utils.Error.Println(err)
		return err
	}

	if response.StatusCode == 201 {
		err = json.Unmarshal(response.Body, &response)
		if err != nil {
			utils.Error.Println(err)
			return err
		}
		return nil
	} else {
		fmt.Println(string(response.Body))
		fmt.Println(string(response.RequestUri))
		return fmt.Errorf("error on submit %s, code: %d", p.Env, response.StatusCode)
	}
}

// TODO: = required
func ImportPuppetClasses(p SweUpdateParams, ctx *user.GlobalCTX) (string, error) {
	//POST /api/environments/:environment_id/smart_proxies/:id/import_puppetclasses
	pID := ApiGetSmartProxy(p.Host, ctx)
	eID := ForemanID(p.Host, p.Environment, ctx)
	pApi, _ := json.Marshal(SweUpdatePOSTParams{
		DryRun: p.DryRun,
		Except: p.Except,
	})

	uri := fmt.Sprintf("environments/%d/smart_proxies/%d/import_puppetclasses", eID, pID)
	response, err := utils.ForemanAPI("POST", p.Host, uri, string(pApi), ctx)
	if err != nil {
		utils.Error.Println(err)
		return "", err
	}

	return string(response.Body), nil
}
