package main

// ===============================
// TYPES & VARS
// ===============================
// Smart Class
type SCParameter struct {
	Parameter           string      `json:"parameter"`
	ID                  int         `json:"id"`
	Description         string      `json:"description"`
	Override            bool        `json:"override"`
	ParameterType       string      `json:"parameter_type"`
	DefaultValue        interface{} `json:"default_value"`
	UsePuppetDefault    bool        `json:"use_puppet_default"`
	Required            bool        `json:"required"`
	ValidatorType       string      `json:"validator_type"`
	ValidatorRule       string      `json:"validator_rule"`
	MergeOverrides      bool        `json:"merge_overrides"`
	AvoidDuplicates     bool        `json:"avoid_duplicates"`
	OverrideValueOrder  string      `json:"override_value_order"`
	OverrideValuesCount int         `json:"override_values_count"`
	//CreatedAt           string            `json:"created_at"`
	//UpdatedAt           string            `json:"updated_at"`
	//PuppetClass         *PClass           `json:"puppetclass"`
	//OverrideValues      []*OverrideValues `json:"override_values"`
}

// Smart Class Container
type SCParameters struct {
	Total    int           `json:"total"`
	SubTotal int           `json:"subtotal"`
	Page     int           `json:"page"`
	PerPage  int           `json:"per_page"`
	Search   string        `json:"search"`
	Results  []SCParameter `json:"results"`
}

//// OverrideValues
//type OverrideValues struct {
//	ID               int         `json:"id"`
//	Match            string      `json:"match"`
//	Value            interface{} `json:"value"`
//	UsePuppetDefault bool        `json:"use_puppet_default"`
//}
//
//// PClass
//type PClass struct {
//	Name       string `json:"name"`
//	ModuleMame string `json:"module_name"`
//	ID         int    `json:"id"`
//}

// ===============
// GET
// ===============
//func getSmartClasses(host string) {
//	var result SCParameters
//
//	//aSC := getAllSWE(host)
//
//	for _, sc := range aSC {
//		uri := fmt.Sprintf("smart_class_parameters/%d", sc.SCID)
//		bodyText := ForemanAPI("GET", host, uri, "")
//		err := json.Unmarshal(bodyText, &result)
//		if err != nil {
//			log.Fatalf("%q:\n %s\n", err, bodyText)
//			return
//		}
//		//insertSCOverride(host, result, sc.SCID)
//	}
//}
