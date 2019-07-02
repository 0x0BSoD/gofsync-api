package utils

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
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
	ctx := middleware.GetContext(r)

	fmt.Println("====================================")
	fmt.Println("Session", ctx.Session)
	fmt.Println("Config", ctx.Config)
	fmt.Println("Sessions", ctx.Sessions)
	fmt.Println("====================================")

	if ctx.Session.Socket == nil {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			Error.Println(err)
			return
		}
		ctx.Session.AddWSConn(conn)
		fmt.Printf("%s connected\n", ctx.Session.UserName)
	} else {
		fmt.Println("WS skipped")
	}
}

func CastMsgToUser(ctx *user.GlobalCTX, msg models.Step) {
	if ctx.Session.Socket != nil && ctx.Session.SocketActive {
		strMsg, _ := json.Marshal(msg)
		ctx.Session.WSMessage <- strMsg
	}
}

func writePump(ctx *user.GlobalCTX) {
	fmt.Println("WS Pump Started")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = ctx.Session.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-ctx.Session.WSMessage:
			_ = ctx.Session.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = ctx.Session.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ctx.Session.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(ctx.Session.WSMessage)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-ctx.Session.WSMessage)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = ctx.Session.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ctx.Session.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
