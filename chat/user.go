package chat

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[CHAT]", log.Ltime|log.Lshortfile)

// User defines one chat user.
type User struct {
	ID          string
	Conn        io.ReadWriter
	DisplayName string
	JWT         string
	Rooms       map[string]*Room

	closeCh chan struct{}
}

// Send transmits to the user
func (u *User) Send(msg OutgoingMsg) error {
	if err := json.NewEncoder(u.Conn).Encode(&msg); err != nil {
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
			if err := json.NewDecoder(u.Conn).Decode(&in); err != nil {
				logger.Printf("[ERROR] Failed to decode incoming message: %s\n", err.Error())
				// NOTE: Should we drop the connection
				continue
			}
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
