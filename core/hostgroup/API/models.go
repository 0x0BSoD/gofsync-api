package API

import "git.ringcentral.com/archops/goFsync/core/hostgroup/DB"

type HostGroups struct {
	Results  []DB.HostGroup `json:"results"`
	Total    int            `json:"total"`
	SubTotal int            `json:"subtotal"`
	Page     int            `json:"page"`
	PerPage  int            `json:"per_page"`
	Search   string         `json:"search"`
}

type Parameters struct {
	Results  []DB.HostGroupParameter `json:"results"`
	Total    int                     `json:"total"`
	SubTotal int                     `json:"subtotal"`
	Page     int                     `json:"page"`
	PerPage  int                     `json:"per_page"`
	Search   string                  `json:"search"`
}
