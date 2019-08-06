package puppetclass

// Special web interface oriented structs ====================
type EditorItem struct {
	ID          int             `json:"id"`
	ForemanID   int             `json:"foreman_id"`
	InHostGroup bool            `json:"in_host_group"`
	Class       string          `json:"class"`
	SubClass    string          `json:"sub_class"`
	Parameters  []ParameterItem `json:"parameters"`
}

type ParameterItem struct {
	ID             int    `json:"id"`
	ForemanID      int    `json:"foreman_id"`
	Name           string `json:"name"`
	DefaultValue   string `json:"default_value"`
	Type           string `json:"type"`
	OverridesCount int    `json:"overrides_count"`
}
