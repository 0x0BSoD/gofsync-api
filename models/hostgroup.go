package models

type HostGroup struct {
	ID                  int    `json:"id"`
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
type HostGroups struct {
	Results  []HostGroup `json:"results"`
	Total    int         `json:"total"`
	SubTotal int         `json:"subtotal"`
	Page     int         `json:"page"`
	PerPage  int         `json:"per_page"`
	Search   string      `json:"search"`
}

//  Host Group parameters
type HostGroupPContainer struct {
	Results  []HostGroupP `json:"results"`
	Total    int          `json:"total"`
	SubTotal int          `json:"subtotal"`
	Page     int          `json:"page"`
	PerPage  int          `json:"per_page"`
	Search   string       `json:"search"`
}
type HostGroupP struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Priority int    `json:"priority"`
}
type HostGroupS struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}
type HgError struct {
	ID        int    `json:"id"`
	HostGroup string `json:"host_group"`
	Host      string `json:"host"`
	Error     string `json:"error"`
}

// HostGroupBase Structure for post
type HostGroupBase struct {
	ParentId       int    `json:"parent_id"`
	Name           string `json:"name"`
	EnvironmentId  int    `json:"environment_id"`
	PuppetClassIds []int  `json:"puppetclass_ids"`
	LocationIds    []int  `json:"location_ids"`
	//Parameters     []int  `json:"group_parameters_attributes"`
}
type HostGroupOverrides struct {
	OvrForemanId int    `json:"ovr_foreman_id"`
	ScForemanId  int    `json:"sc_foreman_id"`
	Match        string `json:"match"`
	Value        string `json:"value"`
}
type HWPostRes struct {
	BaseInfo   HostGroupBase        `json:"hostgroup"`
	Overrides  []HostGroupOverrides `json:"override_value"`
	NotExistPC []int                `json:"not_exist_pc"`
	DBHGExist  int                  `json:"dbhg_exist"`
	ExistId    int                  `json:"exist_id"`
}

// HTTP ============================
type HGElem struct {
	ID            int                           `json:"id"`
	ForemanID     int                           `json:"foreman_id"`
	Name          string                        `json:"name"`
	SourceName    string                        `json:"source_name,omitempty"`
	Status        string                        `json:"status"`
	Environment   string                        `json:"environment"`
	ParentId      string                        `json:"parent_id"`
	Params        []HGParam                     `json:"params,omitempty"`
	PuppetClasses map[string][]PuppetClassesWeb `json:"puppet_classes"`
	Updated       string                        `json:"updated"`
}
type HGListElem struct {
	ID        int    `json:"id"`
	ForemanID int    `json:"foreman_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}
type HGParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HGPost struct {
	SourceHost string `json:"source_host"`
	TargetHost string `json:"target_host"`
	TargetHgId int    `json:"target_hg_id"`
	SourceHgId int    `json:"source_hg_id"`
	DBUpdate   bool   `json:"db_update"`
}
type ErrStruct struct {
	Message string
	State   string
}
type POSTStructBase struct {
	HostGroup HostGroupBase `json:"hostgroup"`
}
type POSTStructOvrVal struct {
	OverrideValue struct {
		Match string `json:"match"`
		Value string `json:"value"`
	} `json:"override_value"`
}

type RackTablesSWE struct {
	Name      string `json:"name"`
	SweStatus string `json:"swestatus"`
}
