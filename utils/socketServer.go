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
	// pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	// pingPeriod = (pongWait * 9) / 10
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

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
		lock.Lock()
		data, _ := json.Marshal(msg)
		p := []byte(data)
		if p != nil {
			_ = cfg.Web.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := cfg.Web.Socket.WriteMessage(websocket.TextMessage, p); err != nil {
				return
			}
		}
		defer lock.Unlock()
	}
}
