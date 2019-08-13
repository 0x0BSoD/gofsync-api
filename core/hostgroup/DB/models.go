package DB

import "git.ringcentral.com/archops/goFsync/core/smartclass/DB"

type HostGroup struct {
	ForemanID           int    `json:"id"`
	Name                string `json:"name"`
	Title               string `json:"title"`
	SubnetID            int    `json:"subnet_id"`
	SubnetName          string `json:"subnet_name"`
	OperatingSystemID   int    `json:"operatingsystem_id"`
	OperatingSystemName string `json:"operatingsystem_name"`
	DomainID            int    `json:"domain_id"`
	DomainName          string `json:"domain_name"`
	EnvironmentID       int    `json:"environment_id"`
	EnvironmentName     string `json:"environment_name"`
	ComputeProfileId    int    `json:"compute_profile_id"`
	ComputeProfileName  string `json:"compute_profile_name"`
	Ancestry            string `json:"ancestry,omitempty"`
	PuppetProxyId       int    `json:"puppet_proxy_id"`
	PuppetCaProxyId     int    `json:"puppet_ca_proxy_id"`
	PTableId            int    `json:"ptable_id"`
	PTableName          string `json:"ptable_name"`
	MediumId            int    `json:"medium_id"`
	MediumName          string `json:"medium_name"`
	ArchitectureId      int    `json:"architecture_id"`
	ArchitectureName    int    `json:"architecture_name"`
	RealmId             int    `json:"realm_id"`
	RealmName           string `json:"realm_name"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

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
	Priority  int    `json:"priority"`
}
