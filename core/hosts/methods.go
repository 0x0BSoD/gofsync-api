package hosts

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/environment"
	"git.ringcentral.com/archops/goFsync/core/hostgroups"
	"git.ringcentral.com/archops/goFsync/core/locations"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/utils"
	"strings"
)

func CreateNewHost(p NewHostParams, ctx *user.GlobalCTX) ([]byte, error) {
	var result NewHost

	result.Host.Name = p.Name
	result.Host.Managed = true
	result.Host.Type = "Host::Managed"
	result.Host.Enabled = true
	result.Host.IsOwned = "93-Users"

	if p.Environment == "" {
		p.Environment = "swe" + strings.Split(p.HostGroup, ".")[1]
	}

	if p.Location == "" {
		return []byte{}, fmt.Errorf("error on addnig host: %s, location not specified", p.Location)
	}

	hgID := hostgroups.ForemanID(p.ForemanHost, p.HostGroup, ctx)
	if hgID == -1 {
		return []byte{}, fmt.Errorf("error on addnig host: %s, not exist", p.HostGroup)
	}
	result.Host.HostgroupID = hgID

	envID := environment.ForemanID(p.ForemanHost, p.Environment, ctx)
	if envID == -1 {
		return []byte{}, fmt.Errorf("error on addnig host: %s, not exist", p.Environment)
	}
	result.Host.EnvironmentID = envID

	locID := locations.ForemanID(p.ForemanHost, p.Location, ctx)
	if locID == -1 {
		return []byte{}, fmt.Errorf("error on addnig host: %s, not exist", p.Location)
	}
	result.Host.LocationID = locID

	fmt.Println(result)

	jDataBase, _ := json.Marshal(result)

	response, _ := utils.ForemanAPI("POST", p.ForemanHost, "hosts", string(jDataBase), ctx)

	return response.Body, nil
}
