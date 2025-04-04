package hosts

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"git.ringcentral.com/archops/goFsync/utils"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"sync"
)

// ===============================
// CHECKS
// ===============================

// ===============================
// GET
// ===============================
func HostById(host string, id int, ctx *user.GlobalCTX) models.Response {
	uri := fmt.Sprintf("hosts/%d", id)
	response, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err != nil {
		logger.Error.Println(err)
	}
	return response
}

func ByHostgroup(host string, hg string, ctx *user.GlobalCTX) models.Response {
	uri := fmt.Sprintf("hostgroups/%s/hosts?format=json&per_page=%d", hg, ctx.Config.Api.GetPerPage)
	response, err := utils.ForemanAPI("GET", host, uri, "", ctx)
	if err != nil {
		logger.Error.Println(err)
	}
	return response
}

// ===============================
// POST
// ===============================

type NewHost struct {
	Host struct {
		// required
		Name string `json:"name"`
		// required for us
		EnvironmentID int `json:"environment_id"`
		// required for us
		HostgroupID int `json:"hostgroup_id"`
		//DomainID          int    `json:"domain_id"`
		//PtableID          int    `json:"ptable_id"`
		//MediumID          int    `json:"medium_id"`
		ArchitectureID    interface{} `json:"architecture_id"`
		OperatingsystemID interface{} `json:"operatingsystem_id"`
		//PuppetProxyID     int    `json:"puppet_proxy_id"`
		//ComputeResourceID int    `json:"compute_resource_id"`
		//RootPass          string `json:"root_pass"`
		// required
		LocationID int `json:"location_id"`
		// required
		//OrganizationID    int    `json:"organization_id"`
		//IP                string `json:"ip"`
		//Mac               string `json:"mac"`
		//ComputeAttributes struct {
		//	Cpus   int `json:"cpus"`
		//	Memory int `json:"memory"`
		//} `json:"compute_attributes"`
		//Enabled bool `json:"enabled"`
	} `json:"host"`
}

//type InterfacesAttributes struct {
//	Destroy int `json:"_destroy"`
//	Type string `json:"type"`
//	Mac string `json:"mac"`
//	Identifier string `json:"identifier"`
//	Name string `json:"name"`
//
//}
//type NewHost struct {
//	Host struct {
//		Name          string `json:"name"`
//		LocationID    int    `json:"location_id"`
//		HostgroupID   int    `json:"hostgroup_id"`
//		EnvironmentID int    `json:"environment_id"`
//		Managed       bool   `json:"managed"`
//		Type          string `json:"type"`
//		IsOwned       string `json:"is_owned"`
//		Enabled       bool   `json:"enabled"`
//		InterfacesAttributes
//	} `json:"host"`
//}

type NewHostParams struct {
	Name        string `json:"name"`
	ForemanHost string `json:"foreman_host"`
	Location    string `json:"location"`
	HostGroup   string `json:"hostgroup"`
	Environment string `json:"environment"`
}

func AddHost(host string, new NewHost, ctx *user.GlobalCTX) error {
	jDataStr, err := json.Marshal(new)
	if err != nil {
		return err
	}
	_, err = logger.ForemanAPI("POST", host, "hosts", string(jDataStr), ctx)
	if err != nil {
		return err
	}
	logger.Info.Printf("created new host: %s on %s", new.Host.Name, host)

	return nil
}

//func ByHostgroupName(hgName string, params url.Values, ctx *user.GlobalCTX) map[string][]Host {
//	result := make(map[string][]Host)
//	for _, host := range ss.Config.Hosts {
//		id := hostgroups.CheckHGID(hgName, host, ss)
//		uri := "hostgroups/%d/hosts?format=json&per_page=%d"
//		if id != -1 {
//			if val, ok := params["changed"]; ok {
//				p := strings.Trim(val[0], " ")
//				if strings.HasPrefix(p, "<") || strings.HasPrefix(p, ">") {
//					uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report+%s", id, ss.Config.Api.GetPerPage, url.QueryEscape(p))
//				} else {
//					uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report%%3D+%s", id, ss.Config.Api.GetPerPage, url.QueryEscape(p))
//				}
//				fmt.Println(uri)
//			} else {
//				uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d", id, ss.Config.Api.GetPerPage)
//			}
//			response, err := utils.ForemanAPI("GET", host, uri, "", ss.Config)
//			if err != nil {
//				logger.Error.Println(err)
//			}
//			if response.StatusCode == 404 {
//				logger.Error.Println("not found")
//			}
//
//			var tmpResult Hosts
//			err = json.Unmarshal(response.Body, &tmpResult)
//			if err != nil {
//				logger.Error.Println(err)
//			}
//
//			for _, i := range tmpResult.Results {
//				result[host] = append(result[host], i)
//			}
//		}
//	}
//	return result
//}

// Result struct
type HResult struct {
	sync.Mutex
	hosts map[string][]string
}

func (r *HResult) Add(foreman string, hostname string) {
	r.Lock()
	r.hosts[foreman] = append(r.hosts[foreman], hostname)
	r.Unlock()
}
func (r *HResult) Init() {
	r.hosts = make(map[string][]string)
}

//func ByHostgroupNameHostNames(hgName string, params url.Values, ctx *user.GlobalCTX) map[string][]string {
//	var r HResult
//	r.Init()
//	// Create a new WorkQueue.
//	wq := utils.New()
//	// This sync.WaitGroup is to make sure we wait until all of our work
//	// is done.
//	var wg sync.WaitGroup
//	for _, host := range ss.Config.Hosts {
//		wg.Add(1)
//		go func(h string) {
//			wq <- func() {
//				defer wg.Done()
//				id := hostgroups.CheckHGID(hgName, h, ss)
//				uri := "hostgroups/%d/hosts?format=json&per_page=%d"
//				if id != -1 {
//					if val, ok := params["changed"]; ok {
//						p := strings.Trim(val[0], " ")
//						if strings.HasPrefix(p, "<") || strings.HasPrefix(p, ">") {
//							uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report+%s", id, ss.Config.Api.GetPerPage, url.QueryEscape(p))
//						} else {
//							uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d&search=last_report%%3D+%s", id, ss.Config.Api.GetPerPage, url.QueryEscape(p))
//						}
//					} else {
//						uri = fmt.Sprintf("hostgroups/%d/hosts?format=json&per_page=%d", id, ss.Config.Api.GetPerPage)
//					}
//					response, err := utils.ForemanAPI("GET", h, uri, "", ss.Config)
//					if err != nil {
//						logger.Error.Println(err)
//						utils.GetErrorContext(err)
//					}
//					if response.StatusCode == 404 {
//						logger.Error.Println("not found")
//					}
//
//					var tmpResult Hosts
//					err = json.Unmarshal(response.Body, &tmpResult)
//					if err != nil {
//						logger.Error.Println(err)
//						utils.GetErrorContext(err)
//					}
//
//					for _, i := range tmpResult.Results {
//						r.Add(h, i.Name)
//					}
//				}
//			}
//		}(host)
//	}
//
//	// Wait for all of the work to finish, then close the WorkQueue.
//	wg.Wait()
//	close(wq)
//	return r.hosts
//}
