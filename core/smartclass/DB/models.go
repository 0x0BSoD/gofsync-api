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
