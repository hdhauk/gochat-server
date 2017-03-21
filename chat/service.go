package chat

// Session is a connection between the service and a client.
type Session interface {
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
	Close() error
	//LocalAddr() net.Addr
	//RemoteAddr() net.Addr
}

// RoomID is an UUID to identify a room.
type RoomID string

// UserID is an UUID to indentify a user.
type UserID string

// SessionID is an UUID to identify a session.
type SessionID string

// Service host chatrooms, and run them.
type Service struct {
	sessions map[SessionID]Session
	rooms    map[RoomID]Room
}
