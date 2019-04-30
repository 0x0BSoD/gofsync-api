package puppetclass

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/core/smartclass"
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
		params := mux.Vars(r)

		puppetClasses := GetAllPCDB(cfg, params["host"])

		var pcObject []models.PCintId

		for _, pc := range puppetClasses {
			pcObject = append(pcObject, models.PCintId{
				Class:     pc.Class,
				Subclass:  pc.Subclass,
				ForemanId: pc.ForemanId,
				SCIDs:     pc.SCIDs,
				ID:        pc.ID,
			})
		}

		var res []TreeView
		id := 0
		for _, data := range pcObject {
			var chRes []TreeView
			for _, scId := range data.SCIDs {
				scData := smartclass.GetSCData(scId, cfg)
				chRes = append(chRes, TreeView{
					BaseID: scData.ID,
					Id:     id,
					Name:   scData.Name,
				})
				id++
			}
			res = append(res, TreeView{
				Name:     data.Subclass,
				Id:       id,
				Children: chRes,
			})
			id++
		}

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Error.Printf("Error on getting all locations: %s", err)
		}
	}
}
