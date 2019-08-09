package DB

type PuppetClass struct {
	ID             int    `json:"id"`
	ForemanID      int    `json:"foreman_id"`
	Class          string `json:"class"`
	Subclass       string `json:"subclass"`
	EnvironmentIDs []int  `json:"environment_ids"`
	SmartClassIDs  []int  `json:"sc_ids"`
}

type Parameters struct {
	ID                   int         `json:"id"`
	Name                 string      `json:"name"`
	ModuleName           string      `json:"module_name"`
	SmartClassParameters []Parameter `json:"smart_class_parameters"`
	EnvironmentsID       []int       `json:"environments_id"`
	HostGroups           []HGList    `json:"hostgroups"`
}

type Parameter struct {
	ForemanID int    `json:"id"`
	Parameter string `json:"parameter"`
}

type HGList struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}
