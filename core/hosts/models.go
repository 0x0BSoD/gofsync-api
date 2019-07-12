package hosts

type Hosts struct {
	Results  []Host `json:"results"`
	Total    int    `json:"total"`
	SubTotal int    `json:"subtotal"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	Search   string `json:"search"`
}

type Host struct {
	Name                 string `json:"name"`
	ForemanID            int    `json:"id"`
	EnvironmentForemanID int    `json:"environment_id"`
	LocationForemanID    int    `json:"location_id"`
	HostGroupForemanID   int    `json:"hostgroup_id"`
	Environment          string `json:"environment_name"`
	LocationName         string `json:"location_name"`
	HostGroup            string `json:"hostgroup_name"`
	IP                   string `json:"ip"`
	LastReport           string `json:"last_report"`
	MAC                  string `json:"mac"`
	DomainName           string `json:"domain_name"`
	OperatingSystem      string `json:"operatingsystem_name"`
	ModelName            string `json:"model_name"`
}

type ForemanHost struct {
	Name string `json:"name"`
	Env  string `json:"env"`
}
