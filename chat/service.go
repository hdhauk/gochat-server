package chat

import "errors"

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
	rooms    map[RoomID]RoomInterface

	// Internal communication.
	addRoomCh chan RoomInterface
	stopCh    chan struct{}
}

// NewService returns a new service.
func NewService() *Service {
	stopCh := make(chan struct{})
	close(stopCh)
	return &Service{
		sessions:  make(map[SessionID]Session),
		rooms:     make(map[RoomID]RoomInterface),
		addRoomCh: make(chan RoomInterface),
		stopCh:    stopCh,
	}
}

// Start the chat service.
func (s *Service) Start() {
	// Do nothing if service is already running.
	select {
	case <-s.stopCh:
	default:
		return
	}

	s.stopCh = make(chan struct{})
	// Room Service.
	go func() {
		for {
			select {
			case r := <-s.addRoomCh:
				s.rooms[r.GetID()] = r
			case <-s.stopCh:
				return
			}
		}
	}()
}

// AddRoom opens a new room in the service.
func (s *Service) AddRoom(r RoomInterface) error {
	select {
	case s.addRoomCh <- r:
		return nil
	case <-s.stopCh:
		return ErrServiceNotRunning
	}

}

// Stop the service.
func (s *Service) Stop() {
	select {
	case <-s.stopCh:
		return
	default:
		close(s.stopCh)
	}
}

var (
	// ErrServiceNotRunning is returned whenever someone try to use the service
	// without it running.
	ErrServiceNotRunning = errors.New("service not running")
	// ErrIllegalToken is returned whenever a user try to access any protected
	// endpoints without a valid JWT.
	ErrIllegalToken = errors.New("token not authorized")
)
