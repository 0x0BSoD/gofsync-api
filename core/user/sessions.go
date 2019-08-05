package user

import (
	"fmt"
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

func (ss *GlobalCTX) Check(token string) bool {
	if _, ok := ss.Sessions.Hub[token]; ok {
		return true
	} else {
		return false
	}
}

func (ss *GlobalCTX) Set(user *Claims, token string) {
	if val, ok := ss.Sessions.Hub[token]; ok {
		ss.GlobalLock.Lock()
		if ss.Session.UserName != ss.Sessions.Hub[token].UserName {
			ss.Session = &val
		}
		ss.GlobalLock.Unlock()
	} else {
		val := ss.Sessions.add(user, token)
		ss.Session = &val
	}
}

func (ss *Sessions) add(user *Claims, token string) Session {
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

func (ss *GlobalCTX) StartPump() {
	ss.GlobalLock.Lock()
	if !ss.Session.PumpStarted {
		fmt.Println("WS PUMP for", ss.Session.UserName)
		go writePump(ss.Session)
		fmt.Println("Pump Started")

		ss.Session.PumpStarted = true
	}
	ss.GlobalLock.Unlock()
}

func (s *Session) SendMsg(msg []byte) {
	//fmt.Println("Session got message:", string(msg))
	s.WSMessage <- msg
}

func (ss *GlobalCTX) Broadcast(msg []byte) {
	for _, s := range ss.Sessions.Hub {
		s.SendMsg(msg)
	}
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
			fmt.Println("PUMP got message:", string(message))
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
