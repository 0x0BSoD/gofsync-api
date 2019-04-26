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
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Children []TreeView `json:"children,omitempty"`
}

//type TreeViewRoot struct {
//	Items []TreeViewChildren `json:"items"`
//	Search string `json:"search"`
//	CaseSensitive bool `json:"caseSensitive"`
//}

func GetAllPCHttp(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var res []TreeView
		mapRes := make(map[string][]string)
		puppetClasses := GetAllPCDB(cfg)
		for _, pc := range puppetClasses {
			var children []string
			for _, id := range pc.SCIDs {
				scData := smartclass.GetSCData(id, cfg)
				children = append(children, scData.Name)
			}
			mapRes[pc.Class] = children
		}

		id := 0
		for class, subclasses := range mapRes {
			var tmpRes []TreeView
			for _, subclass := range subclasses {
				tmpRes = append(tmpRes, TreeView{
					Id:   id,
					Name: subclass,
				})
				id++
			}
			id++
			res = append(res, TreeView{
				Id:       id,
				Name:     class,
				Children: tmpRes,
			})
		}

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Error.Printf("Error on getting all locations: %s", err)
		}
	}
}
