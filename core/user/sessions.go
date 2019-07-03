package user

import (
	"github.com/gorilla/websocket"
	"sort"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 1 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var newline = []byte{'\n'}

func (ss *Sessions) Check(token string) bool {
	if _, ok := ss.Hub[token]; ok {
		return true
	} else {
		return false
	}
}

//func (ss *Sessions) Get(token string) Session {
//	return ss.Hub[token]
//}

func (ss *Sessions) Get(user *Claims, token string) Session {
	if val, ok := ss.Hub[token]; ok {
		return val
	} else {
		return ss.Add(user, token)
	}
}

func (ss *Sessions) Add(user *Claims, token string) Session {
	ID := ss.calcID()
	newSession := Session{
		ID:          ID,
		UserName:    user.Username,
		PumpStarted: false,
		TTL:         24 * time.Hour,
		Created:     time.Now(),
		WSMessage:   make(chan []byte),
	}
	ss.Hub[token] = newSession
	return newSession
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
	s.Socket = conn
}

func (s *Session) StartWSPump() {
	go writePump(s)
	s.PumpStarted = true
}

func writePump(s *Session) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = s.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-s.WSMessage:
			_ = s.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = s.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := s.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(s.WSMessage)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-s.WSMessage)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = s.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
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
