package models

// PuppetClasses container
type PuppetClasses struct {
	Results  map[string][]PuppetClass `json:"results"`
	Total    int                      `json:"total"`
	SubTotal int                      `json:"subtotal"`
	Page     int                      `json:"page"`
	PerPage  int                      `json:"per_page"`
	Search   string                   `json:"search"`
}

// PuppetClass structure
type PuppetClass struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	SmartClassesId []int
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
