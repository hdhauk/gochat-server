package chat

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type room interface {
	// Authorization
	AddUser(UserID) error
	KickUser(UserID) error

	// Session handeling
	AddSession(UserID, Session) error
	RemoveSession(UserID, Session) error

	// Information retrival
	ListUsers() map[UserID]AuthLevel
	//GetHistory() []entry

	// Management
	Start()
	Stop()
	Reset()
}

// Room is a chatroom.
type Room struct {
	// Room details
	ID    RoomID
	Name  string
	Topic string

	// Users
	AuthUsers   map[UserID]AuthLevel
	ActiveUsers map[UserID]Session

	// Internal communication
	newSessionCh  chan newSession
	newUserCh     chan newUser
	kickUserCh    chan newUser
	userListCh    chan userList
	toBroadcastCh chan interface{}
	stopCh        chan struct{}
}

// AuthLevel determine a users access level.
type AuthLevel int

type newSession struct {
	s   Session
	id  UserID
	err chan<- error
}
type newUser struct {
	id  UserID
	err chan<- error
}
type userList struct {
	u     []UserID
	retCh chan map[UserID]AuthLevel
}

// NewRoom returns a new room object along with its RoomID.
func NewRoom(authUsers map[UserID]AuthLevel, name, topic string) (*Room, RoomID) {
	id := RoomID(uuid.NewV4().String())
	return &Room{
		ID:           RoomID(id),
		Name:         name,
		Topic:        topic,
		AuthUsers:    authUsers,
		ActiveUsers:  make(map[UserID]Session),
		newSessionCh: make(chan newSession),
		newUserCh:    make(chan newUser),
		kickUserCh:   make(chan newUser),
		userListCh:   make(chan userList),
	}, id
}

// Start initialized and start the room, allowing for chat services.
func (r *Room) Start() {
	r.stopCh = make(chan struct{})
	go func() {
		for {
			select {

			// Add Session.
			case new := <-r.newSessionCh:
				if _, ok := r.AuthUsers[new.id]; !ok {
					new.err <- ErrUserNotAuth
					continue
				}
				r.ActiveUsers[new.id] = new.s
				new.err <- nil

			// Add authorized user.
			case new := <-r.newUserCh:
				r.AuthUsers[new.id] = AuthLevel(1)
				new.err <- nil

				// Remove authorized user.
			case kick := <-r.kickUserCh:
				_, ok := r.AuthUsers[kick.id]
				if !ok {
					kick.err <- ErrUserNotFound
					continue
				}
				delete(r.AuthUsers, kick.id)
				kick.err <- nil

			// Dump authorized users.
			case l := <-r.userListCh:
				copy := make(map[UserID]AuthLevel)
				for k, v := range r.AuthUsers {
					copy[k] = v
				}
				l.retCh <- copy

			// Stop service.
			case <-r.stopCh:
				return
			}
		}
	}()
}

// ListUsers returns the ID of all authorized users.
func (r *Room) ListUsers() (map[UserID]AuthLevel, error) {
	retCh := make(chan map[UserID]AuthLevel)
	select {
	case r.userListCh <- userList{retCh: retCh}:
		return <-retCh, nil
	case <-time.After(timeout):
		return nil, ErrRoomNotRunning
	}
}

// Stop the room service.
func (r *Room) Stop() {
	close(r.stopCh)
}

// AddUser adds a user to the list of authorized users.
func (r *Room) AddUser(id UserID) error {
	errCh := make(chan error)
	select {
	case r.newUserCh <- newUser{id: id, err: errCh}:
		return <-errCh
	case <-time.After(timeout):
		return ErrRoomNotRunning
	}
}

// KickUser removes the user from the list of authorized users. Returns an
// error if the user does not exist.
func (r *Room) KickUser(id UserID) error {
	errCh := make(chan error)
	select {
	case r.kickUserCh <- newUser{id: id, err: errCh}:
		return <-errCh
	case <-time.After(timeout):
		return ErrRoomNotRunning
	}
}

// AddSession adds the session to the room.
func (r *Room) AddSession(id UserID, s Session) error {
	errCh := make(chan error)
	select {
	case r.newSessionCh <- newSession{s, id, errCh}:
		return <-errCh
	case <-time.After(timeout):
		return ErrRoomNotRunning
	}
}

// Helper methods
// =============================================================================

// func removeFromSlice(val userID, slice []UserID) error {
// 	for _, i := range slice {
//
// 	}
// }

// Variables
// =============================================================================

const timeout = 200 * time.Millisecond

// ErrRoomNotRunning is returned whenever the service function run by Start()
// is not responding within a second.
var (
	ErrRoomNotRunning = errors.New("room not running")
	ErrUserNotAuth    = errors.New("user not authorized")
	ErrUserNotFound   = errors.New("user not found")
)
