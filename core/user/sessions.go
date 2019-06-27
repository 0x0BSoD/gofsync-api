package user

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"sort"
	"time"
)

func CreateHub() models.Sessions {
	//if cfg.Redis {
	//	response, err := cache.Do("GET", "usersHub")
	//	if err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	if response == nil {
	//		w.WriteHeader(http.StatusUnauthorized)
	//		return
	//	}
	//	return response
	//}
	fmt.Println("New HUB")
	return models.Sessions{
		Hub: make(map[string]models.Session),
	}
}

func Start(user *models.Claims, token string, cfg *models.Config) models.Session {

	fmt.Println(cfg.Sessions)

	if val, ok := cfg.Sessions.Hub[token]; ok {
		//fmt.Println("[X] Old Session")
		//fmt.Println(val)
		return val
	} else {
		ID := 0
		if len(cfg.Sessions.Hub) > 0 {
			type kv struct {
				Key   string
				Value int
			}
			var ss []kv
			for k, v := range cfg.Sessions.Hub {
				ss = append(ss, kv{k, v.ID})
			}
			sort.Slice(ss, func(i, j int) bool {
				return ss[i].Value > ss[j].Value
			})
			ID = ss[len(ss)-1].Value
		}
		sa := true
		if user.Username == "srv_foreman" {
			sa = false
		}
		newSession := models.Session{
			ID:           ID,
			UserName:     user.Username,
			SocketActive: sa,
			Config:       cfg,
			TTL:          24 * time.Hour,
			Created:      time.Now(),
			WSMessage:    make(chan []byte),
		}
		//fmt.Println("[X] New Session")
		//fmt.Println(newSession)
		cfg.Sessions.Hub[token] = newSession
		return newSession
	}
}
