package main

import (
	"encoding/json"
)

// ===============================
// TYPES & VARS
// ===============================
// PuppetClasses container
type Environments struct {
	Results  []*Environment         `json:"results"`
	Total    int                    `json:"total"`
	SubTotal int                    `json:"subtotal"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
	Search   string                 `json:"search"`
	Sort     map[string]interface{} `json:"sort"`
}

// PuppetClass structure
type Environment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ===============
// GET
// ===============
func environments(host string) (Environments, error) {

	var result Environments
	bodyText := ForemanAPI("GET", host, "environments", "")

	err := json.Unmarshal(bodyText, &result)
	if err != nil {
		return Environments{}, err
	}

	return result, nil
}
