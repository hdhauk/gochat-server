package chat

import "time"

// OutgoingMsg is the data structure the server send to the clients.
type OutgoingMsg struct {
	ID         string    `json:"msg-id"`
	SenderName string    `json:"sender-name"`
	SenderID   string    `json:"sender-id"`
	TimeStamp  time.Time `json:"timestamp"`

	IncomingMsg
}
