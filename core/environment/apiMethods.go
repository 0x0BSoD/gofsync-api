package environment

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
)

// ===============
// GET
// ===============
func ApiGetAll(hostname string, ctx *user.GlobalCTX) (Environments, error) {
	var result Environments

	bodyText, err := utils.ForemanAPI("GET", hostname, "environments", "", ctx)
	if err != nil {
		return Environments{}, err
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return Environments{}, err
	}

	return result, nil
}

func ApiGet(hostname, envName string, ctx *user.GlobalCTX) (*Environment, error) {
	var result Environments
	uri := fmt.Sprintf("environments?search=name+=+%s", envName)
	bodyText, err := utils.ForemanAPI("GET", hostname, uri, "", ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyText.Body, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Results) == 1 {
		return result.Results[0], nil
	} else if len(result.Results) > 1 {
		return nil, fmt.Errorf("to mutch results for %s on %s", envName, hostname)
	} else {
		return nil, fmt.Errorf("error on getting env data from %s", hostname)
	}
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
	_, err := ApiGet(p.Host, p.Env, ctx)
	if err != nil {
		locationIDs := locations.DbAllForemanID(ctx.Config.Hosts[p.Host], ctx)

		params := NewEnv{
			Environment: NewEnvParams{
				Name:         p.Env,
				LocationsIDs: locationIDs,
			},
		}

		obj, _ := json.Marshal(params)

		response, err := utils.ForemanAPI("POST", p.Host, "environments", string(obj), ctx)
		if err != nil {
			return err
		}

		if response.StatusCode == 201 || response.StatusCode == 200 {
			err = json.Unmarshal(response.Body, &response)
			if err != nil {
				return err
			}

			err = AddNewEnv(p.Host, p.Env, ctx)
			if err != nil {
				return err
			}

			return nil
		} else {
			return fmt.Errorf("error on submit %s, code: %d", p.Env, response.StatusCode)
		}
	}

	return nil
}

func ImportPuppetClasses(p SweUpdateParams, ctx *user.GlobalCTX) (string, error) {
	pID := ApiGetSmartProxy(p.Host, ctx)
	eID := ForemanID(ctx.Config.Hosts[p.Host], p.Environment, ctx)
	pApi, _ := json.Marshal(SweUpdatePOSTParams{
		DryRun: p.DryRun,
		Except: p.Except,
	})

	uri := fmt.Sprintf("environments/%d/smart_proxies/%d/import_puppetclasses", eID, pID)
	response, err := utils.ForemanAPI("POST", p.Host, uri, string(pApi), ctx)
	if err != nil {
		return "", err
	}

	return string(response.Body), nil
}
