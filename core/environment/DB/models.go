package DB

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

type AllEnvSvn struct {
	Info map[string][]SvnInfo `json:"info"`
}
