package API

import "sync"

type Parameters struct {
	Total    int         `json:"total"`
	SubTotal int         `json:"subtotal"`
	Page     int         `json:"page"`
	PerPage  int         `json:"per_page"`
	Search   string      `json:"search"`
	Results  []Parameter `json:"results"`
}

type Parameter struct {
	ID          int    `json:"id"`
	Parameter   string `json:"parameter"`
	PuppetClass struct {
		ForemanID  int    `json:"id"`
		Name       string `json:"name"`
		ModuleName string `json:"module_name"`
	} `json:"puppetclass"`
	ForemanID           int             `json:"id"`
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

type OverrideValues struct {
	Total    int             `json:"total"`
	SubTotal int             `json:"subtotal"`
	Page     int             `json:"page"`
	PerPage  int             `json:"per_page"`
	Search   string          `json:"search"`
	Results  []OverrideValue `json:"results"`
}

type OverrideValue struct {
	ForemanID        int         `json:"id"`
	Match            string      `json:"match"`
	Value            interface{} `json:"value"`
	UsePuppetDefault bool        `json:"use_puppet_default"`
}

type Result struct {
	sync.Mutex
	parameters []Parameter
}

func (r *Result) Add(ID Parameter) {
	r.Lock()
	r.parameters = append(r.parameters, ID)
	r.Unlock()
}
