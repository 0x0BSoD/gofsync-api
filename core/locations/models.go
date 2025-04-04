package locations

type Locations struct {
	Results  []*Location            `json:"results"`
	Total    int                    `json:"total"`
	SubTotal int                    `json:"subtotal"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
	Search   string                 `json:"search"`
	Sort     map[string]interface{} `json:"sort"`
}
type Location struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// HTTP
type AllLocations struct {
	Host      string `json:"host"`
	Env       string `json:"env"`
	Locations []Loc  `json:"locations"`
	Open      []bool `json:"open"`
}
type Loc struct {
	Name        string `json:"name"`
	Highlighted bool   `json:"highlighted"`
}
