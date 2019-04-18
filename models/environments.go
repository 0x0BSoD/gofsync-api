package models

// PuppetClasses container
type Environments struct {
	Results  []*Environment         `json:"results"`
	Total    int                    `json:"total"`
	SubTotal int                    `json:"subtotal"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
	Search   string                 `json:"search"`
	Sort     map[string]interface{} `json:"sort"`
}

// PuppetClass structure
type Environment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// HTTP
type EnvCheckP struct {
	Host string `json:"host"`
	Env  string `json:"env"`
}
