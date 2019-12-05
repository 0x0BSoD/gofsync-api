package utils

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"net"
	"strings"
)

// ClientManager used for storing connected clients
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	ctx        *user.GlobalCTX
}

// Client - it's a client obliviously
type Client struct {
	name   string
	socket net.Conn
	data   chan []byte
}

func (manager *ClientManager) start(ctx *user.GlobalCTX) {
	manager.ctx = ctx
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			Info.Println("[git] added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				Info.Println("[git] a connection has terminated!")
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

type gotMessage struct {
	String     string `json:"string"`
	ClientName string `json:"client_name"`
}

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			messageRight := message[:length]
			Info.Println("[git] RECEIVED: " + string(messageRight))
			var data gotMessage
			err := json.Unmarshal([]byte(strings.TrimRight(string(messageRight), "\n")), &data)
			if err != nil {
				panic(err)
			}
			if data.String == "connected" {
				client.name = data.ClientName
			}
			for cl := range manager.clients {
				fmt.Println(cl)
			}
			manager.broadcast <- message
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

// StartGitServer - listener for all git clients
func StartGitServer(ctx *user.GlobalCTX) {
	Info.Println("[git] starting server ...")
	listener, err := net.Listen("tcp", ":13666")
	if err != nil {
		fmt.Println(err)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start(ctx)
	for {
		connection, _ := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}
