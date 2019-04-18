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

// PuppetclassesNI for getting from base
type PuppetclassesNI struct {
	Class     string
	SubClass  string
	ForemanID int
}

type PC struct {
	ID        int
	ForemanId int
	Class     string
	Subclass  string
	SCIDs     string
}
type PuppetClassesWeb struct {
	Subclass     string      `json:"subclass"`
	SmartClasses []string    `json:"smart_classes,omitempty"`
	Overrides    []SCOParams `json:"overrides,omitempty"`
}
