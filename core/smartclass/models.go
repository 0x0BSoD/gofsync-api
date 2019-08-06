package smartclass

//
//import (
//	"git.ringcentral.com/archops/goFsync/core/environment"
//)
//
//// Smart Class Container

//
//// Smart Class
//
//
//// Smart Class with string def parameter
//type SCParameterDef struct {
//	Parameter           string          `json:"parameter"`
//	PuppetClass         PuppetClassInSc `json:"puppetclass"`
//	ID                  int             `json:"id"`
//	Description         string          `json:"description"`
//	Override            bool            `json:"override"`
//	ParameterType       string          `json:"parameter_type"`
//	DefaultValue        string          `json:"default_value"`
//	UsePuppetDefault    bool            `json:"use_puppet_default"`
//	Required            bool            `json:"required"`
//	ValidatorType       string          `json:"validator_type"`
//	ValidatorRule       string          `json:"validator_rule"`
//	MergeOverrides      bool            `json:"merge_overrides"`
//	AvoidDuplicates     bool            `json:"avoid_duplicates"`
//	OverrideValueOrder  string          `json:"override_value_order"`
//	OverrideValuesCount int             `json:"override_values_count"`
//	OverrideValues      []OverrideValue `json:"override_values"`
//}
//
//// PC for old Foreman
//type PuppetClassInSc
//
//// OverrideValues Container

//
//type PCSCParameters struct {
//	ID                   int                       `json:"id"`
//	Name                 string                    `json:"name"`
//	ModuleName           string                    `json:"module_name"`
//	SmartClassParameters []PCSCParameter           `json:"smart_class_parameters"`
//	Environments         []environment.Environment `json:"environments"`
//	HostGroups           []HGList                  `json:"hostgroups"`
//}
//type HGList struct {
//	ID    int    `json:"id"`
//	Name  string `json:"name"`
//	Title string `json:"title"`
//}
//
//type PCSCParameter struct {
//	ID        int    `json:"id"`
//	Parameter string `json:"parameter"`
//}
//
//// Return From Base

//type LogStatus struct {
//	Name          string `json:"name"`
//	Host          string `json:"host"`
//	Current       int    `json:"current"`
//	CurrentThread int    `json:"current_thread,omitempty"`
//	TotalInThread int    `json:"total_in_thread,omitempty"`
//	Total         int    `json:"total"`
//}
//
//type OvrParams struct {
//	OvrForemanId   int    `json:"ovr_foreman_id"`
//	SCForemanId    int    `json:"sc_foreman_id"`
//	DefaultValue   string `json:"default_value"`
//	Type           string `json:"type"`
//	PuppetClass    string `json:"puppet_class"`
//	SmartClassName string `json:"smart_class_name"`
//	Value          string `json:"value"`
//}
//
//type OverrideParameters struct {
//	PuppetClass string              `json:"puppet_class"`
//	Parameters  []OverrideParameter `json:"parameters"`
//}
//
//type OverrideParameter struct {
//	ParameterForemanId int    `json:"parameter_foreman_id"`
//	OverrideForemanId  int    `json:"override_foreman_id"`
//	Name               string `json:"name"`
//	Value              string `json:"value"`
//	Type               string `json:"type"`
//	DefaultValue       string `json:"default_value"`
//}
