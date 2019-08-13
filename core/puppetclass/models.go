package puppetclass

import (
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"sync"
)

// Special web interface oriented structs ====================
type EditorItem struct {
	ID          int             `json:"id"`
	ForemanID   int             `json:"foreman_id"`
	InHostGroup bool            `json:"in_host_group"`
	Class       string          `json:"class"`
	SubClass    string          `json:"sub_class"`
	Parameters  []ParameterItem `json:"parameters"`
}

type ParameterItem struct {
	ID             int    `json:"id"`
	ForemanID      int    `json:"foreman_id"`
	Name           string `json:"name"`
	DefaultValue   string `json:"default_value"`
	Type           string `json:"type"`
	OverridesCount int    `json:"overrides_count"`
}

// Result struct
type PCResult struct {
	sync.Mutex
	resSlice []smartclass.PCSCParameters
}

func (r *PCResult) Add(pc smartclass.PCSCParameters) {
	r.Lock()
	r.resSlice = append(r.resSlice, pc)
	r.Unlock()
}
