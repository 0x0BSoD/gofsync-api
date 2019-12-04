package foremans

type ForemanHost struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Env       string    `json:"env"`
	Dashboard Dashboard `json:"dashboard"`
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

type Trend struct {
	Labels []int `json:"labels"`
	Values []int `json:"values"`
}

type Dashboard struct {
	Trend    Trend  `json:"trend"`
	Success  int    `json:"success"`
	RFailed  int    `json:"r_failed"`
	Failed   int    `json:"failed"`
	LastHost string `json:"last_host"`
	Summary  int    `json:"summary"`
}
