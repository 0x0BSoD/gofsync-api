package DB

type SmartClass struct {
	ID                  int         `json:"id"`
	ForemanID           int         `json:"foreman_id"`
	Name                string      `json:"name"`
	OverrideValuesCount int         `json:"override_values_count"`
	ValueType           string      `json:"value_type"`
	DefaultVal          interface{} `json:"default_val"`
	PuppetClass         string      `json:"puppet_class"`
	Override            []Override  `json:"override"`
	Dump                string      `json:"dump"`
}

type APISmartClass struct {
	Parameter   string `json:"parameter"`
	PuppetClass struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ModuleName string `json:"module_name"`
	} `json:"puppetclass"`
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

type OverrideValue struct {
	ForemanID        int         `json:"id"`
	Match            string      `json:"match"`
	Value            interface{} `json:"value"`
	UsePuppetDefault bool        `json:"use_puppet_default"`
}

type Override struct {
	ID         int         `json:"override_id"`
	ForemanID  int         `json:"foreman_id"`
	Parameter  string      `json:"parameter"`
	Match      string      `json:"match"`
	Value      string      `json:"value"`
	SmartClass *SmartClass `json:"smart_class"`
}
