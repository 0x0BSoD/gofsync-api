package hostgroups

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
	Parameters []HGParam            `json:"parameters"`
	NotExistPC []int                `json:"not_exist_pc"`
	DBHGExist  int                  `json:"dbhg_exist"`
	ExistId    int                  `json:"exist_id"`
}

// HTTP ============================

type HGListElem struct {
	ID        int    `json:"id"`
	ForemanID int    `json:"foreman_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
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

type POSTStructParameter struct {
	HGParam struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"parameter"`
}

type RackTablesSWE struct {
	Name      string `json:"name"`
	SweStatus string `json:"swestatus"`
}

type BatchPostStruct struct {
	ID          int    `json:"id"`
	HGName      string `json:"hgName"`
	THost       string `json:"tHost"`
	SHost       string `json:"sHost"`
	Environment struct {
		Name     string `json:"name"`
		TargetID int    `json:"targetId"`
	} `json:"environment"`
	Foreman struct {
		TargetID int `json:"targetId"`
		SourceID int `json:"sourceId"`
	} `json:"foreman"`
	InProgress bool `json:"in_progress"`
	Done       bool `json:"done"`
	HTTPResp   string
}
