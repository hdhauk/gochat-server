package chat

import (
	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
)

var logger = log.New(os.Stdout, "[CHAT]", log.Ltime|log.Lshortfile)

// ReadWriteJSON defines the interface required for any connection to the  server.
type ReadWriteJSON interface {
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
}

// User defines one chat user.
type User struct {
	ID          string
	Conn        ReadWriteJSON
	DisplayName string
	JWT         string
	Rooms       map[string]*Room

	closeCh chan struct{}
}

// NewUser returns a new user object.
func NewUser(c ReadWriteJSON, name string) *User {
	return &User{
		ID:          uuid.NewV4().String(),
		Conn:        c,
		DisplayName: name,
		Rooms:       make(map[string]*Room),
		closeCh:     make(chan struct{}),
	}
}

// Send transmits to the user
func (u *User) Send(msg OutgoingMsg) error {
	if err := u.Conn.WriteJSON(msg); err != nil {
		return err
	}
	return nil
}

// Listen listen on the connection for incoming messages.
func (u *User) Listen() {
	rx := make(chan IncomingMsg)
	go func() {
		var in IncomingMsg
		for {
			if err := u.Conn.ReadJSON(&in); err != nil {
				logger.Printf("[ERROR] Failed to decode incoming message: %s\n", err.Error())
				// NOTE: Should we drop the connection
				continue
			}
			fmt.Println(in)
			in.SenderID = u.ID
			rx <- in
		}

	}()

	for {
		select {
		case <-u.closeCh:
			logger.Printf("[INFO] Stopping listening for user: %s\n", u.DisplayName)
			return
		case in := <-rx:
			if room, ok := u.Rooms[in.RoomID]; ok {
				// Send message to room.
				room.rx <- in
				continue
			}
			logger.Println("[ERROR] User attempted to send to non-existing room.")
		}
	}
}
