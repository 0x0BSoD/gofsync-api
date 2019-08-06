package DB

import "git.ringcentral.com/archops/goFsync/core/environment"

type PuppetClass struct {
	ID            int    `json:"id"`
	ForemanID     int    `json:"foreman_id"`
	Class         string `json:"class"`
	Subclass      string `json:"subclass"`
	SmartClassIDs []int  `json:"sci_ds"`
}

type Parameters struct {
	ID                   int                       `json:"id"`
	Name                 string                    `json:"name"`
	ModuleName           string                    `json:"module_name"`
	SmartClassParameters []Parameter               `json:"smart_class_parameters"`
	Environments         []environment.Environment `json:"environments"`
	HostGroups           []HGList                  `json:"hostgroups"`
}

type Parameter struct {
	ID        int    `json:"id"`
	Parameter string `json:"parameter"`
}

type HGList struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}
