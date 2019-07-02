package puppetclass

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/middleware"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type PCHttp struct {
	Subclass     string   `json:"subclass"`
	SmartClasses []string `json:"smart_classes,omitempty"`
}

type TreeView struct {
	BaseID   int        `json:"base_id,omitempty"`
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Children []TreeView `json:"children,omitempty"`
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)

	puppetClasses := DbAll(params["host"], ctx)

	pcObject := make(map[string][]PuppetClassEditor)
	for _, pc := range puppetClasses {
		var paramsPC []ParameterEditor
		var dumpObj smartclass.SCParameterDef
		for _, scId := range pc.SCIDs {
			scData := smartclass.GetSCData(scId, ctx)
			_ = json.Unmarshal([]byte(scData.Dump), &dumpObj)
			paramsPC = append(paramsPC, ParameterEditor{
				ID:             scData.ID,
				ForemanID:      scData.ForemanId,
				Name:           scData.Name,
				DefaultValue:   dumpObj.DefaultValue,
				OverridesCount: scData.OverrideValuesCount,
				Type:           scData.ValueType,
			})
		}
		pcObject[pc.Class] = append(pcObject[pc.Class], PuppetClassEditor{
			ID:          pc.ID,
			ForemanID:   pc.ForemanId,
			Class:       pc.Class,
			SubClass:    pc.Subclass,
			InHostGroup: false,
			Parameters:  paramsPC,
		})
	}

	err := json.NewEncoder(w).Encode(pcObject)
	if err != nil {
		logger.Error.Printf("Error on getting all locations: %s", err)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	Sync(params["host"], ctx)
	err := json.NewEncoder(w).Encode("submitted")
	if err != nil {
		logger.Error.Printf("Error on EnvCheck: %s", err)
	}
}
