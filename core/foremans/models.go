package foremans

type ForemanHost struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Env  string `json:"env"`
	//Dashboard DashboardSend `json:"dashboard"`
}

type Container struct {
	Results  []Item `json:"results"`
	Total    int    `json:"total"`
	SubTotal int    `json:"subtotal"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	Search   string `json:"search"`
}

type Item struct {
	ID         int    `json:"id"`
	HostID     int    `json:"host_id"`
	HostName   string `json:"host_name"`
	ReportedAt string `json:"reported_at"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Status     Status `json:"status"`
}

type Status struct {
	Applied        int `json:"applied"`
	Restarted      int `json:"restarted"`
	Failed         int `json:"failed"`
	FailedRestarts int `json:"failed_restarts"`
	Skipped        int `json:"skipped"`
	Pending        int `json:"pending"`
}

type Dashboard struct {
	Trend     map[int]int     `json:"trend"`
	Success   int             `json:"success"`
	RFailed   int             `json:"r_failed"`
	Failed    int             `json:"failed"`
	LastHosts map[string]bool `json:"last_hosts"`
	Summary   int             `json:"summary"`
}
