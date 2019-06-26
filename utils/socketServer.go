package utils

import (
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/gorilla/websocket"
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

var (
	newline = []byte{'\n'}
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub  *Hub
	name string
	id   int
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

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func Serve(cfg *models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if cfg.Web.Logged && cfg.Web.SocketActive {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				Error.Println(err)
			}
			cfg.Web.Socket = conn
			//newID := 0
			//if len(hub.clients) > 0 {
			//	newID = hub.clients[-1]
			//}
			//
			//client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), name: cfg.Api.Username, id:newID}
			//client.hub.register <- client
			//client.send <- []byte("TEST")
			// Allow collection of memory referenced by the caller by doing all work in
			// new goroutines.
			//go client.writePump()
			go func(conn *websocket.Conn) {
				for {
					_, _, err = conn.ReadMessage()
					if err != nil {
						_ = conn.Close()
					}
				}
			}(cfg.Web.Socket)
			//fmt.Println("Client connected")
			//fmt.Println(hub.clients)
		}
	}
}

func BroadCastMsg(cfg *models.Config, msg models.Step) {
	var lock sync.Mutex
	if cfg.Web.Logged && cfg.Web.SocketActive {

		lock.Lock()
		defer lock.Unlock()

		data, _ := json.Marshal(msg)
		p := []byte(data)
		if p != nil {
			//err := conn.SetWriteDeadline(time.Now().Add(writeWait))
			//if err != nil {
			//	Error.Println(err)
			//}
			if err := cfg.Web.Socket.WriteMessage(1, p); err != nil {
				Error.Println(err)
				return
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
