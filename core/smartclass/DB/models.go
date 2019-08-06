package DB

type SmartClass struct {
	ID                  int         `json:"id"`
	ForemanId           int         `json:"foreman_id"`
	Name                string      `json:"name"`
	OverrideValuesCount int         `json:"override_values_count"`
	ValueType           string      `json:"value_type"`
	DefaultVal          interface{} `json:"default_val"`
	PuppetClass         string      `json:"puppet_class"`
	Override            []Overrides `json:"override"`
	Dump                string      `json:"dump"`
}

type Overrides struct {
	SmartClassId int    `json:"smart_class_id"`
	OverrideId   int    `json:"override_id"`
	Parameter    string `json:"parameter"`
	Match        string `json:"match"`
	Value        string `json:"value"`
}
