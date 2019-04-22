package utils

import (
	"fmt"
	"log"

	"github.com/googollee/go-socket.io"
)

func Serve() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "state", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("event", "have "+msg)
		s.Close()
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go server.Serve()

	return server
}
