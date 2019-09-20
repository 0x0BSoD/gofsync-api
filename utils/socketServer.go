package utils

import (
	"encoding/json"
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

func WSSessions(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.GetContext(r)

	//type JSONSess struct {
	//	ID int `json:"id"`
	//	Name string `json:"name"`
	//	Sockets []struct{
	//		ID int `json:"id"`
	//		PumpStarted bool `json:"pump_started"`
	//	} `json:"sockets"`
	//}

	for _, s := range ctx.Sessions.Hub {
		fmt.Println("SessionID:", s.ID)
		fmt.Println("User:", s.UserName)
		for _, ss := range s.Sockets {
			fmt.Println("SocketID:", ss.ID)
			fmt.Println("MsgSender:", ss.PumpStarted)
		}
	}

	err := json.NewEncoder(w).Encode("data")
	if err != nil {
		fmt.Printf("Error on getting sessions info : %s", err)
	}
}

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

			ctx.GlobalLock.Lock()
			ID := ctx.Session.Add(conn)
			ctx.GlobalLock.Unlock()

			ctx.GlobalLock.Lock()
			ctx.StartPump(ID)
			ctx.GlobalLock.Unlock()
		}
	} else {
		fmt.Println(ctx.Sessions.Hub)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("no available sessions in the hub"))
	}
}
