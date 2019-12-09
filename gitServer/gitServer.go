package gitServer

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type SendMessage struct {
	HostName string `json:"host_name"`
	SWE      string `json:"swe"`
	Action   string `json:"action"`
}

// ClientManager used for storing connected clients
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	lock       *sync.Mutex
}

// Client - it's a client obliviously
type Client struct {
	name   string
	socket net.Conn
	data   chan []byte
}

func (c *ClientManager) start() {
	for {
		select {
		case connection := <-c.register:
			c.clients[connection] = true
			log.Println("[git] added new connection!")
		case connection := <-c.unregister:
			if _, ok := c.clients[connection]; ok {
				close(connection.data)
				delete(c.clients, connection)
				log.Println("[git] a connection has terminated!")
			}
		case message := <-c.broadcast:
			for connection := range c.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					c.lock.Lock()
					delete(c.clients, connection)
					c.lock.Unlock()
				}
			}
		}
	}
}

type gotMessage struct {
	String     string `json:"string"`
	ClientName string `json:"client_name"`
}

func (c *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 8192)
		length, err := client.socket.Read(message)
		if err != nil {
			c.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			messageRight := message[:length]
			log.Println("[git] RECEIVED: " + string(messageRight))
			var data gotMessage
			err := json.Unmarshal([]byte(strings.TrimRight(string(messageRight), "\n")), &data)
			if err != nil {
				panic(err)
			}
			if data.String == "connected" {
				client.name = data.ClientName
			}
			for cl := range c.clients {
				fmt.Println(cl)
			}
			//c.broadcast <- message
		}
	}
}

func (c *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			fmt.Println("sender", string(message))
			client.socket.Write(message)
		}
	}
}

// Clone - clone target SWE
func (c *ClientManager) Cmd(cmd SendMessage) {
	b, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println("error:", err)
	}

	c.lock.Lock()
	fmt.Println("l")
	for cl := range c.clients {
		fmt.Println(cl)
		if cl.name == cmd.HostName {
			fmt.Println(cmd)
			cl.data <- b
		}
	}
	fmt.Println("unl")
	c.lock.Unlock()
}

// StartGitServer - listener for all git clients
func StartGitServer() *ClientManager {
	log.Println("[git] starting server ...")

	listener, err := net.Listen("tcp", ":13666")
	if err != nil {
		fmt.Println(err)
	}

	manager := &ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		lock:       &sync.Mutex{},
	}

	go func() {
		go manager.start()
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
	}()

	return manager
}
