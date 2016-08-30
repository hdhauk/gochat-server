package main

import "time"

// Message is a chat message with
type Message struct {
	ID string // Primary key (UUIDv4)

	// These will be set automatically by GORM
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	SenderID string
	Priority int
	Text     string // NOTE: Needs to be able to be arbitrary large
	//Pos      Position
}

// User ...
type User struct {
	ID string

	CreatedAt time.Time
	UpdatedAt time.Time

	Nick string
}

// HACK:MOCK DATA
var users = []User{
	User{ID: "1", CreatedAt: time.Now(), UpdatedAt: time.Now(), Nick: "Ola"},
	User{ID: "2", CreatedAt: time.Now(), UpdatedAt: time.Now(), Nick: "Kari"},
	User{ID: "3", CreatedAt: time.Now(), UpdatedAt: time.Now(), Nick: "Petter"},
	User{ID: "4", CreatedAt: time.Now(), UpdatedAt: time.Now(), Nick: "Lise"},
	User{ID: "5", CreatedAt: time.Now(), UpdatedAt: time.Now(), Nick: "Ali"},
}

// Position ...
type Position struct {
	Lat   float32
	Lon   float32
	Accur int
	Elev  int
}

// Chat ...
type Chat struct {
	ID string // Primary key (UUIDv4)

	// These will be set automatically by GORM
	CreatedAt time.Time
	UpdatedAt time.Time

	// Fields set by the users
	Name  string
	Descr string
	Owner User
	//MemberIDs []string // <- FIXME: Sqlite3 cannot handle slices
	Clk int
}

// HACK: MOCK DATA
var chats = []Chat{
	Chat{ID: "1", CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "Chat1", Descr: "This is chat 1", Owner: users[0], Clk: 4},
	Chat{ID: "2", CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "Chat2", Descr: "This is chat 2", Owner: users[1], Clk: 44},
	Chat{ID: "3", CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "Chat3", Descr: "This is chat 3", Owner: users[2], Clk: 2342},
	Chat{ID: "4", CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "Chat4", Descr: "This is chat 4", Owner: users[3], Clk: 22},
}
