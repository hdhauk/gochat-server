package chat

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Room ...
type Room struct {
	ID        string
	Name      string
	Topic     string
	Members   map[string]*User
	Key       string
	History   map[string]*Entry
	NextClock int

	rx      chan IncomingMsg
	closeCh chan struct{}
}

// NewRoom returns a new bearbone room.
func NewRoom(name string) *Room {
	return &Room{
		ID:      uuid.NewV4().String(),
		Name:    name,
		Members: make(map[string]*User),
		History: make(map[string]*Entry),
		rx:      make(chan IncomingMsg),
		closeCh: make(chan struct{}),
	}
}

// Start initalized the chatroom and listen for messages, and broadcasts to all
// participants
func (r *Room) Start() {
	logger.Printf("[INFO] Starting room: %s\n", r.Name)
	for {
		select {
		case in := <-r.rx:
			// Determine sender identity
			senderID := in.SenderID
			sender, ok := r.Members[senderID]
			if !ok {
				logger.Println("[WARN] Sender not a member of the channel. Discarding entry.")
				continue
			}

			// Forge entry
			entry := createEntry(in)
			entry.Clock = r.NextClock
			r.NextClock++
			entry.RoomID = r.ID
			entry.Sender = sender

			// Save entry to room history
			r.History[entry.ID] = entry

			// Create outgoing message
			out := OutgoingMsg{
				ID:          entry.ID,
				SenderName:  sender.DisplayName,
				SenderID:    sender.ID,
				TimeStamp:   entry.TimeStamp,
				IncomingMsg: in,
			}

			// Broadcast to all room participants.
			for _, member := range r.Members {
				member.Send(out)
			}

		case <-r.closeCh:
			logger.Printf("[INFO] Closing room: %s\n", r.Name)
			return
		}
	}
}

// AddMember adds user to the room, unless it's already a member.
func (r *Room) AddMember(user *User) error {
	_, added := r.Members[user.ID]
	if added {
		err := fmt.Errorf("user id already registered in room (id = %s )", user.ID)
		logger.Printf("[INFO] Cannot add user: %s\n ", err.Error())
		return err
	}
	logger.Printf("[INFO] User \"%s\" added to room \"%s\"\n ", user.DisplayName, r.Name)
	r.Members[user.ID] = user

	// Add room to user
	user.Rooms[r.ID] = r

	return nil
}
