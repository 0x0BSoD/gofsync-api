package entitys

// PuppetClasses container
type PuppetClasses struct {
	Results  map[string][]*PuppetClass `json:"results"`
	Total    int                       `json:"total"`
	SubTotal int                       `json:"subtotal"`
	Page     int                       `json:"page"`
	PerPage  int                       `json:"per_page"`
	Search   string                    `json:"search"`
}

// PuppetClass structure
type PuppetClass struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PuppetClassName
type PuppetClassName struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	ModuleName           string                 `json:"module_name"`
	CreatedAt            string                 `json:"created_at"`
	UpdatedAt            string                 `json:"updated_at"`
	SmartVariables       []string               `json:"smart_variables"`
	SmartClassParameters []*SmartClassParameter `json:"smart_class_parameters"`
	Environments         []*Env                 `json:"environments"`
	HostGroups           []*HG                  `json:"hostgroups"`
}

// SmartClassParameter
type SmartClassParameter struct {
	Parameter string `json:"parameter"`
	ID        int    `json:"id"`
}

// SmartClassParameter
type HG struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	ID    int    `json:"id"`
}

// SmartClassParameter
type Env struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// Result for dbActions
type Result struct {
	ClassID           int
	ClassName         string
	PuppetSCOverrides string
	SCID              int
}

// Override Variables
type SCPOverride struct {
	Parameter           string            `json:"parameter"`
	ID                  int               `json:"id"`
	Description         string            `json:"description"`
	Override            bool              `json:"override"`
	ParameterType       string            `json:"parameter_type"`
	DefaultValue        interface{}       `json:"default_value"`
	UsePuppetDefault    bool              `json:"use_puppet_default"`
	Required            bool              `json:"required"`
	ValidatorType       string            `json:"validator_type"`
	ValidatorRule       string            `json:"validator_rule"`
	MergeOverrides      bool              `json:"merge_overrides"`
	AvoidDuplicates     bool              `json:"avoid_duplicates"`
	OverrideValueOrder  string            `json:"override_value_order"`
	OverrideValuesCount int               `json:"override_values_count"`
	CreatedAt           string            `json:"created_at"`
	UpdatedAt           string            `json:"updated_at"`
	PuppetClass         *PClass           `json:"puppetclass"`
	OverrideValues      []*OverrideValues `json:"override_values"`
}

// OverrideValues
type OverrideValues struct {
	ID               int         `json:"id"`
	Match            string      `json:"match"`
	Value            interface{} `json:"value"`
	UsePuppetDefault bool        `json:"use_puppet_default"`
}

// PClass
type PClass struct {
	Name       string `json:"name"`
	ModuleMame string `json:"module_name"`
	ID         int    `json:"id"`
}

// For inserting into base
type SCPOverrideForBase struct {
	Name               string
	ClassID int
	Override           bool
	ValidatorType      string
	OverrideValueOrder string
	DefaultValue       interface{}
	OverrideValues     []*OverrideValues
}
