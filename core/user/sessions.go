package user

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sort"
	"time"
)

func (ss *Sessions) Check(token string) bool {
	if val, ok := ss.Hub[token]; ok {
		fmt.Println("Old Session:")
		fmt.Println(val)
		return true
	} else {
		fmt.Println("New Session:")
		return false
	}
}

func (ss *Sessions) Get(token string) Session {
	return ss.Hub[token]
}

func (ss *Sessions) Set(token string, ctx *GlobalCTX) {
	fmt.Println("Setting session to ctx ....")
	ctx.Session = ss.Get(token)
}

func (ss *Sessions) Add(user *Claims, token string) {
	ID := ss.calcID()
	newSession := Session{
		ID:           ID,
		UserName:     user.Username,
		SocketActive: false,
		TTL:          24 * time.Hour,
		Created:      time.Now(),
		WSMessage:    make(chan []byte),
	}
	ss.Hub[token] = newSession
}

func (ss *Sessions) calcID() int {
	ID := 0
	if len(ss.Hub) > 0 {
		type kv struct {
			Key   string
			Value int
		}
		var sessions []kv
		for k, v := range ss.Hub {
			sessions = append(sessions, kv{k, v.ID})
		}
		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].Value > sessions[j].Value
		})
		ID = sessions[len(sessions)-1].Value
	}
	return ID
}

func (s *Session) AddWSConn(conn *websocket.Conn) {
	fmt.Println("ADD CONN, ", conn)
	s.SocketActive = true
	s.Socket = conn
}

func CreateHub() Sessions {
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
	return Sessions{
		Hub: make(map[string]Session),
	}
}
