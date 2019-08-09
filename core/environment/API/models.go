package API

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
	ForemanID int    `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	Loading   bool   `json:"loading"`
	Repo      string `json:"repo"`
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
