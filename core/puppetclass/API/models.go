package API

type PuppetClasses struct {
	Results  map[string][]PuppetClass `json:"results"`
	Total    int                      `json:"total"`
	SubTotal int                      `json:"subtotal"`
	Page     int                      `json:"page"`
	PerPage  int                      `json:"per_page"`
	Search   string                   `json:"search"`
}

type PuppetClass struct {
	ForemanID     int    `json:"id"`
	Name          string `json:"name"`
	SmartClassIDs []int  `json:"smart_class_ids"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type PuppetClassDetailed struct {
	ForemanID    int    `json:"id"`
	Name         string `json:"name"`
	ModuleName   string `json:"module_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Environments []struct {
		ForemanID int    `json:"id"`
		Name      string `json:"name"`
	} `json:"environments"`
	SmartClassParameters []struct {
		ForemanID int    `json:"id"`
		Parameter string `json:"parameter"`
	} `json:"smart_class_parameters"`
}
