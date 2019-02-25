package entitys

// SWEs container
type SWEs struct {
	Hostgroup *SWE `json:"hostgroup"`
}

// SWE structure
type SWE struct {
	Name              string                 `json:"name"`
	ID                int                    `json:"id"`
	SubnetID          int                    `json:"subnet_id"`
	OperatingsystemID int                    `json:"operatingsystem_id"`
	DomainID          int                    `json:"domain_id"`
	EnvironmentID     int                    `json:"environment_id"`
	Ancestry          string                 `json:"ancestry,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	PuppetclassIds    []int                  `json:"puppetclass_ids"`
}
