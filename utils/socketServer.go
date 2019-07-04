package utils

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// For DEV ===
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mutex = &sync.Mutex{}
)

func WSServe(ctx *user.GlobalCTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx != nil {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("WS failed: %s", err)))
				return
			}
			ctx.Session.AddWSConn(conn)
			mutex.Lock()
			ctx.Session.PumpStarted = true
			mutex.Unlock()
		}
	}
}
