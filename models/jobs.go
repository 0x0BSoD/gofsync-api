package models

type Job struct {
	Action     string `json:"action"`
	SourceHost string `json:"source_host,omitempty"`
	TargetHost string `json:"target_host,omitempty"`
	HostGroup  string `json:"host_group,omitempty"`
	User       string `json:"user"`
}
