package utils

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// For DEV ===
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Serve(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		cfg.Web.Socket = conn
		if err != nil {
			log.Println(err)
			return
		}
	}
}
func BroadCastMsg(cfg *models.Config, msg models.Step) {
	var lock sync.Mutex
	if cfg.Web.Logged {
		data, _ := json.Marshal(msg)
		p := []byte(data)
		if p != nil {
			lock.Lock()
			_ = cfg.Web.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := cfg.Web.Socket.WriteMessage(websocket.TextMessage, p); err != nil {
				return
			}
			lock.Unlock()
		}
	}
}
