package environment

import "encoding/xml"

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
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	State     string      `json:"state"`
	Loading   bool        `json:"loading"`
	WSMessage interface{} `json:"ws_message"`
	Repo      string
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

type AllEnvSvn struct {
	Info map[string][]SvnDirInfo `json:"info"`
}

type NewEnvParams struct {
	Name         string `json:"name"`
	LocationsIDs []int  `json:"location_ids"`
}
type NewEnv struct {
	Environment NewEnvParams `json:"environment"`
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

type SvnUrlInfo struct {
	XMLName xml.Name `xml:"info" json:"xml_name"`
	Text    string   `xml:",chardata" json:"text"`
	Entry   struct {
		Text       string `xml:",chardata" json:"text"`
		Kind       string `xml:"kind,attr" json:"kind"`
		Path       string `xml:"path,attr" json:"path"`
		Revision   string `xml:"revision,attr" json:"revision"`
		URL        string `xml:"url" json:"url"`
		Repository struct {
			Text string `xml:",chardata" json:"text"`
			Root string `xml:"root" json:"root"`
			Uuid string `xml:"uuid" json:"uuid"`
		} `xml:"repository" json:"repository"`
		Commit struct {
			Text     string `xml:",chardata" json:"text"`
			Revision string `xml:"revision,attr" json:"revision"`
			Author   string `xml:"author" json:"author"`
			Date     string `xml:"date" json:"date"`
		} `xml:"commit" json:"commit"`
	} `xml:"entry" json:"entry"`
}

type SvnDirInfo struct {
	XMLName xml.Name `xml:"info" json:"xml_name"`
	Text    string   `xml:",chardata" json:"text"`
	Entry   struct {
		Text       string `xml:",chardata" json:"text"`
		Kind       string `xml:"kind,attr" json:"kind"`
		Path       string `xml:"path,attr" json:"path"`
		Revision   string `xml:"revision,attr" json:"revision"`
		URL        string `xml:"url" json:"url"`
		Repository struct {
			Text string `xml:",chardata" json:"text"`
			Root string `xml:"root" json:"root"`
			Uuid string `xml:"uuid" json:"uuid"`
		} `xml:"repository" json:"repository"`
		Commit struct {
			Text     string `xml:",chardata" json:"text"`
			Revision string `xml:"revision,attr" json:"revision"`
			Author   string `xml:"author" json:"author"`
			Date     string `xml:"date" json:"date"`
		} `xml:"commit" json:"commit"`
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
