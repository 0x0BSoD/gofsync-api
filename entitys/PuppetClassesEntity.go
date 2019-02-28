package entitys

// PuppetClasses container
type PuppetClasses struct {
	Results  map[string][]*PuppetClass `json:"results"`
	Total    int                       `json:"total"`
	SubTotal int                       `json:"subtotal"`
	Page     int                       `json:"page"`
	PerPage  int                       `json:"per_page"`
	Search   string                    `json:"search"`
}

// PuppetClass structure
type PuppetClass struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PuppetClassName
type PuppetClassName struct {
	ID                   int                               `json:"id"`
	Name                 string                            `json:"name"`
	ModuleName           string                            `json:"module_name"`
	CreatedAt            string                            `json:"created_at"`
	UpdatedAt            string                            `json:"updated_at"`
	SmartVariables       []string                          `json:"smart_variables"`
	SmartClassParameters []*SmartClassParameter              `json:"smart_class_parameters"`
	Environments         []*Env          `json:"environments"`
	HostGroups           []*HG            `json:"hostgroups"`
}

// SmartClassParameter
type SmartClassParameter struct {
	Parameter string `json:"parameter"`
	ID        int    `json:"id"`
}

// SmartClassParameter
type HG struct {
	Name string `json:"name"`
	Title string `json:"title"`
	ID        int    `json:"id"`
}

// SmartClassParameter
type Env struct {
	Name string `json:"name"`
	ID        int    `json:"id"`
}
