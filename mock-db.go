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

// Position ...
type Position struct {
	Lat   float32
	Lon   float32
	Accur int
	Elev  int
}
