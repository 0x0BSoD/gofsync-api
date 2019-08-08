package environment

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
