package environment

// Environment container
type Environments struct {
	Results  []*Environment         `json:"results"`
	Total    int                    `json:"total"`
	SubTotal int                    `json:"subtotal"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
	Search   string                 `json:"search"`
	Sort     map[string]interface{} `json:"sort"`
}

// Environment structure
type Environment struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Loading bool   `json:"loading"`
	Repo    string
}

// smart_proxies container
type SmartProxies struct {
	Results  []*SmartProxy          `json:"results"`
	Total    int                    `json:"total"`
	SubTotal int                    `json:"subtotal"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
	Search   string                 `json:"search"`
	Sort     map[string]interface{} `json:"sort"`
}

// Environment structure
type SmartProxy struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

// HTTP
type EnvCheckP struct {
	Host string `json:"host"`
	Env  string `json:"env"`
}

// swe update
type SweUpdateParams struct {
	Host        string `json:"host"`
	Environment string `json:"environment"`
	DryRun      bool   `json:"dry_run"`
	Except      string `json:"except"`
}
type SweUpdatePOSTParams struct {
	DryRun bool   `json:"dryrun"`
	Except string `json:"except,omitempty"`
}

type SvnLog struct {
	LogEntry []LogEntry `xml:"logentry" json:"log_entry"`
}
type LogEntry struct {
	Revision string `xml:"revision,attr" json:"revision"`
	Author   string `xml:"author" json:"author"`
	Date     string `xml:"date" json:"date"`
	Msg      string `xml:"msg" json:"msg"`
}

type SvnInfo struct {
	Entry struct {
		Revision    string    `xml:"revision,attr" json:"revision"`
		Kind        string    `xml:"kind,attr" json:"kind"`
		Path        string    `xml:"path,attr" json:"path"`
		Url         string    `xml:"url" json:"url"`
		RelativeUrl string    `xml:"relative-url" json:"relative-url"`
		Repository  SvnRepo   `xml:"repository" json:"repository"`
		WcInfo      SvnWcInfo `xml:"wc-info" json:"wc-info"`
		Commit      SvnCommit `xml:"commit" json:"commit"`
	} `xml:"entry" json:"entry"`
}

type SvnCommit struct {
	Revision string `xml:"revision,attr" json:"revision"`
	Author   string `xml:"author" json:"author"`
	Date     string `xml:"date" json:"date"`
}

type SvnRepo struct {
	Root string `xml:"root" json:"root"`
	UUID string `xml:"uuid" json:"uuid"`
}

type SvnWcInfo struct {
	WcRootAbspath string `xml:"wcroot-abspath" json:"wcroot-abspath"`
	Schedule      string `xml:"schedule" json:"schedule"`
	Depth         string `xml:"depth" json:"depth"`
}
