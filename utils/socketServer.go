package utils

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/middleware"
	"github.com/gorilla/websocket"
	"net/http"
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
)

func WSServe(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.GetContext(r)
	if len(ctx.Sessions.Hub) > 0 {
		if ctx != nil && ctx.Session.UserName != "" {

			ctx.GlobalLock.Lock()
			conn, err := upgrader.Upgrade(w, r, nil)
			ctx.GlobalLock.Unlock()

			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("WS failed: %s", err)))
				return
			}

			ctx.Session.Lock.Lock()
			ctx.Session.Socket = conn
			ctx.Session.Lock.Unlock()

			ctx.StartPump()
		}
	} else {
		fmt.Println(ctx.Sessions.Hub)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("no available sessions in the hub"))
	}
}
