package user

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/gorilla/websocket"
	"sort"
	"sync"
	"time"
)

const (
	writeWait  = 1 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var newline = []byte{'\n'}

// =====================================================================================================================
// Global Context
// =====================================================================================================================

// Session already exist?
func (s *GlobalCTX) Check(token string) bool {
	if _, ok := s.Sessions.Hub[token]; ok {
		return true
	} else {
		return false
	}
}

// Add a new session or return pointer to existing
func (s *GlobalCTX) Set(user *Claims, token string) {
	if val, ok := s.Sessions.Hub[token]; ok {
		s.GlobalLock.Lock()
		if s.Session.UserName != s.Sessions.Hub[token].UserName {
			s.Session = &val
		}
		s.GlobalLock.Unlock()
	} else {
		s.GlobalLock.Lock()
		val := s.Sessions.add(user, token)
		s.Session = &val
		s.GlobalLock.Unlock()
	}
}

// Send the message to all connected users
func (s *GlobalCTX) Broadcast(wsMessage models.WSMessage) {
	s.GlobalLock.Lock()
	for _, s := range s.Sessions.Hub {
		s.SendMsg(wsMessage)
	}
	s.GlobalLock.Unlock()
}

func (s *GlobalCTX) StartPump(ID int) {
	if !s.Session.Sockets[ID].PumpStarted {
		fmt.Println("starting WS consumer for ", s.Session.UserName, ID)
		go writePump(s.Session.Sockets[ID], s.GlobalLock)
		time.Sleep(1 * time.Second)
	}

	s.Session.Sockets[ID].Lock.Lock()
	s.Session.Sockets[ID].PumpStarted = true
	s.Session.Sockets[ID].Lock.Unlock()
}

// =====================================================================================================================
// Sessions
// =====================================================================================================================

func CreateHub() Sessions {
	return Sessions{
		Hub: make(map[string]Session),
	}
}

func (ss *Sessions) add(user *Claims, token string) Session {
	ID := ss.calcID()
	ss.Hub[token] = Session{
		ID:       ID,
		UserName: user.Username,
		Sockets:  make(map[int]*SocketData),
		Lock:     &sync.Mutex{},
	}
	return ss.Hub[token]
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

// =====================================================================================================================
// Session
// =====================================================================================================================

func (s *Session) Add(conn *websocket.Conn) int {
	s.Lock.Lock()
	ID := s.calcID()
	s.Sockets[ID] = &SocketData{
		PumpStarted: false,
		Socket:      conn,
		WSMessage:   make(chan []byte),
		Lock:        &sync.Mutex{},
	}
	s.Lock.Unlock()
	return ID
}

func (s *Session) SendMsg(wsMessage models.WSMessage) {
	if s != nil {
		s.Lock.Lock()
		defer s.Lock.Unlock()

		for _, s := range s.Sockets {
			if s.PumpStarted {
				msg, err := json.Marshal(wsMessage)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("[WS] ", string(msg))
				s.WSMessage <- msg
			}
		}
	}
}

func (s *Session) calcID() int {
	ID := 0
	if len(s.Sockets) > 0 {
		var IDs []int
		for i := range s.Sockets {
			IDs = append(IDs, i)
		}
		sort.Ints(IDs)
		if IDs != nil {
			last := IDs[len(IDs)-1]
			ID = last + 1
		}

	}
	return ID
}

func writePump(s *SocketData, GlobalLock *sync.Mutex) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		fmt.Println("stopping WS consumer ... ")
		ticker.Stop()
		s.Lock.Lock()
		s.PumpStarted = false
		s.Lock.Lock()
		_ = s.Socket.Close()
		close(s.WSMessage)
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
			s.Lock.Lock()
			_ = s.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
			s.Lock.Unlock()
		}
	}
}
