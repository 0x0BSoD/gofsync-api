package utils

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/middleware"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/gorilla/websocket"
	"net/http"
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

var (
	newline  = []byte{'\n'}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// For DEV ===
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WSServe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := middleware.GetConfig(r)
	if cfg.SocketActive && cfg.Socket == nil {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			Error.Println(err)
			return
		}
		cfg.Socket = conn
		go writePump(&cfg)
		fmt.Printf("%s connected\n", cfg.UserName)
	} else {
		fmt.Println("WS skipped")
	}
}

func CastMsgToUser(ss *models.Session, msg models.Step) {
	if ss.SocketActive {
		strMsg, _ := json.Marshal(msg)
		ss.WSMessage <- strMsg
	}
}

func writePump(ss *models.Session) {
	fmt.Println("WS Pump Started")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ss.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-ss.WSMessage:
			ss.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				ss.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ss.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(ss.WSMessage)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-ss.WSMessage)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ss.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ss.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
