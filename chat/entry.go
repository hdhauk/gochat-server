package chat

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Entry is one chat message in a chatroom.
type Entry struct {
	ID            string
	Sender        *User
	TimeStamp     time.Time
	Clock         int
	LastClockSeen map[string]int

	IncomingMsg
}

func createEntry(in IncomingMsg) *Entry {
	return &Entry{
		ID:            uuid.NewV4().String(),
		TimeStamp:     time.Now(),
		LastClockSeen: make(map[string]int),
		IncomingMsg:   in,
	}
}
