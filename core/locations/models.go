package locations

import "git.ringcentral.com/archops/goFsync/core/locations/info"

// HTTP
type AllLocations struct {
	Host      string `json:"host"`
	Env       string `json:"env"`
	Locations []Loc  `json:"locations"`
	Open      []bool `json:"open"`
	info.Dashboard
}
type Loc struct {
	Name        string `json:"name"`
	Highlighted bool   `json:"highlighted"`
}
