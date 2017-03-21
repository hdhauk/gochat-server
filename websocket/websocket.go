package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Session is a session between the server and an authenticated user
// over a websocket.
type Session *websocket.Conn

// NewSession returns a new session.
func NewSession(w http.ResponseWriter, r *http.Request) (Session, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
