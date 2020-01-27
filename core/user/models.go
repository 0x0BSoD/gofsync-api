package user

import (
	"git.ringcentral.com/archops/goFsync/gitServer"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"sync"
)

type GlobalCTX struct {
	Sessions   Sessions
	Session    *Session
	Config     models.Config
	GitSrv     *gitServer.ClientManager
	GlobalLock *sync.Mutex
	SyncWIP    bool
}

type Credentials struct {
	Password   string `json:"password"`
	Username   string `json:"username"`
	RememberMe bool   `json:"remember_me,omitempty"`
}
type Claims struct {
	Username   string `json:"username"`
	RememberMe bool   `json:"remember_me,omitempty"`
	jwt.StandardClaims
}
type User struct {
	UUID     string `json:"uuid" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// Struct for store user sessions, key it user token
type Sessions struct {
	Hub map[string]Session
}

type Session struct {
	ID       int
	UserName string
	Lock     *sync.Mutex
	Sockets  map[int]*SocketData
}

type SocketData struct {
	ID          int
	PumpStarted bool
	WSMessage   chan []byte
	Lock        *sync.Mutex
	Socket      *websocket.Conn
}
