package puppetclass

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
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

func madeChildren(cfg *models.Config, obj []models.PCintId, idMask *int) []TreeView {
	var res []TreeView
	for _, i := range obj {
		var chRes []TreeView
		for _, scId := range i.SCIDs {
			scData := smartclass.GetSCData(scId, cfg)
			chRes = append(chRes, TreeView{
				BaseID: scData.ID,
				Id:     *idMask,
				Name:   scData.Name,
			})
			*idMask++
		}
		res = append(res, TreeView{
			BaseID:   i.ID,
			Id:       *idMask,
			Name:     i.Subclass,
			Children: chRes,
		})
		*idMask++
	}
	return res
}

func GetAllPCHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		puppetClasses := GetAllPCDB(cfg)

		pcObject := make(map[string][]models.PCintId)

		for _, pc := range puppetClasses {
			pcObject[pc.Class] = append(pcObject[pc.Class], models.PCintId{
				Class:     pc.Class,
				Subclass:  pc.Subclass,
				ForemanId: pc.ForemanId,
				SCIDs:     pc.SCIDs,
				ID:        pc.ID,
			})
		}

		var res []TreeView
		id := 0
		for pc, data := range pcObject {
			ch := madeChildren(cfg, data, &id)
			res = append(res, TreeView{
				Name:     pc,
				Id:       id,
				Children: ch,
			})
			id++
		}

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Error.Printf("Error on getting all locations: %s", err)
		}
	}
}
