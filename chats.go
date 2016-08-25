package main

import "time"

// Chat is a chat with one or more members.
// Clk is the clock of the last message in the channel.
type Chat struct {
	ID string // Primary key (UUIDv4)

	// These will be set automatically by GORM
	CreatedAt time.Time
	UpdatedAt time.Time

	// Fields set by the users
	Name    string
	Descr   string
	Owner   User
	Members []User
	Clk     int
}
