package models

// Smart Class Container
type SCParameters struct {
	Total    int           `json:"total"`
	SubTotal int           `json:"subtotal"`
	Page     int           `json:"page"`
	PerPage  int           `json:"per_page"`
	Search   string        `json:"search"`
	Results  []SCParameter `json:"results"`
}

// Smart Class
type SCParameter struct {
	Parameter           string          `json:"parameter"`
	PuppetClass         PuppetClassInSc `json:"puppetclass"`
	ID                  int             `json:"id"`
	Description         string          `json:"description"`
	Override            bool            `json:"override"`
	ParameterType       string          `json:"parameter_type"`
	DefaultValue        interface{}     `json:"default_value"`
	UsePuppetDefault    bool            `json:"use_puppet_default"`
	Required            bool            `json:"required"`
	ValidatorType       string          `json:"validator_type"`
	ValidatorRule       string          `json:"validator_rule"`
	MergeOverrides      bool            `json:"merge_overrides"`
	AvoidDuplicates     bool            `json:"avoid_duplicates"`
	OverrideValueOrder  string          `json:"override_value_order"`
	OverrideValuesCount int             `json:"override_values_count"`
	OverrideValues      []OverrideValue `json:"override_values"`
}

// PC for old Foreman
type PuppetClassInSc struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ModuleName string `json:"module_name"`
}

// OverrideValues Container
type OverrideValues struct {
	Total    int             `json:"total"`
	SubTotal int             `json:"subtotal"`
	Page     int             `json:"page"`
	PerPage  int             `json:"per_page"`
	Search   string          `json:"search"`
	Results  []OverrideValue `json:"results"`
}
type OverrideValue struct {
	ID               int         `json:"id"`
	Match            string      `json:"match"`
	Value            interface{} `json:"value"`
	UsePuppetDefault bool        `json:"use_puppet_default"`
}
type PCSCParameters struct {
	ID                   int             `json:"id"`
	Name                 string          `json:"name"`
	ModuleName           string          `json:"module_name"`
	SmartClassParameters []PCSCParameter `json:"smart_class_parameters"`
	Environments         []Environment   `json:"environments"`
	HostGroups           []HostGroupS    `json:"hostgroups"`
}
type PCSCParameter struct {
	ID        int    `json:"id"`
	Parameter string `json:"parameter"`
}

// Return From Base
type SCGetResAdv struct {
	ID                  int
	ForemanId           int
	Name                string
	OverrideValuesCount int
	ValueType           string
	DefaultVal          interface{}
	PuppetClass         string
	Override            []SCOParams
	Dump                string
}
type SmartClass struct {
	Id        int    `json:"id"`
	ForemanId int    `json:"foreman_id"`
	Name      string `json:"name"`
}
type SCOParams struct {
	SmartClassId int    `json:"smart_class_id"`
	OverrideId   int    `json:"override_id"`
	Parameter    string `json:"parameter"`
	Match        string `json:"match"`
	Value        string `json:"value"`
}
type LogStatus struct {
	Name          string `json:"name"`
	Host          string `json:"host"`
	Current       int    `json:"current"`
	CurrentThread int    `json:"current_thread,omitempty"`
	TotalInThread int    `json:"total_in_thread,omitempty"`
	Total         int    `json:"total"`
}

type OvrParams struct {
	OvrForemanId int `json:"ovr_foreman_id"`
	SCForemanId  int `json:"sc_foreman_id"`
	//Parameter      string `json:"parameter"`
	Type           string `json:"type"`
	PuppetClass    string `json:"puppet_class"`
	SmartClassName string `json:"smart_class_name"`
	Value          string `json:"value"`
}
