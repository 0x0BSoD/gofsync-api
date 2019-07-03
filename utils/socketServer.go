package utils

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/middleware"
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
	if ctx != nil {

		fmt.Println(ctx.Session.PumpStarted)
		fmt.Println(ctx.Session.Socket)
		fmt.Println(ctx.Session.UserName)

		if ctx.Session.Socket == nil {

			for v, h := range r.Header {
				fmt.Printf("%s\t\t\t%s\n", v, h)
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				Error.Println(err)
				return
			}

			ctx.Session.AddWSConn(conn)

			fmt.Println("====================================")
			fmt.Printf("%s connected\n", ctx.Session.UserName)
			fmt.Println("Session Socket:", ctx.Session.Socket)
			fmt.Println("Session Socket Active:", ctx.Session.PumpStarted)
			fmt.Println("Session message channel:", ctx.Session.WSMessage)
			fmt.Println("Config:", ctx.Config)
			fmt.Println("Sessions:", ctx.Sessions)
			fmt.Println("====================================")

		} else {
			fmt.Println("WS skipped")
		}
	}
}
