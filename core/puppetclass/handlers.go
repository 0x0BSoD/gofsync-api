package puppetclass

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/core/puppetclass/DB"
	"git.ringcentral.com/archops/goFsync/core/smartclass"
	"git.ringcentral.com/archops/goFsync/middleware"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// VARS
	ctx := middleware.GetContext(r)
	params := mux.Vars(r)
	var DBGet DB.Get

	// Get All puppet classes
	puppetClasses := DBGet.All(params["host"], ctx)

	pcObject := make(map[string][]EditorItem)
	for _, pc := range puppetClasses {
		var paramsPC []ParameterItem
		var dumpObj smartclass.SCParameterDef
		for _, scId := range pc.SmartClassIDs {
			scData := smartclass.GetSCData(scId, ctx)
			_ = json.Unmarshal([]byte(scData.Dump), &dumpObj)
			paramsPC = append(paramsPC, ParameterItem{
				ID:             scData.ID,
				ForemanID:      scData.ForemanId,
				Name:           scData.Name,
				DefaultValue:   dumpObj.DefaultValue,
				OverridesCount: scData.OverrideValuesCount,
				Type:           scData.ValueType,
			})
		}
		pcObject[pc.Class] = append(pcObject[pc.Class], EditorItem{
			ID:          pc.ID,
			ForemanID:   pc.ForemanID,
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
