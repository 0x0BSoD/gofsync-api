package DB

import (
	"git.ringcentral.com/archops/goFsync/core/puppetclass/API"
)

type HostGroupJSON struct {
	ID            int                          `json:"id,omitempty"`
	ForemanID     int                          `json:"foreman_id,omitempty"`
	Name          string                       `json:"name,omitempty"`
	SourceName    string                       `json:"source_name,omitempty"`
	Status        string                       `json:"status,omitempty"`
	Environment   string                       `json:"environment,omitempty"`
	ParentId      string                       `json:"parent_id,omitempty"`
	Params        []HostGroupParameter         `json:"params,omitempty,omitempty"`
	PuppetClasses map[string][]API.PuppetClass `json:"puppet_classes,omitempty"`
	Updated       string                       `json:"updated,omitempty"`
}

type HostGroupParameter struct {
	ForemanID int    `json:"foreman_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}
