package puppetclass

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/middleware"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
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
	cfg := middleware.GetConfig(r)
	params := mux.Vars(r)

	puppetClasses := DbAll(params["host"], cfg)

	pcObject := make(map[string][]models.PuppetClassEditor)
	for _, pc := range puppetClasses {
		var paramsPC []models.ParameterEditor
		for _, scId := range pc.SCIDs {
			scData := smartclass.GetSCData(scId, cfg)
			paramsPC = append(paramsPC, models.ParameterEditor{
				ForemanID:      scData.ForemanId,
				Name:           scData.Name,
				DefaultValue:   "",
				OverridesCount: scData.OverrideValuesCount,
				Type:           scData.ValueType,
			})
		}
		pcObject[pc.Class] = append(pcObject[pc.Class], models.PuppetClassEditor{
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
