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
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SmartClassIDs []int  `json:"smart_class_ids"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
