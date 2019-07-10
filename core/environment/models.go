package environment

// Environment container
type Environments struct {
	Results  []*Environment         `json:"results"`
	Total    int                    `json:"total"`
	SubTotal int                    `json:"subtotal"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
	Search   string                 `json:"search"`
	Sort     map[string]interface{} `json:"sort"`
}

// Environment structure
type Environment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// smart_proxies container
type SmartProxies struct {
	Results  []*SmartProxy          `json:"results"`
	Total    int                    `json:"total"`
	SubTotal int                    `json:"subtotal"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
	Search   string                 `json:"search"`
	Sort     map[string]interface{} `json:"sort"`
}

// Environment structure
type SmartProxy struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

// HTTP
type EnvCheckP struct {
	Host string `json:"host"`
	Env  string `json:"env"`
}

// swe update
type SweUpdateParams struct {
	Host        string `json:"host"`
	Environment string `json:"environment"`
	DryRun      bool   `json:"dry_run"`
	Except      string `json:"except"`
}
type SweUpdatePOSTParams struct {
	DryRun bool   `json:"dryrun"`
	Except string `json:"except,omitempty"`
}
