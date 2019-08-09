package DB

import "git.ringcentral.com/archops/goFsync/core/smartclass/DB"

type HostGroupJSON struct {
	ID            int                               `json:"id,omitempty"`
	ForemanID     int                               `json:"foreman_id,omitempty"`
	Name          string                            `json:"name,omitempty"`
	SourceName    string                            `json:"source_name,omitempty"`
	Status        string                            `json:"status,omitempty"`
	Environment   string                            `json:"environment,omitempty"`
	ParentId      string                            `json:"parent_id,omitempty"`
	Params        []HostGroupParameter              `json:"params,omitempty,omitempty"`
	PuppetClasses map[string][]HostGroupPuppetCLass `json:"puppet_classes,omitempty"`
	Updated       string                            `json:"updated,omitempty"`
}

type HostGroupPuppetCLass struct {
	ID           int             `json:"id"`
	ForemanID    int             `json:"foreman_id"`
	Class        string          `json:"class"`
	Subclass     string          `json:"subclass"`
	SmartClasses []DB.SmartClass `json:"smart_classes"`
}

type HostGroupParameter struct {
	ForemanID int    `json:"foreman_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}
