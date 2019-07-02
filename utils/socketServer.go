package utils

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
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

func WSServe(ctx *user.GlobalCTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if ctx.Session.Socket == nil {

			for v, h := range r.Header {
				fmt.Printf("%s\t%s", v, h)
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				Error.Println(err)
				return
			}
			ctx.Session.AddWSConn(conn)
			fmt.Printf("%s connected\n", ctx.Session.UserName)
			go writePump(ctx.Session.Socket, ctx.Session.WSMessage)
			fmt.Println("WS Pump Started")

		} else {
			fmt.Println("WS skipped")
		}

		fmt.Println("====================================")
		fmt.Println("Session Socket:", ctx.Session.Socket)
		fmt.Println("Session Socket Active:", ctx.Session.SocketActive)
		fmt.Println("Session:", ctx.Session.WSMessage)
		fmt.Println("Config:", ctx.Config)
		fmt.Println("Sessions:", ctx.Sessions)
		fmt.Println("====================================")

	}
}

func CastMsgToUser(ctx *user.GlobalCTX, msg models.Step) {
	strMsg, _ := json.Marshal(msg)
	ctx.Session.WSMessage <- strMsg
}

func writePump(conn *websocket.Conn, msg chan []byte) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = conn.Close()
	}()
	for {
		select {
		case message, ok := <-msg:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(msg)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-msg)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
