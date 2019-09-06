package puppetclass

import (
	"git.ringcentral.com/archops/goFsync/core/smartclass"
)

// PuppetClasses container
type PuppetClasses struct {
	Results  map[string][]PuppetClass `json:"results"`
	Total    int                      `json:"total"`
	SubTotal int                      `json:"subtotal"`
	Page     int                      `json:"page"`
	PerPage  int                      `json:"per_page"`
	Search   string                   `json:"search"`
}

// PuppetClass structure
type PuppetClass struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	SmartClassesId []int
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// PuppetclassesNI for getting from base
type PuppetclassesNI struct {
	Class     string
	SubClass  string
	ForemanID int
}

type PC struct {
	ID        int    `json:"id,omitempty"`
	ForemanId int    `json:"foreman_id,omitempty"`
	Class     string `json:"class"`
	Subclass  string `json:"subclass"`
	SCIDs     string `json:"sci_ds"`
}

type PCintId struct {
	ID        int    `json:"id,omitempty"`
	ForemanId int    `json:"foreman_id,omitempty"`
	Class     string `json:"class"`
	Subclass  string `json:"subclass"`
	SCIDs     []int  `json:"sci_ds,omitempty"`
}

type PuppetClassesWeb struct {
	Subclass     string                  `json:"subclass"`
	SmartClasses []smartclass.SmartClass `json:"smart_classes,omitempty"`
	Overrides    []smartclass.SCOParams  `json:"overrides,omitempty"`
}

// Type for an editor ====================
type PuppetClassesEditor map[int]PuppetClassEditor

type PuppetClassEditor struct {
	ID          int               `json:"id"`
	ForemanID   int               `json:"foreman_id"`
	InHostGroup bool              `json:"in_host_group"`
	Class       string            `json:"class"`
	SubClass    string            `json:"sub_class"`
	Parameters  []ParameterEditor `json:"parameters"`
}
type ParameterEditor struct {
	ID             int    `json:"id"`
	ForemanID      int    `json:"foreman_id"`
	Name           string `json:"name"`
	DefaultValue   string `json:"default_value"`
	Type           string `json:"type"`
	OverridesCount int    `json:"overrides_count"`
}
